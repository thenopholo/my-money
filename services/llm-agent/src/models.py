from pydantic import BaseModel, Field


class RawTransaction(BaseModel):
    """Transação bruta recebida do backend Go."""

    description: str
    amount: float
    transaction_date: str  # ISO 8601
    transaction_type: str  # "income" | "expense"
    installments: int | None = None
    current_installment: int | None = None


class ExistingCategory(BaseModel):
    """Categoria existente do usuário."""

    id: str
    name: str
    type: str  # "income" | "expense"


class CategorizationRequest(BaseModel):
    """Request recebida do backend Go."""

    transactions: list[RawTransaction]
    existing_categories: list[ExistingCategory] = []
    import_type: str  # "bank_account" | "credit_card"


class CategorizedTransaction(BaseModel):
    """Transação categorizada pela LLM."""

    original_description: str
    cleaned_description: str
    amount: float
    transaction_date: str
    transaction_type: str
    category_id: str | None = None
    suggested_category_name: str | None = None
    suggested_category_type: str | None = None
    confidence: float = Field(ge=0.0, le=1.0)
    installments: int | None = None
    current_installment: int | None = None


class SuggestedCategory(BaseModel):
    """Nova categoria sugerida pela LLM."""

    name: str
    category_type: str


class CategorizationResponse(BaseModel):
    """Response retornada para o backend Go."""

    transactions: list[CategorizedTransaction]
    new_categories_suggested: list[SuggestedCategory]
    summary: str
