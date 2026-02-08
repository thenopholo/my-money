# 💰 My Money

**API REST de planejamento financeiro pessoal**, construída em Go com Clean Architecture.

> Quanto você ganha? Quanto você gasta? Quanto sobra (se sobrar)?

My Money resolve um problema real: a maioria dos apps financeiros sobrecarrega o usuário com dezenas de features que só fazem sentido pra quem já entende de finanças. Gráficos de rendimento composto, projeções de portfólio, análise de volatilidade… o básico fica enterrado em menus infinitos.

Esse projeto parte de uma premissa simples: a maioria das pessoas só precisa registrar **o que entra, o que sai e acompanhar o saldo**. Contas bancárias, cartões de crédito, receitas previstas, despesas planejadas — tudo organizado de um jeito que qualquer pessoa entende. Sem MBA em finanças como pré-requisito.

---

## 📋 Índice

- [Funcionalidades](#-funcionalidades)
- [Stack Tecnológica](#-stack-tecnológica)
- [Arquitetura](#-arquitetura)
- [Estrutura de Diretórios](#-estrutura-de-diretórios)
- [Padrões e Decisões Técnicas](#-padrões-e-decisões-técnicas)
- [Segurança e Prevenção](#-segurança-e-prevenção)
- [Testes](#-testes)
- [Endpoints da API](#-endpoints-da-api)
- [Rodando o Projeto](#-rodando-o-projeto)
- [Variáveis de Ambiente](#-variáveis-de-ambiente)

---

## ✨ Funcionalidades

| Domínio | O que faz |
|---------|-----------|
| **Autenticação** | Registro, login com JWT (HMAC-SHA256), troca de senha com validação de força |
| **Contas Bancárias** | CRUD de contas corrente/poupança com controle de saldo automático |
| **Categorias** | Categorização de receitas e despesas por tipo |
| **Cartões de Crédito** | Gerenciamento de cartões com limite, dia de fechamento e vencimento |
| **Transações** | Registro de receitas/despesas com impacto automático no saldo da conta |
| **Faturas** | Fechamento mensal de fatura do cartão, pagamento via conta bancária |
| **Transações de Cartão** | Compras no cartão com parcelamento e vínculo automático à fatura |
| **Receitas/Despesas Planejadas** | Agendamento de receitas e despesas recorrentes (mensal, anual, avulsa) com geração de transações no vencimento |

---

## 🛠 Stack Tecnológica

| Categoria | Tecnologia |
|-----------|------------|
| **Linguagem** | Go 1.25 |
| **HTTP Router** | [chi v5](https://github.com/go-chi/chi) |
| **Banco de Dados** | PostgreSQL 17 |
| **Driver DB** | [pgx v5](https://github.com/jackc/pgx) + `pgxpool` |
| **Geração de Queries** | [sqlc](https://sqlc.dev) |
| **Migrations** | [tern](https://github.com/jackc/tern) |
| **Autenticação** | [golang-jwt v5](https://github.com/golang-jwt/jwt) (HMAC-SHA256) |
| **Senhas** | `bcrypt` via `golang.org/x/crypto` |
| **Decimais** | [shopspring/decimal](https://github.com/shopspring/decimal) |
| **UUIDs** | [google/uuid](https://github.com/google/uuid) |
| **CORS** | [go-chi/cors](https://github.com/go-chi/cors) |
| **Hot Reload** | [air](https://github.com/air-verse/air) |
| **Container** | Docker + Docker Compose |

---

## 🏗 Arquitetura

O projeto segue **Clean Architecture** com separação estrita de camadas e **Dependency Inversion Principle** — as dependências sempre apontam para o domínio.

```
cmd/api (composição / DI manual)
    ↓
handler → service → domain
    ↓        ↓
middleware  repository → domain
               ↓
           postgres (sqlc)
```

### Fluxo de uma Request

```
HTTP Request
    → Middleware (Logger, Recoverer, CORS, Auth JWT)
        → Handler (decode, ownership check, delegate)
            → Service (regras de negócio, orquestração)
                → Repository (persistência, conversão de tipos)
                    → PostgreSQL (via sqlc)
```

### Princípios Aplicados

- **Domain não conhece infraestrutura** — apenas stdlib, `uuid` e `decimal`
- **Service define interfaces** — o consumidor define o contrato (`service/interface.go`)
- **Repository implementa interfaces** — traduz tipos do banco para domínio
- **Handler é fino** — só faz decode, verificação de ownership e chama o service
- **DI manual no `main.go`** — sem frameworks (wire, fx, dig), composição explícita

---

## 📁 Estrutura de Diretórios

```
my-money/
├── cmd/
│   └── api/
│       └── main.go                  # Entrypoint: bootstrap, DI, server
├── internal/                        # Código privado (não importável externamente)
│   ├── auth/                        # JWT (generate, verify)
│   ├── config/                      # Carregamento de variáveis de ambiente
│   ├── domain/                      # Entidades, validações, tipos, erros sentinela
│   ├── handler/                     # Handlers HTTP, router, helpers de response
│   │   └── middleware/              # Middleware de autenticação
│   ├── migrations/                  # SQL migrations (tern)
│   ├── repository/                  # Implementações concretas de repositório
│   │   ├── postgres/                # Código gerado pelo sqlc (NÃO editar)
│   │   └── queries/                 # Queries SQL fonte para o sqlc
│   └── service/                     # Lógica de negócio, interfaces de repositório
├── docker-compose.yml
├── Dockerfile                       # Multi-stage build
├── Makefile
├── sqlc.yml
└── go.mod
```

> Cada camada vive no seu diretório. Cada entidade tem seu arquivo. Cada arquivo tem seu teste. Sem surpresas.

---

## 🔧 Padrões e Decisões Técnicas

### Entidades de Domínio com Validação no Constructor

Toda entidade é criada via constructor `New<Entidade>(...)` que valida os dados antes de retornar. Se a validação falhar, retorna um erro sentinela — nunca um objeto inválido.

```go
// internal/domain/bank_account.go

type AccountType string

const (
    AccountTypeChecking AccountType = "checking"
    AccountTypeSavings  AccountType = "savings"
)

func (a AccountType) IsValid() bool {
    return a == AccountTypeChecking || a == AccountTypeSavings
}

func NewBankAccount(userID uuid.UUID, accountType AccountType, name, bankName string, balance decimal.Decimal) (*BankAccount, error) {
    if !accountType.IsValid() {
        return nil, ErrInvalidAccountType
    }
    // ... validações e construção
}
```

O tipo `AccountType` é um **tipo semântico** (`type AccountType string`) com método `IsValid()` — padrão aplicado em todos os enumerados do projeto (`TransactionType`, `InvoiceStatus`, `Recurrence`).

### Lógica de Negócio na Entidade

Operações que alteram o estado da entidade ficam na própria struct, não no service:

```go
// internal/domain/bank_account.go

func (ba *BankAccount) ApplyIncome(amount decimal.Decimal) error {
    if amount.LessThanOrEqual(decimal.Zero) {
        return ErrInvalidAmount
    }
    return ba.SetBalance(ba.Balance.Add(amount))
}

func (ba *BankAccount) ApplyExpense(amount decimal.Decimal) error {
    if amount.LessThanOrEqual(decimal.Zero) {
        return ErrInvalidAmount
    }
    return ba.SetBalance(ba.Balance.Sub(amount))
}
```

O service orquestra, a entidade decide.

### Erros Sentinela Centralizados

Todos os erros de domínio vivem em `internal/domain/erros.go` como variáveis `var`:

```go
// internal/domain/erros.go

var (
    // Erros de User
    ErrInvalidEmail       = errors.New("invalid email format")
    ErrUserNotFound       = errors.New("user not found in DB")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrSamePassword       = errors.New("new password must be different from current password")

    // Erros de BankAccount
    ErrNegativeBalance     = errors.New("balance cannot be negative for savings account")
    ErrBankAccountNotFound = errors.New("bank account not found")

    // Erros de Transaction
    ErrCategoryTypeMismatch = errors.New("category type does not match transaction type")
    ErrInsufficientBalance  = errors.New("insufficient balance")
    // ...
)
```

Comparação sempre com `errors.Is()`, propagação com `fmt.Errorf("contexto: %w", err)`.

### Interface Definida pelo Consumidor

As interfaces de repositório ficam no pacote `service` (quem consome), não no pacote `repository` (quem implementa). Isso segue o princípio de inversão de dependência do Go idiomático:

```go
// internal/service/interface.go

type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type BankAccountRepository interface {
    Create(ctx context.Context, ba *domain.BankAccount) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.BankAccount, error)
    GetByUserID(ctx context.Context, id uuid.UUID) ([]*domain.BankAccount, error)
    Update(ctx context.Context, ba *domain.BankAccount) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Tradução de Erros no Repository

O repositório encapsula erros de infraestrutura e traduz para erros de domínio:

```go
// internal/repository/user_repository.go

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    u, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrUserNotFound  // tradução pgx → domínio
        }
        return nil, err
    }

    return &domain.User{
        ID:           u.ID,
        Name:         u.Name,
        Email:        u.Email,
        PasswordHash: u.PasswordHash,
        CreatedAt:    u.CreatedAt.Time,
        UpdatedAt:    u.UpdatedAt.Time,
    }, nil
}
```

Nenhuma camada acima do repository sabe que o PostgreSQL existe.

### Injeção de Dependências Manual

Toda composição acontece no `main.go`, sem frameworks. A ordem é previsível e explícita:

```go
// cmd/api/main.go

func main() {
    cfg, err := config.Load()
    // ...
    pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
    defer pool.Close()

    queries := postgres.New(pool)
    jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTDuration)

    // Repositories
    userRepo := repository.NewUserRepository(queries)
    bankAccountRepo := repository.NewBankAccountRepository(queries)
    // ...

    // Services
    userService := service.NewUserService(userRepo)
    bankAccountService := service.NewBankAccountService(bankAccountRepo)
    transactionService := service.NewTransactionService(
        transactionRepo, bankAccountRepo, categoryRepo,
        plannedIncomeRepo, plannedExpenseRepo, invoiceRepo,
    )

    // Handlers (ownership indireta: recebe múltiplos services)
    transactionHandler := handler.NewTransactionHandler(transactionService, bankAccountService)
    invoiceHandler := handler.NewInvoiceHandler(invoiceService, creditCardService)

    // Router
    router := handler.NewRouter(
        userHandler, bankAccountHandler, categoryHandler, /* ... */
        handlermw.Auth(jwtManager), cfg.CORSAllowedOrigins,
    )

    http.ListenAndServe(":"+cfg.Port, router)
}
```

### Geração Type-Safe de SQL com sqlc

Queries SQL são escritas manualmente e o `sqlc` gera código Go type-safe:

```sql
-- internal/repository/queries/users.sql

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash)
VALUES (sqlc.arg(name), sqlc.arg(email), sqlc.arg(password_hash))
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at, updated_at
FROM users WHERE email = $1;
```

A configuração mapeia tipos do PostgreSQL para tipos Go:

```yaml
# sqlc.yml
overrides:
  - db_type: "uuid"
    go_type:
      import: "github.com/google/uuid"
      type: "UUID"
  - db_type: "numeric"
    go_type:
      import: "github.com/shopspring/decimal"
      type: "Decimal"
```

### Propagação de `context.Context`

O contexto é propagado da request HTTP até o repositório, garantindo cancelamento e timeout em todas as camadas:

```
r.Context() → handler → service(ctx) → repository(ctx) → sqlc queries(ctx)
```

Toda função que faz I/O recebe `ctx context.Context` como primeiro parâmetro — sem exceção.

### Docker Multi-Stage Build

O `Dockerfile` usa multi-stage build para gerar uma imagem mínima (~15MB):

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/api ./cmd/api

# Stage 2: Runtime
FROM alpine:3.21
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/api .
USER appuser
EXPOSE 4235
ENTRYPOINT ["./api"]
```

Flags `-s -w` no `ldflags` removem a tabela de símbolos e info de debug. O runtime roda como **usuário não-root**.

---

## 🛡 Segurança e Prevenção

Segurança não é um feature isolado — está espalhada em cada camada do projeto. Abaixo estão os mecanismos concretos implementados.

### Autenticação JWT com Verificação de Algoritmo

O `JWTManager` assina tokens com **HMAC-SHA256** e, no `Verify()`, verifica explicitamente o método de assinatura para prevenir **algorithm confusion attacks** (um atacante trocar o algoritmo para `none` ou `RS256`):

```go
// internal/auth/jwt.go

func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected token signing method")
        }
        return m.secretKey, nil
    })
    // ...
}
```

Claims customizadas incluem `UserID` e `Email`, com expiração configurável. O token é extraído do header `Authorization: Bearer <token>` pelo middleware.

### Middleware de Autenticação

Toda rota sob `/api/*` passa pelo middleware que valida o JWT, extrai o `userID` e injeta no `context.Context` via `context.WithValue`:

```go
// internal/handler/middleware/auth.go

func Auth(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Verifica existência do header Authorization
            // 2. Valida formato "Bearer <token>"
            // 3. Verifica e decodifica o JWT
            // 4. Parsea o userID das claims
            // 5. Injeta no contexto via context.WithValue

            ctx := context.WithValue(r.Context(), userIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

A chave do contexto é um **tipo privado** (`type contextKey string`) — impede colisão com outras bibliotecas.

### Autorização via Ownership (Direta e Indireta)

Após autenticação, todo handler verifica que o recurso pertence ao usuário do token **antes** de qualquer leitura ou escrita. O projeto implementa dois padrões:

**Ownership direta** — o recurso tem `UserID` próprio (ex.: `BankAccount`, `Category`, `CreditCard`):

```go
account, _ := h.bankAccountService.GetByID(r.Context(), accountID)
if account.UserID != userID {
    Error(w, http.StatusForbidden, "forbidden")
    return
}
```

**Ownership indireta** — o recurso pertence a outro recurso do usuário (ex.: `Transaction` → `BankAccount`, `Invoice` → `CreditCard`):

```go
// Transaction: verifica ownership via BankAccount
tx, _ := h.service.GetByID(r.Context(), id)
account, _ := h.accountService.GetByID(r.Context(), tx.AccountID)
if account.UserID != userID {
    Error(w, http.StatusForbidden, "forbidden")
    return
}
```

Para ownership indireta, o handler **recebe o service da entidade-pai** como dependência adicional no constructor:

```go
type TransactionHandler struct {
    service        *service.TransactionService
    accountService *service.BankAccountService  // para verificar ownership
}

// Mesmo padrão em:
// InvoiceHandler         → CreditCardService
// CreditCardTxHandler    → CreditCardService
```

> Nunca se confia em IDs enviados no body para autorização — a verificação é sempre feita pelo token JWT.

### Senhas com bcrypt e Validação de Força

Senhas são hasheadas com **bcrypt** (`DefaultCost`) e validadas com regras de complexidade antes do hash:

```go
// internal/service/password_validation.go

func validatePasswordStrength(password string) error {
    if len(password) < 8 {
        return domain.ErrPasswordTooShort
    }

    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):  hasUpper = true
        case unicode.IsLower(char):  hasLower = true
        case unicode.IsDigit(char):  hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char): hasSpecial = true
        }
    }

    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        return domain.ErrPasswordTooWeak
    }
    return nil
}
```

**Regras**: mínimo 8 caracteres, maiúscula, minúscula, número e caractere especial. Na troca de senha, a nova senha não pode ser igual à atual (`ErrSamePassword`).

### Proteção de Dados Sensíveis via `json:"-"`

Campos sensíveis **nunca** são serializados em responses JSON. A tag `json:"-"` impede o `json.Marshal()` de incluí-los:

```go
type User struct {
    ID           uuid.UUID `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`  // NUNCA serializado
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

Mesmo que o handler passe o user inteiro para `JSON()`, o `PasswordHash` não aparece na response. No `Register`, o handler retorna explicitamente apenas os campos públicos:

```go
JSON(w, http.StatusCreated, map[string]any{
    "id":    user.ID,
    "name":  user.Name,
    "email": user.Email,
})
```

### Validação Rígida de Input

Todos os inputs passam por validação antes de serem processados:

- **`DisallowUnknownFields()`** — rejeita campos JSON desconhecidos (previne mass assignment):
  ```go
  func decodeJSON(r *http.Request, dst any) error {
      dec := json.NewDecoder(r.Body)
      dec.DisallowUnknownFields()
      return dec.Decode(dst)
  }
  ```
- **Helpers tipados** — `parseUUID()`, `parseDecimal()`, `parseDate()` validam o formato antes de prosseguir
- **Structs anônimas inline** — cada handler define exatamente os campos que aceita, nada a mais

### Prevenção de SQL Injection

Todas as queries SQL são geradas pelo **sqlc** com parâmetros bind (`$1`, `$2` ou `sqlc.arg()`). Nenhuma query é construída por concatenação de strings:

```sql
-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at, updated_at
FROM users WHERE email = $1;
```

### Stack de Middlewares de Segurança

O router aplica uma cadeia de middlewares globais antes de qualquer handler:

```go
r.Use(chimiddleware.Logger)     // log de todas as requests
r.Use(chimiddleware.Recoverer)  // captura panics → 500 em vez de crash
r.Use(chimiddleware.RequestID)  // rastreabilidade com X-Request-ID
r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   corsAllowedOrigins,  // configurável via env
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
    AllowCredentials: true,
    MaxAge:           300,
}))
```

Origens CORS são configuráveis via `CORS_ALLOWED_ORIGINS` (nunca `*` em produção).

### Secrets e Variáveis de Ambiente

- `.env` no `.gitignore` — nunca commitado
- `config.Load()` valida variáveis obrigatórias com erro claro:
  ```go
  if cfg.DatabaseURL == "" {
      return nil, fmt.Errorf("missing required env var: DATABASE_URL")
  }
  ```
- Secrets nunca logados — apenas presença/ausência

### Docker: Usuário Não-Root

O container de produção roda como **usuário não-root** com binário estático:

```dockerfile
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
```

### Resumo de Proteções

| Ameaça | Prevenção |
|--------|-----------|
| Algorithm confusion (JWT) | Verificação explícita do método de assinatura no `Verify()` |
| Token forjado/expirado | Middleware valida JWT em toda rota `/api/*` |
| Acesso a recurso alheio | Ownership direta + indireta verificada em todo handler |
| Senha fraca | Validação de complexidade (8+ chars, upper, lower, digit, special) |
| Vazamento de hash de senha | `json:"-"` na struct + response explícita no Register/Login |
| Mass assignment | `DisallowUnknownFields()` + structs anônimas inline |
| SQL Injection | Queries geradas pelo sqlc com bind parameters |
| Panic não tratado | `chimiddleware.Recoverer` no router |
| Secrets no código | `.env` no `.gitignore`, validação no `config.Load()` |
| CORS aberto | Origens configuráveis, nunca `*` em produção |
| Container como root | `USER appuser` no Dockerfile |

---

## 🧪 Testes

O projeto usa **exclusivamente a stdlib `testing`** — sem testify, sem gomega, sem mockgen. Isso é intencional: menos dependências, mais controle, e demonstra domínio das ferramentas nativas do Go.

```bash
# Rodar todos os testes
make test

# Testes com cobertura
make test-cover

# Testes com race detector
go test -race -cover -count=1 ./internal/...
```

### Table-Driven Tests com Subtests Paralelos

O padrão principal de teste no projeto é **table-driven** com `t.Run()` e `t.Parallel()`. Cada cenário é um caso na tabela, com nome descritivo e comparação via `errors.Is()`:

```go
// internal/domain/bank_account_test.go

func TestNewBankAccount_Validations(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name        string
        accountType AccountType
        bankName    string
        balance     decimal.Decimal
        wantErr     error
    }{
        {
            name:        "deve criar conta checking com saldo positivo",
            accountType: AccountTypeChecking,
            bankName:    "Nubank",
            balance:     decimal.NewFromInt(1000),
            wantErr:     nil,
        },
        {
            name:        "deve criar conta checking com saldo negativo (permite overdraft)",
            accountType: AccountTypeChecking,
            bankName:    "Nubank",
            balance:     decimal.NewFromInt(-500),
            wantErr:     nil,
        },
        {
            name:        "deve retornar erro para saldo negativo em poupança",
            accountType: AccountTypeSavings,
            bankName:    "Inter",
            balance:     decimal.NewFromInt(-100),
            wantErr:     ErrNegativeBalance,
        },
        // ... mais cenários
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            ba, err := NewBankAccount(userID, tt.accountType, "Conta", tt.bankName, tt.balance)

            if !errors.Is(err, tt.wantErr) {
                t.Errorf("NewBankAccount() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Mocks Manuais com Function Fields

Services são testados com **mocks manuais** que implementam as interfaces de `service/interface.go`. Cada método é um **function field**, configurado por cenário — sem libs externas:

```go
// internal/service/mock_user_repo_test.go

type mockUserRepository struct {
    createFn     func(ctx context.Context, user *domain.User) error
    getByIDFn    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
    getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
    updateFn     func(ctx context.Context, user *domain.User) error
    deleteFn     func(ctx context.Context, id uuid.UUID) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
    return m.createFn(ctx, user)
}
// ... demais métodos delegam para o function field
```

No teste, só configura o que o cenário precisa:

```go
repo := &mockUserRepository{
    getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
        return nil, domain.ErrUserNotFound  // simula usuário não existente
    },
    createFn: func(_ context.Context, _ *domain.User) error {
        return nil  // simula persistência com sucesso
    },
}
svc := NewUserService(repo)
user, err := svc.Register(ctx, "João", "joao@email.com", "Str0ng!Pass")
```

### Testes de Service com Cenários de Erro

O `UserService` é testado cobrindo sucesso e **todos os caminhos de erro**:

```go
// internal/service/user_service_test.go

func TestUserService_Register(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        inputPass string
        mockRepo func() *mockUserRepository
        wantErr  error
    }{
        {"deve registrar usuário com sucesso",    "Str0ng!Pass",  mockOk(),     nil},
        {"deve retornar erro para senha fraca",    "weak",         mockEmpty(),  domain.ErrPasswordTooShort},
        {"deve retornar erro para senha sem especial", "Str0ngPass", mockEmpty(), domain.ErrPasswordTooWeak},
        {"deve retornar erro para email já em uso", "Str0ng!Pass", mockEmailExists(), domain.ErrEmailInUse},
    }
    // ...
}
```

Há teste específico que **valida que o hash gerado é bcrypt válido**:

```go
func TestUserService_Register_GeneratesBcryptHash(t *testing.T) {
    // ... registra usuário, captura o user salvo no mock
    if !strings.HasPrefix(savedUser.PasswordHash, "$2") {
        t.Error("PasswordHash não é bcrypt")
    }
}
```

### Testes de Validação de Senha

A validação de força da senha tem cobertura exaustiva com table-driven tests:

```go
// internal/service/password_validation_test.go

func TestValidatePasswordStrength(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name     string
        password string
        wantErr  error
    }{
        {"deve aceitar senha forte",                   "Str0ng!Pass",  nil},
        {"deve retornar erro para senha curta",        "Aa1@bc",       domain.ErrPasswordTooShort},
        {"deve retornar erro para senha sem maiúscula", "aaa1@bcde",   domain.ErrPasswordTooWeak},
        {"deve retornar erro para senha sem minúscula", "AAA1@BCDE",   domain.ErrPasswordTooWeak},
        {"deve retornar erro para senha sem número",    "Aaaa@bcde",   domain.ErrPasswordTooWeak},
        {"deve retornar erro para senha sem especial",  "Aaaa1bcde",   domain.ErrPasswordTooWeak},
        {"deve retornar erro para senha vazia",         "",            domain.ErrPasswordTooShort},
        {"deve retornar erro para apenas números",      "12345678",    domain.ErrPasswordTooWeak},
    }
    // ...
}
```

### Testes de JWT (Geração, Expiração, Assinatura)

O `auth/jwt_test.go` cobre os cenários críticos de segurança do token:

```go
// internal/auth/jwt_test.go

func TestJWTManager_GenerateAndVerify(t *testing.T) {
    // Gera token → Verifica → Confere claims (UserID, Email, Subject, ExpiresAt)
}

func TestJWTManager_Verify_TokenExpirado(t *testing.T) {
    manager := NewJWTManager("secret", -1*time.Hour)  // duração negativa = já expirado
    token, _ := manager.Generate(uuid.New(), "test@email.com")
    _, err := manager.Verify(token)
    // err != nil ✅
}

func TestJWTManager_Verify_SecretDiferente(t *testing.T) {
    manager1 := NewJWTManager("secret-1", 1*time.Hour)
    manager2 := NewJWTManager("secret-2", 1*time.Hour)
    token, _ := manager1.Generate(uuid.New(), "test@email.com")
    _, err := manager2.Verify(token)
    // err != nil ✅ — assinatura não bate
}

func TestJWTManager_Verify_TokenInvalido(t *testing.T) {
    // Testa: "not-a-jwt-token", string vazia, token manipulado
}
```

### Testes de Middleware de Autenticação

O middleware é testado com `httptest.NewRequest` e `httptest.NewRecorder`, cobrindo todos os cenários de falha:

```go
// internal/handler/middleware/auth_test.go

func TestAuth_Sucesso(t *testing.T) {
    // Gera token válido → envia no header → verifica que o handler recebeu o userID
    jwtManager := auth.NewJWTManager("test-secret", 1*time.Hour)
    token, _ := jwtManager.Generate(userID, "test@email.com")

    req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)
    // w.Code == 200, capturedUserID == userID ✅
}

func TestAuth_SemHeader(t *testing.T)       { /* → 401 "missing authorization header" */ }
func TestAuth_HeaderInvalido(t *testing.T)   { /* "sem Bearer", "Basic abc", "Bearer" sem token → 401 */ }
func TestAuth_TokenExpirado(t *testing.T)    { /* → 401 "invalid or expired token" */ }
func TestAuth_JWTManagerNil(t *testing.T)    { /* → 500 "jwt manager not configured" */ }
func TestUserIDFromContext_SemValor(t *testing.T) { /* contexto vazio → false */ }
```

### Testes de Helpers HTTP

Os helpers de decode e parsing são testados isoladamente:

```go
// internal/handler/response_test.go

func TestDecodeJSON(t *testing.T) {
    tests := []struct {
        name    string
        body    string
        wantErr bool
    }{
        {"JSON válido",         `{"name":"teste"}`,                          false},
        {"JSON inválido",       `{invalid}`,                                 true},
        {"campo desconhecido",  `{"name":"teste","unknown_field":"value"}`,   true},  // DisallowUnknownFields!
    }
    // ...
}

func TestParseUUID(t *testing.T) { /* válido, inválido, vazio */ }
func TestParseDecimal(t *testing.T) { /* "123.45", "100", "abc" */ }
func TestParseDate(t *testing.T) { /* RFC3339, yyyy-mm-dd, vazio, inválido */ }
```

### Testes de Config com `t.Setenv()`

Os testes de configuração usam **`t.Setenv()`** para simular variáveis de ambiente de forma segura (restaura automaticamente no cleanup):

```go
// internal/config/config_test.go

func TestLoad_Sucesso(t *testing.T) {
    t.Setenv("DATABASE_URL", "postgres://localhost/mydb")
    t.Setenv("JWT_SECRET", "super-secret")
    t.Setenv("PORT", "8080")

    cfg, err := Load()
    // cfg.DatabaseURL == "postgres://localhost/mydb" ✅
}

func TestLoad_SemDatabaseURL(t *testing.T) {
    t.Setenv("DATABASE_URL", "")
    t.Setenv("JWT_SECRET", "super-secret")
    _, err := Load()
    // err != nil ✅ — variável obrigatória ausente
}

func TestLoad_PortaPadrao(t *testing.T) {
    // PORT vazio → default "4235" ✅
}
```

### Testes de Domain (Regras de Negócio)

As entidades de domínio são testadas cobrindo **constructors, métodos de negócio, validações e tipos enumerados**:

```go
// internal/domain/bank_account_test.go

func TestAccountType_IsValid(t *testing.T) {
    // "checking" → true, "savings" → true, "" → false, "investment" → false
}

func TestBankAccount_ApplyIncome(t *testing.T) {
    // Somar R$500 a R$1000 → R$1500 ✅
    // Valor zero → ErrInvalidAmount ✅
    // Valor negativo → ErrInvalidAmount ✅
}

func TestBankAccount_ApplyExpense(t *testing.T) {
    // Checking com saldo insuficiente → permite overdraft (saldo negativo) ✅
    // Savings com saldo insuficiente → ErrNegativeBalance ✅
}

func TestBankAccount_AllowsOverdraft(t *testing.T) {
    // Checking → true, Savings → false
}
```

### Estratégia de Testes por Camada

| Camada | O que é testado | Técnicas usadas | Cobertura alvo |
|--------|-----------------|-----------------|----------------|
| `domain/` | Constructors, validações, métodos de negócio, `IsValid()`, `json:"-"` | Table-driven, `t.Parallel()`, `errors.Is()` | ≥ 80% |
| `service/` | Lógica de negócio, cenários de sucesso e erro | Mocks manuais com function fields, table-driven | ≥ 70% |
| `auth/` | Generate, Verify, expiração, assinatura errada | Table-driven, duração negativa para expiração | ≥ 80% |
| `config/` | Env vars presentes, ausentes, defaults | `t.Setenv()` | ≥ 60% |
| `handler/` | Helpers (decodeJSON, parseUUID, parseDecimal, parseDate) | `httptest`, table-driven | Helpers testados |
| `middleware/` | Auth (sem header, formato inválido, expirado, válido, nil manager) | `httptest.NewRequest`, `httptest.NewRecorder` | ≥ 70% |

### Filosofia de Teste

- **Stdlib only** — sem testify, gomega ou libs de assertion
- **Mocks manuais** — sem mockgen, mockery; structs com function fields
- **Table-driven** — todo cenário é um caso na tabela, não um teste separado
- **Paralelos** — `t.Parallel()` em todo teste e subtest para detectar race conditions
- **`errors.Is()`** — comparação de erros sentinela, não strings
- **`httptest`** — handlers e middleware testados sem subir servidor real
- **Assertions descritivas** — `t.Errorf("got %v, want %v", got, want)` com contexto claro
- **Race detector** — `go test -race` nos testes para garantir segurança em concorrência

---

## 🌐 Endpoints da API

### Autenticação (público)

| Método | Rota | Descrição |
|--------|------|-----------|
| `POST` | `/auth/register` | Registro de usuário |
| `POST` | `/auth/login` | Login (retorna JWT) |

### Rotas Protegidas (`/api/*` — requer `Authorization: Bearer <token>`)

| Método | Rota | Descrição |
|--------|------|-----------|
| `PUT` | `/api/me/password` | Troca de senha |
| | | |
| `POST` | `/api/accounts` | Criar conta bancária |
| `GET` | `/api/accounts` | Listar contas do usuário |
| `GET` | `/api/accounts/{id}` | Buscar conta por ID |
| `PUT` | `/api/accounts/{id}` | Atualizar conta |
| `DELETE` | `/api/accounts/{id}` | Remover conta |
| | | |
| `POST` | `/api/categories` | Criar categoria |
| `GET` | `/api/categories` | Listar categorias do usuário |
| `GET` | `/api/categories/{id}` | Buscar categoria por ID |
| `PUT` | `/api/categories/{id}` | Atualizar categoria |
| `DELETE` | `/api/categories/{id}` | Remover categoria |
| | | |
| `POST` | `/api/credit-cards` | Criar cartão de crédito |
| `GET` | `/api/credit-cards` | Listar cartões do usuário |
| `GET` | `/api/credit-cards/{id}` | Buscar cartão por ID |
| `PUT` | `/api/credit-cards/{id}` | Atualizar cartão |
| `DELETE` | `/api/credit-cards/{id}` | Remover cartão |
| | | |
| `POST` | `/api/transactions` | Criar transação manual |
| `GET` | `/api/transactions/{id}` | Buscar transação por ID |
| `PUT` | `/api/transactions/{id}` | Atualizar transação |
| `DELETE` | `/api/transactions/{id}` | Remover transação (reverte saldo) |
| `GET` | `/api/accounts/{accountID}/transactions` | Listar transações da conta |
| `POST` | `/api/transactions/planned-income/{id}` | Gerar transação de receita planejada |
| `POST` | `/api/transactions/planned-expense/{id}` | Gerar transação de despesa planejada |
| `POST` | `/api/transactions/pay-invoice` | Pagar fatura via conta bancária |
| | | |
| `POST` | `/api/credit-cards/{cardID}/transactions` | Criar compra no cartão |
| `GET` | `/api/credit-cards/{cardID}/transactions` | Listar transações do cartão |
| `GET` | `/api/credit-card-transactions/{id}` | Buscar transação do cartão por ID |
| `PUT` | `/api/credit-card-transactions/{id}` | Atualizar transação do cartão |
| `DELETE` | `/api/credit-card-transactions/{id}` | Remover transação do cartão |
| | | |
| `POST` | `/api/credit-cards/{cardID}/invoices/close-month` | Fechar fatura do mês |
| `GET` | `/api/credit-cards/{cardID}/invoices` | Listar faturas do cartão |
| `GET` | `/api/invoices/{id}` | Buscar fatura por ID |
| `DELETE` | `/api/invoices/{id}` | Remover fatura |
| | | |
| `POST` | `/api/planned-incomes` | Criar receita planejada |
| `GET` | `/api/planned-incomes` | Listar receitas planejadas |
| `GET` | `/api/planned-incomes/{id}` | Buscar receita planejada por ID |
| `PUT` | `/api/planned-incomes/{id}` | Atualizar receita planejada |
| `DELETE` | `/api/planned-incomes/{id}` | Remover receita planejada |
| | | |
| `POST` | `/api/planned-expenses` | Criar despesa planejada |
| `GET` | `/api/planned-expenses` | Listar despesas planejadas |
| `GET` | `/api/planned-expenses/{id}` | Buscar despesa planejada por ID |
| `PUT` | `/api/planned-expenses/{id}` | Atualizar despesa planejada |
| `DELETE` | `/api/planned-expenses/{id}` | Remover despesa planejada |

### Health Check

| Método | Rota | Descrição |
|--------|------|-----------|
| `GET` | `/health` | Status da API |

---

## 🚀 Rodando o Projeto

### Pré-requisitos

- Go 1.25+
- Docker e Docker Compose
- Make

### Setup rápido

```bash
# 1. Clone o repositório
git clone https://github.com/thenopholo/my-money.git
cd my-money

# 2. Setup inicial (instala ferramentas, copia .env)
make setup

# 3. Configure o .env com seus valores
# (edite o arquivo .env gerado a partir do .env.example)

# 4. Suba o banco de dados
make db-up

# 5. Rode as migrations
make migrate

# 6. Gere o código do sqlc
make sqlc

# 7. Rode a aplicação
make dev   # com hot-reload (air)
# ou
make run   # sem hot-reload
```

### Com Docker (produção)

```bash
# Build e start completo
make docker-build
make docker-run

# Parar tudo
make docker-stop
```

### Comandos disponíveis

```bash
make help  # Lista todos os comandos disponíveis
```

| Comando | Descrição |
|---------|-----------|
| `make dev` | Roda com hot-reload (air) |
| `make run` | Roda a aplicação |
| `make build` | Compila o binário |
| `make test` | Roda os testes |
| `make test-cover` | Testes com relatório de cobertura |
| `make lint` | Roda o linter |
| `make db-up` | Inicia o PostgreSQL |
| `make db-down` | Para o PostgreSQL |
| `make db-shell` | Acessa o psql |
| `make migrate` | Roda migrations pendentes |
| `make migrate-down` | Reverte última migration |
| `make migrate-new name=xxx` | Cria nova migration |
| `make sqlc` | Gera código Go das queries |
| `make sqlc-check` | Verifica se o código gerado está atualizado |

---

## 🔐 Variáveis de Ambiente

| Variável | Obrigatória | Descrição |
|----------|:-----------:|-----------|
| `DATABASE_URL` | ✅ | URL de conexão do PostgreSQL |
| `JWT_SECRET` | ✅ | Secret para assinatura dos tokens JWT |
| `PORT` | ❌ | Porta da API (default: `4235`) |
| `CORS_ALLOWED_ORIGINS` | ❌ | Origens permitidas, separadas por vírgula (default: `http://localhost:3000`) |
| `DB_USER` | ✅ | Usuário do PostgreSQL (usado pelo docker-compose) |
| `DB_PASSWORD` | ✅ | Senha do PostgreSQL (usado pelo docker-compose) |
| `DB_NAME` | ✅ | Nome do banco (usado pelo docker-compose) |
| `DB_PORT` | ✅ | Porta do PostgreSQL (usado pelo docker-compose) |

> ⚠️ O arquivo `.env` **nunca** é commitado. Use o `.env.example` como referência.

---

## 📝 Licença

Este projeto é de uso pessoal e para fins de portfólio.

---

<p align="center">
  Feito com Go, café e a vontade de não precisar de planilha pra saber se sobrou dinheiro no mês. ☕
</p>
