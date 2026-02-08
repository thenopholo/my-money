package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/thenopholo/my-money/internal/domain"
)

const (
	llmTimeout    = 120 * time.Second // 2 minutos — CSVs com muitas transações podem demorar na OpenAI
	llmMaxRetries = 3
)

// LLMClient é o cliente HTTP para comunicação com o agente Python de categorização.
type LLMClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewLLMClient cria uma nova instância de LLMClient.
func NewLLMClient(baseURL string) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: llmTimeout,
		},
	}
}

// llmCategorizationRequest é o body enviado para o agente Python.
type llmCategorizationRequest struct {
	Transactions       []llmRawTransaction    `json:"transactions"`
	ExistingCategories []llmExistingCategory  `json:"existing_categories"`
	ImportType         string                 `json:"import_type"`
}

type llmRawTransaction struct {
	Description        string  `json:"description"`
	Amount             float64 `json:"amount"`
	TransactionDate    string  `json:"transaction_date"`
	TransactionType    string  `json:"transaction_type"`
	Installments       *int    `json:"installments,omitempty"`
	CurrentInstallment *int    `json:"current_installment,omitempty"`
}

type llmExistingCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// llmCategorizationResponse é a response recebida do agente Python.
type llmCategorizationResponse struct {
	Transactions           []llmCategorizedTransaction `json:"transactions"`
	NewCategoriesSuggested []llmSuggestedCategory      `json:"new_categories_suggested"`
	Summary                string                      `json:"summary"`
}

type llmCategorizedTransaction struct {
	OriginalDescription   string   `json:"original_description"`
	CleanedDescription    string   `json:"cleaned_description"`
	Amount                float64  `json:"amount"`
	TransactionDate       string   `json:"transaction_date"`
	TransactionType       string   `json:"transaction_type"`
	CategoryID            *string  `json:"category_id"`
	SuggestedCategoryName *string  `json:"suggested_category_name"`
	SuggestedCategoryType *string  `json:"suggested_category_type"`
	Confidence            float64  `json:"confidence"`
	Installments          *int     `json:"installments,omitempty"`
	CurrentInstallment    *int     `json:"current_installment,omitempty"`
}

type llmSuggestedCategory struct {
	Name         string `json:"name"`
	CategoryType string `json:"category_type"`
}

// Categorize envia transações para o agente Python e retorna categorizadas.
func (c *LLMClient) Categorize(
	ctx context.Context,
	transactions []domain.RawCSVTransaction,
	categories []*domain.Category,
	importType domain.ImportType,
) (*llmCategorizationResponse, error) {
	// Monta request — inicializa slices como arrays vazios para evitar null no JSON
	reqBody := llmCategorizationRequest{
		Transactions:       make([]llmRawTransaction, 0, len(transactions)),
		ExistingCategories: make([]llmExistingCategory, 0, len(categories)),
		ImportType:         string(importType),
	}

	for _, tx := range transactions {
		amount, _ := tx.Amount.Float64()
		reqBody.Transactions = append(reqBody.Transactions, llmRawTransaction{
			Description:        tx.Description,
			Amount:             amount,
			TransactionDate:    tx.TransactionDate.Format("2006-01-02"),
			TransactionType:    string(tx.TransactionType),
			Installments:       tx.Installments,
			CurrentInstallment: tx.CurrentInstall,
		})
	}

	for _, cat := range categories {
		reqBody.ExistingCategories = append(reqBody.ExistingCategories, llmExistingCategory{
			ID:   cat.ID.String(),
			Name: cat.Name,
			Type: string(cat.CategoryType),
		})
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling LLM request: %w", err)
	}

	// Retry com backoff exponencial
	var lastErr error
	for attempt := 0; attempt < llmMaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			log.Printf("LLM retry %d/%d após %v", attempt+1, llmMaxRetries, backoff)

			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("LLM categorization: %w", ctx.Err())
			case <-time.After(backoff):
			}
		}

		start := time.Now()
		resp, err := c.doRequest(ctx, jsonBody)
		elapsed := time.Since(start)
		log.Printf("LLM categorization request took %v", elapsed)

		if err != nil {
			lastErr = err
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("LLM categorization failed after %d retries: %w", llmMaxRetries, lastErr)
}

// doRequest executa a request HTTP para o agente Python.
func (c *LLMClient) doRequest(ctx context.Context, jsonBody []byte) (*llmCategorizationResponse, error) {
	url := c.baseURL + "/categorize"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating LLM request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLLMUnavailable, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading LLM response: %w", err)
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: status %d", domain.ErrLLMUnavailable, resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("LLM agent error: status=%d body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("LLM returned status %d: %s", resp.StatusCode, string(body))
	}

	var result llmCategorizationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding LLM response: %w", err)
	}

	return &result, nil
}

// toLLMPreviewResponse converte a resposta da LLM para o formato de domínio.
func toLLMPreviewResponse(
	resp *llmCategorizationResponse,
	importType domain.ImportType,
	targetID uuid.UUID,
) *domain.ImportPreviewResponse {
	preview := &domain.ImportPreviewResponse{
		ImportType:        importType,
		TargetID:          targetID,
		Summary:           resp.Summary,
		TotalTransactions: len(resp.Transactions),
		TotalAmount:       decimal.Zero,
	}

	for _, tx := range resp.Transactions {
		ct := domain.CategorizedTransaction{
			OriginalDescription: tx.OriginalDescription,
			CleanedDescription:  tx.CleanedDescription,
			Amount:              decimal.NewFromFloat(tx.Amount),
			TransactionType:     domain.TransactionType(tx.TransactionType),
			Confidence:          tx.Confidence,
			Installments:        tx.Installments,
			CurrentInstallment:  tx.CurrentInstallment,
		}

		// Parseia data
		if t, err := time.Parse("2006-01-02", tx.TransactionDate); err == nil {
			ct.TransactionDate = t
		}

		// Parseia category_id
		if tx.CategoryID != nil && *tx.CategoryID != "" && *tx.CategoryID != "null" {
			if id, err := uuid.Parse(*tx.CategoryID); err == nil {
				ct.CategoryID = &id
			}
		}

		ct.SuggestedCategoryName = tx.SuggestedCategoryName
		if tx.SuggestedCategoryType != nil {
			catType := domain.CategoryType(*tx.SuggestedCategoryType)
			ct.SuggestedCategoryType = &catType
		}

		preview.Transactions = append(preview.Transactions, ct)
		preview.TotalAmount = preview.TotalAmount.Add(ct.Amount)
	}

	for _, sc := range resp.NewCategoriesSuggested {
		preview.NewCategoriesSuggested = append(preview.NewCategoriesSuggested, domain.SuggestedCategory{
			Name:         sc.Name,
			CategoryType: domain.CategoryType(sc.CategoryType),
		})
	}

	return preview
}
