package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/thenopholo/my-money/internal/domain"
)

// ImportService orquestra o fluxo de importação: parse CSV → LLM → validação → persistência.
type ImportService struct {
	llmClient       *LLMClient
	transactionRepo TransactionRepository
	ccTxRepo        CreditCardTransactionRepository
	categoryRepo    CategoryRepository
	bankAccountRepo BankAccountRepository
	creditCardRepo  CreditCardRepository
}

// NewImportService cria uma nova instância de ImportService.
func NewImportService(
	llmClient *LLMClient,
	transactionRepo TransactionRepository,
	ccTxRepo CreditCardTransactionRepository,
	categoryRepo CategoryRepository,
	bankAccountRepo BankAccountRepository,
	creditCardRepo CreditCardRepository,
) *ImportService {
	return &ImportService{
		llmClient:       llmClient,
		transactionRepo: transactionRepo,
		ccTxRepo:        ccTxRepo,
		categoryRepo:    categoryRepo,
		bankAccountRepo: bankAccountRepo,
		creditCardRepo:  creditCardRepo,
	}
}

// Preview parseia CSV, categoriza via LLM e retorna preview sem persistir.
func (s *ImportService) Preview(
	ctx context.Context,
	file io.Reader,
	importType domain.ImportType,
	targetID uuid.UUID,
	userID uuid.UUID,
) (*domain.ImportPreviewResponse, error) {
	if !importType.IsValid() {
		return nil, domain.ErrInvalidImportType
	}

	// Parseia CSV
	var rawTransactions []domain.RawCSVTransaction
	var err error

	switch importType {
	case domain.ImportTypeBank:
		rawTransactions, err = ParseBankCSV(file)
	case domain.ImportTypeCreditCard:
		rawTransactions, err = ParseCreditCardCSV(file)
	}
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}

	if len(rawTransactions) == 0 {
		return nil, domain.ErrNoTransactionsToImport
	}

	// Busca categorias do usuário
	categories, err := s.categoryRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching categories: %w", err)
	}

	// Envia para o agente LLM
	llmResp, err := s.llmClient.Categorize(ctx, rawTransactions, categories, importType)
	if err != nil {
		return nil, fmt.Errorf("LLM categorization: %w", err)
	}

	// Converte resposta para formato de domínio
	preview := toLLMPreviewResponse(llmResp, importType, targetID)

	return preview, nil
}

// Confirm valida e persiste as transações confirmadas pelo usuário.
func (s *ImportService) Confirm(
	ctx context.Context,
	userID uuid.UUID,
	req domain.ImportConfirmRequest,
) (*domain.ImportResult, error) {
	if !req.ImportType.IsValid() {
		return nil, domain.ErrInvalidImportType
	}

	if len(req.Transactions) == 0 {
		return nil, domain.ErrNoTransactionsToImport
	}

	result := &domain.ImportResult{}

	// Cache de categorias criadas nesta sessão (evita duplicatas)
	createdCategories := make(map[string]uuid.UUID) // key: "name|type"

	for i, ct := range req.Transactions {
		// Resolve category_id: cria nova se necessário
		categoryID, created, err := s.resolveCategoryID(ctx, userID, ct, createdCategories)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: %v", i+1, err))
			continue
		}
		if created {
			result.CategoriesCreated++
		}

		// Verifica duplicatas e persiste
		switch req.ImportType {
		case domain.ImportTypeBank:
			isDup, err := s.isDuplicateBankTransaction(ctx, req.TargetID, ct)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: checking duplicate: %v", i+1, err))
				continue
			}
			if isDup {
				result.DuplicatesSkipped++
				continue
			}

			if err := s.createBankTransaction(ctx, req.TargetID, categoryID, ct); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: %v", i+1, err))
				continue
			}
			result.Created++

		case domain.ImportTypeCreditCard:
			isDup, err := s.isDuplicateCCTransaction(ctx, req.TargetID, ct)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: checking duplicate: %v", i+1, err))
				continue
			}
			if isDup {
				result.DuplicatesSkipped++
				continue
			}

			if err := s.createCCTransaction(ctx, req.TargetID, categoryID, ct); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("transaction %d: %v", i+1, err))
				continue
			}
			result.Created++
		}
	}

	return result, nil
}

