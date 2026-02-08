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

EXPOSE 8001

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8001"]
