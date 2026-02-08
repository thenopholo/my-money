package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ImportType representa o tipo de importação
type ImportType string

const (
	ImportTypeBank       ImportType = "bank_account"
	ImportTypeCreditCard ImportType = "credit_card"
)

func (t ImportType) IsValid() bool {
	return t == ImportTypeBank || t == ImportTypeCreditCard
}

// RawCSVTransaction — transação bruta parseada do CSV
type RawCSVTransaction struct {
	Description     string          `json:"description"`
	Amount          decimal.Decimal `json:"amount"`
	TransactionDate time.Time       `json:"transaction_date"`
	TransactionType TransactionType `json:"transaction_type"`
	Installments    *int            `json:"installments,omitempty"`
	CurrentInstall  *int            `json:"current_installment,omitempty"`
}

// CategorizedTransaction — transação categorizada pela LLM (preview)
type CategorizedTransaction struct {
	OriginalDescription   string          `json:"original_description"`
	CleanedDescription    string          `json:"cleaned_description"`
	Amount                decimal.Decimal `json:"amount"`
	TransactionDate       time.Time       `json:"transaction_date"`
	TransactionType       TransactionType `json:"transaction_type"`
	CategoryID            *uuid.UUID      `json:"category_id"`
	SuggestedCategoryName *string         `json:"suggested_category_name"`
	SuggestedCategoryType *CategoryType   `json:"suggested_category_type"`
	Confidence            float64         `json:"confidence"`
	Installments          *int            `json:"installments,omitempty"`
	CurrentInstallment    *int            `json:"current_installment,omitempty"`
}

// ImportPreviewResponse — resposta do endpoint de preview
type ImportPreviewResponse struct {
	ImportType             ImportType               `json:"import_type"`
	TargetID               uuid.UUID                `json:"target_id"`
	Transactions           []CategorizedTransaction `json:"transactions"`
	NewCategoriesSuggested []SuggestedCategory      `json:"new_categories_suggested"`
	Summary                string                   `json:"summary"`
	TotalTransactions      int                      `json:"total_transactions"`
	TotalAmount            decimal.Decimal          `json:"total_amount"`
}

// SuggestedCategory — categoria nova sugerida pela LLM
type SuggestedCategory struct {
	Name         string       `json:"name"`
	CategoryType CategoryType `json:"category_type"`
}

// ImportConfirmRequest — request do endpoint de confirmação
type ImportConfirmRequest struct {
	ImportType   ImportType           `json:"import_type"`
	TargetID     uuid.UUID            `json:"target_id"`
	Transactions []ConfirmTransaction `json:"transactions"`
}

// ConfirmTransaction — transação confirmada pelo usuário (pode ter sido editada)
type ConfirmTransaction struct {
	Description        string          `json:"description"`
	Amount             decimal.Decimal `json:"amount"`
	TransactionDate    time.Time       `json:"transaction_date"`
	TransactionType    TransactionType `json:"transaction_type"`
	CategoryID         *uuid.UUID      `json:"category_id"`
	NewCategoryName    *string         `json:"new_category_name"`
	NewCategoryType    *CategoryType   `json:"new_category_type"`
	Installments       *int            `json:"installments,omitempty"`
	CurrentInstallment *int            `json:"current_installment,omitempty"`
}

// ImportResult — resultado da importação
type ImportResult struct {
	Created           int      `json:"created"`
	DuplicatesSkipped int      `json:"duplicates_skipped"`
	CategoriesCreated int      `json:"categories_created"`
	Errors            []string `json:"errors,omitempty"`
}
