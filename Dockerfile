# ==================== STAGE 1: Build ====================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copia dependências primeiro para aproveitar cache de layers
COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte
COPY . .

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/api ./cmd/api

# ==================== STAGE 2: Runtime ====================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

# Cria usuário não-root para segurança
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copia apenas o binário do stage de build
COPY --from=builder /app/api .

# Copia as migrations para rodar no deploy se necessário
COPY --from=builder /app/internal/migrations ./migrations

# Define o usuário não-root
USER appuser

# Porta padrão da aplicação
EXPOSE 4235

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:4235/health || exit 1

ENTRYPOINT ["./api"]
