import logging

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware

from .agent import categorize_transactions
from .config import OPENAI_API_KEY, LLM_AGENT_PORT
from .models import CategorizationRequest, CategorizationResponse

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="My Money LLM Agent",
    description="Agente de categorização de transações financeiras via LLM",
    version="0.1.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
async def health():
    """Health check endpoint."""
    has_key = bool(OPENAI_API_KEY)
    return {
        "status": "ok" if has_key else "degraded",
        "openai_key_configured": has_key,
    }


@app.post("/categorize", response_model=CategorizationResponse)
async def categorize(request: CategorizationRequest):
    """Categoriza transações usando LLM."""
    if not OPENAI_API_KEY:
        raise HTTPException(
            status_code=503,
            detail="OPENAI_API_KEY not configured",
        )

    if not request.transactions:
        raise HTTPException(
            status_code=400,
            detail="No transactions provided",
        )

    try:
        result = await categorize_transactions(request)
        return result
    except Exception as e:
        logger.error("Erro na categorização: %s", e, exc_info=True)
        raise HTTPException(
            status_code=500,
            detail="Internal categorization error",
        ) from e


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "src.main:app",
        host="0.0.0.0",
        port=LLM_AGENT_PORT,
        reload=True,
    )
