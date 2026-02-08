import json
import logging

from langchain_openai import ChatOpenAI
from langchain_core.messages import SystemMessage, HumanMessage

from .config import OPENAI_API_KEY, LLM_MODEL
from .models import (
    CategorizedTransaction,
    CategorizationRequest,
    CategorizationResponse,
    SuggestedCategory,
)

logger = logging.getLogger(__name__)

SYSTEM_PROMPT = """Você é um assistente financeiro especializado em categorização de transações bancárias brasileiras.

Sua tarefa é analisar cada transação e atribuir uma categoria adequada.

## Regras:
1. Para cada transação, tente mapear para uma das categorias existentes fornecidas pelo usuário.
2. Se nenhuma categoria existente fizer sentido semântico, sugira uma nova categoria com um nome descritivo em português.
3. Ao sugerir novas categorias, use nomes padronizados como: "Alimentação", "Transporte", "Saúde", "Lazer", "Educação", "Moradia", "Streaming", "Assinaturas", "Vestuário", "Salário", "Freelance", etc.
4. Limpe e normalize a descrição (capitalização adequada, remover caracteres desnecessários).
5. Determine o tipo da transação baseado no valor e na descrição:
   - Valores negativos ou descrições com "pagamento", "compra", "débito" → expense
   - Valores positivos ou descrições com "recebido", "salário", "transferência recebida" → income
6. Atribua uma confiança (0.0 a 1.0) para cada categorização.
7. Respeite o tipo da categoria: transações income devem ter categoria income, expense devem ter categoria expense.

## Formato de saída:
Retorne APENAS um JSON válido com a seguinte estrutura (sem markdown, sem texto adicional):
{
    "transactions": [
        {
            "original_description": "DESCRICAO ORIGINAL",
            "cleaned_description": "Descrição Limpa",
            "amount": 123.45,
            "transaction_date": "2026-02-01",
            "transaction_type": "expense",
            "category_id": "uuid-da-categoria-existente-ou-null",
            "suggested_category_name": "Nome da Nova Categoria ou null",
            "suggested_category_type": "expense ou null",
            "confidence": 0.95,
            "installments": null,
            "current_installment": null
        }
    ],
    "new_categories_suggested": [
        {"name": "Alimentação", "category_type": "expense"}
    ],
    "summary": "Resumo da análise"
}
"""


def create_llm() -> ChatOpenAI:
    """Cria instância do ChatOpenAI."""
    return ChatOpenAI(
        model=LLM_MODEL,
        api_key=OPENAI_API_KEY,
        temperature=0.1,
    )


async def categorize_transactions(
    request: CategorizationRequest,
) -> CategorizationResponse:
    """Categoriza transações usando LLM."""
    llm = create_llm()

    # Monta a mensagem com os dados do usuário
    categories_str = ""
    if request.existing_categories:
        categories_list = [
            f"- ID: {cat.id}, Nome: {cat.name}, Tipo: {cat.type}"
            for cat in request.existing_categories
        ]
        categories_str = "Categorias existentes do usuário:\n" + "\n".join(
            categories_list
        )
    else:
        categories_str = "O usuário ainda não possui categorias. Sugira novas categorias para todas as transações."

    transactions_list = []
    for tx in request.transactions:
        tx_dict = {
            "description": tx.description,
            "amount": tx.amount,
            "date": tx.transaction_date,
            "type": tx.transaction_type,
        }
        if tx.installments is not None:
            tx_dict["installments"] = tx.installments
            tx_dict["current_installment"] = tx.current_installment
        transactions_list.append(tx_dict)

    user_message = f"""Tipo de importação: {request.import_type}

{categories_str}

Transações para categorizar:
{json.dumps(transactions_list, ensure_ascii=False, indent=2)}

Categorize cada transação seguindo as regras do sistema. Retorne APENAS o JSON."""

    messages = [
        SystemMessage(content=SYSTEM_PROMPT),
        HumanMessage(content=user_message),
    ]

    logger.info(
        "Enviando %d transações para categorização via LLM", len(request.transactions)
    )

    response = await llm.ainvoke(messages)
    content = response.content.strip()

    # Remove possíveis blocos de código markdown
    if content.startswith("```"):
        lines = content.split("\n")
        # Remove primeira e última linha (```json e ```)
        lines = [l for l in lines if not l.strip().startswith("```")]
        content = "\n".join(lines)

    try:
        data = json.loads(content)
    except json.JSONDecodeError as e:
        logger.error("Falha ao parsear resposta da LLM: %s", e)
        logger.error("Conteúdo recebido: %s", content[:500])
        # Fallback: retorna transações sem categorização
        return _fallback_response(request)

    # Valida e parseia a resposta
    try:
        categorized = CategorizationResponse(
            transactions=[
                CategorizedTransaction(**tx) for tx in data.get("transactions", [])
            ],
            new_categories_suggested=[
                SuggestedCategory(**cat)
                for cat in data.get("new_categories_suggested", [])
            ],
            summary=data.get("summary", "Categorização concluída."),
        )
    except Exception as e:
        logger.error("Erro ao validar resposta da LLM: %s", e)
        return _fallback_response(request)

    logger.info("Categorização concluída: %d transações", len(categorized.transactions))
    return categorized


def _fallback_response(request: CategorizationRequest) -> CategorizationResponse:
    """Resposta de fallback quando a LLM falha — retorna transações sem categorização."""
    transactions = []
    for tx in request.transactions:
        transactions.append(
            CategorizedTransaction(
                original_description=tx.description,
                cleaned_description=tx.description,
                amount=tx.amount,
                transaction_date=tx.transaction_date,
                transaction_type=tx.transaction_type,
                category_id=None,
                suggested_category_name="Outros",
                suggested_category_type=tx.transaction_type,
                confidence=0.0,
                installments=tx.installments,
                current_installment=tx.current_installment,
            )
        )

    return CategorizationResponse(
        transactions=transactions,
        new_categories_suggested=[
            SuggestedCategory(name="Outros", category_type="expense"),
            SuggestedCategory(name="Outros", category_type="income"),
        ],
        summary="Categorização automática falhou. Transações retornadas sem categoria para revisão manual.",
    )