// resolveCategoryID retorna o category_id a usar: existente ou criado.
func (s *ImportService) resolveCategoryID(
	ctx context.Context,
	userID uuid.UUID,
	ct domain.ConfirmTransaction,
	cache map[string]uuid.UUID,
) (uuid.UUID, bool, error) {
	// Se já tem category_id, usa diretamente
	if ct.CategoryID != nil && *ct.CategoryID != uuid.Nil {
		return *ct.CategoryID, false, nil
	}

	// Se tem nome de nova categoria, cria
	if ct.NewCategoryName != nil && *ct.NewCategoryName != "" {
		catType := domain.CategoryTypeExpense
		if ct.NewCategoryType != nil {
			catType = *ct.NewCategoryType
		}

		// Verifica cache para não criar duplicatas na mesma importação
		cacheKey := fmt.Sprintf("%s|%s", *ct.NewCategoryName, catType)
		if cachedID, ok := cache[cacheKey]; ok {
			return cachedID, false, nil
		}

		category, err := domain.NewCategory(userID, catType, *ct.NewCategoryName, nil, nil)
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("creating category %q: %w", *ct.NewCategoryName, err)
		}

		if err := s.categoryRepo.Create(ctx, category); err != nil {
			return uuid.Nil, false, fmt.Errorf("persisting category %q: %w", *ct.NewCategoryName, err)
		}

		cache[cacheKey] = category.ID
		return category.ID, true, nil
	}

	return uuid.Nil, false, fmt.Errorf("transaction has no category_id and no new_category_name")
}

// isDuplicateBankTransaction verifica se já existe transação idêntica na conta (±1 dia).
func (s *ImportService) isDuplicateBankTransaction(
	ctx context.Context,
	accountID uuid.UUID,
	ct domain.ConfirmTransaction,
) (bool, error) {
	existing, err := s.transactionRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return false, err
	}

	for _, tx := range existing {
		if tx.Description == ct.Description &&
			tx.Amount.Equal(ct.Amount) &&
			isWithinOneDay(tx.TransactionDate, ct.TransactionDate) {
			return true, nil
		}
	}

	return false, nil
}

// isDuplicateCCTransaction verifica se já existe transação idêntica no cartão (±1 dia).
func (s *ImportService) isDuplicateCCTransaction(
	ctx context.Context,
	cardID uuid.UUID,
	ct domain.ConfirmTransaction,
) (bool, error) {
	existing, err := s.ccTxRepo.GetByCardID(ctx, cardID)
	if err != nil {
		return false, err
	}

	for _, tx := range existing {
		if tx.Description == ct.Description &&
			tx.Amount.Equal(ct.Amount) &&
			isWithinOneDay(tx.TransactionDate, ct.TransactionDate) {
			return true, nil
		}
	}

	return false, nil
}

// createBankTransaction cria uma transação bancária.
func (s *ImportService) createBankTransaction(
	ctx context.Context,
	accountID, categoryID uuid.UUID,
	ct domain.ConfirmTransaction,
) error {
	tx, err := domain.NewTransaction(
		accountID,
		categoryID,
		nil, // invoiceID
		nil, // plannedIncome
		nil, // plannedExpense
		ct.Amount,
		ct.TransactionType,
		ct.Description,
		ct.TransactionDate,
	)
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	// Atualiza saldo da conta
	account, err := s.bankAccountRepo.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("fetching account: %w", err)
	}

	switch ct.TransactionType {
	case domain.TransactionTypeIncome:
		if err := account.ApplyIncome(ct.Amount); err != nil {
			return err
		}
	case domain.TransactionTypeExpense:
		if err := account.ApplyExpense(ct.Amount); err != nil {
			return err
		}
	}

	if err := s.bankAccountRepo.Update(ctx, account); err != nil {
		return fmt.Errorf("updating account balance: %w", err)
	}

	if err := s.transactionRepo.Create(ctx, tx); err != nil {
		return fmt.Errorf("persisting transaction: %w", err)
	}

	return nil
}

// createCCTransaction cria uma transação de cartão de crédito.
func (s *ImportService) createCCTransaction(
	ctx context.Context,
	cardID, categoryID uuid.UUID,
	ct domain.ConfirmTransaction,
) error {
	installments := 1
	currentInstallment := 1
	var installmentValue *decimal.Decimal

	if ct.Installments != nil && *ct.Installments > 0 {
		installments = *ct.Installments
	}
	if ct.CurrentInstallment != nil && *ct.CurrentInstallment > 0 {
		currentInstallment = *ct.CurrentInstallment
	}
	if installments > 1 {
		v := ct.Amount.Div(decimal.NewFromInt(int64(installments)))
		installmentValue = &v
	}

	ccTx, err := domain.NewCreditCardTransaction(
		cardID,
		categoryID,
		uuid.Nil, // invoiceID — ficará pendente
		ct.Amount,
		ct.Description,
		installments,
		currentInstallment,
		installmentValue,
	)
	if err != nil {
		return fmt.Errorf("creating cc transaction: %w", err)
	}
	ccTx.TransactionDate = ct.TransactionDate

	if err := s.ccTxRepo.Create(ctx, ccTx); err != nil {
		return fmt.Errorf("persisting cc transaction: %w", err)
	}

	return nil
}

// isWithinOneDay verifica se duas datas estão dentro de uma janela de ±1 dia.
func isWithinOneDay(a, b time.Time) bool {
	diff := a.Sub(b)
	if diff < 0 {
		diff = -diff
	}
	return diff <= 24*time.Hour
}
