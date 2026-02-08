FROM python:3.14-slim

WORKDIR /app

# Instala uv para gerenciamento de dependências
RUN pip install --no-cache-dir uv

# Copia arquivos de dependência
COPY pyproject.toml .
COPY .python-version .

# Instala dependências
RUN uv pip install --system -r pyproject.toml

# Copia código-fonte
COPY src/ src/

# Porta padrão do agente
ENV LLM_AGENT_PORT=8001
EXPOSE ${LLM_AGENT_PORT}

# Healthcheck para Docker Compose / orquestradores
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:${LLM_AGENT_PORT}/health')" || exit 1

# Inicia o agente LLM
CMD uvicorn src.main:app --host 0.0.0.0 --port ${LLM_AGENT_PORT}
