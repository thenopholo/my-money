package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ResetService struct {
	transactionRepo TransactionRepository
	ccTxRepo        CreditCardTransactionRepository
	accountRepo     BankAccountRepository
	creditCardRepo  CreditCardRepository
}

func NewResetService(
	transactionRepo TransactionRepository,
	ccTxRepo CreditCardTransactionRepository,
	accountRepo BankAccountRepository,
	creditCardRepo CreditCardRepository,
) *ResetService {
	return &ResetService{
		transactionRepo: transactionRepo,
		ccTxRepo:        ccTxRepo,
		accountRepo:     accountRepo,
		creditCardRepo:  creditCardRepo,
	}
}

func (s *ResetService) ResetAllTransactions(ctx context.Context, userID uuid.UUID) error {
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if err := s.transactionRepo.DeleteAllByAccountID(ctx, account.ID); err != nil {
			return err
		}

		if err := account.SetBalance(decimal.Zero); err != nil {
			return err
		}

		if err := s.accountRepo.Update(ctx, account); err != nil {
			return err
		}
	}

	cards, err := s.creditCardRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, card := range cards {
		if err := s.ccTxRepo.DeleteAllByCardID(ctx, card.ID); err != nil {
			return err
		}
	}

	return nil
}
