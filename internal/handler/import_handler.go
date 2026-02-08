package handler

import (
	"errors"
	"net/http"

	"github.com/thenopholo/my-money/internal/domain"
	"github.com/thenopholo/my-money/internal/service"
)

// maxUploadSize limita o tamanho do upload a 5MB.
const maxUploadSize = 5 * 1024 * 1024

// ImportHandler lida com os endpoints de importação de CSV.
type ImportHandler struct {
	importService  *service.ImportService
	accountService *service.BankAccountService
	cardService    *service.CreditCardService
}

// NewImportHandler cria uma nova instância de ImportHandler.
func NewImportHandler(
	importService *service.ImportService,
	accountService *service.BankAccountService,
	cardService *service.CreditCardService,
) *ImportHandler {
	return &ImportHandler{
		importService:  importService,
		accountService: accountService,
		cardService:    cardService,
	}
}

// Preview recebe um CSV via multipart/form-data, parseia e categoriza via LLM.
func (h *ImportHandler) Preview(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	// Limita tamanho do request
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		Error(w, http.StatusBadRequest, "file too large (max 5MB)")
		return
	}

	// Lê parâmetros do form
	importTypeStr := r.FormValue("import_type")
	importType := domain.ImportType(importTypeStr)
	if !importType.IsValid() {
		Error(w, http.StatusBadRequest, domain.ErrInvalidImportType.Error())
		return
	}

	targetIDStr := r.FormValue("target_id")
	targetID, err := parseUUID(targetIDStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid target_id")
		return
	}

	// Verifica ownership do target
	switch importType {
	case domain.ImportTypeBank:
		account, err := h.accountService.GetByID(r.Context(), targetID)
		if err != nil {
			Error(w, http.StatusNotFound, "account not found")
			return
		}
		if account.UserID != userID {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
	case domain.ImportTypeCreditCard:
		card, err := h.cardService.GetByID(r.Context(), targetID)
		if err != nil {
			Error(w, http.StatusNotFound, "credit card not found")
			return
		}
		if card.UserID != userID {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	// Lê arquivo CSV
	file, _, err := r.FormFile("file")
	if err != nil {
		Error(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	// Chama o service de preview
	preview, err := h.importService.Preview(r.Context(), file, importType, targetID, userID)
	if err != nil {
		h.handleImportError(w, err)
		return
	}

	JSON(w, http.StatusOK, preview)
}

// Confirm recebe as transações categorizadas (editadas pelo usuário) e persiste.
func (h *ImportHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserIDFromContext(w, r)
	if !ok {
		return
	}

	var req domain.ImportConfirmRequest
	if err := decodeJSON(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !req.ImportType.IsValid() {
		Error(w, http.StatusBadRequest, domain.ErrInvalidImportType.Error())
		return
	}

	// Verifica ownership do target
	switch req.ImportType {
	case domain.ImportTypeBank:
		account, err := h.accountService.GetByID(r.Context(), req.TargetID)
		if err != nil {
			Error(w, http.StatusNotFound, "account not found")
			return
		}
		if account.UserID != userID {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
	case domain.ImportTypeCreditCard:
		card, err := h.cardService.GetByID(r.Context(), req.TargetID)
		if err != nil {
			Error(w, http.StatusNotFound, "credit card not found")
			return
		}
		if card.UserID != userID {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	result, err := h.importService.Confirm(r.Context(), userID, req)
	if err != nil {
		h.handleImportError(w, err)
		return
	}

	JSON(w, http.StatusCreated, result)
}

// handleImportError mapeia erros de domínio para HTTP status codes.
func (h *ImportHandler) handleImportError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidImportType):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidCSVFormat):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmptyCSV):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrCSVTooLarge):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNoTransactionsToImport):
		Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrLLMUnavailable):
		Error(w, http.StatusServiceUnavailable, "categorization service unavailable")
	case errors.Is(err, domain.ErrLLMTimeout):
		Error(w, http.StatusGatewayTimeout, "categorization service timed out")
	default:
		Error(w, http.StatusInternalServerError, "internal error")
	}
}
