package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thenopholo/my-money/internal/auth"
	"github.com/thenopholo/my-money/internal/config"
	"github.com/thenopholo/my-money/internal/handler"
	handlermw "github.com/thenopholo/my-money/internal/handler/middleware"
	"github.com/thenopholo/my-money/internal/repository"
	"github.com/thenopholo/my-money/internal/repository/postgres"
	"github.com/thenopholo/my-money/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping database:", err)
	}
	log.Println("Connected to database")

	queries := postgres.New(pool)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTDuration)

	userRepo := repository.NewUserRepository(queries)
	bankAccountRepo := repository.NewBankAccountRepository(queries)
	categoryRepo := repository.NewCategoryRepository(queries)
	creditCardRepo := repository.NewCreditCardRepository(queries)
	invoiceRepo := repository.NewInvoiceRepository(queries)
	transactionRepo := repository.NewTransactionRepository(queries)
	ccTxRepo := repository.NewCreditCardTransactionRepository(queries)
	plannedIncomeRepo := repository.NewPlannedIncomeRepository(queries)
	plannedExpenseRepo := repository.NewPlannedExpenseRepository(queries)

	userService := service.NewUserService(userRepo)
	bankAccountService := service.NewBankAccountService(bankAccountRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	creditCardService := service.NewCreditCardService(creditCardRepo)
	plannedIncomeService := service.NewPlannedIncomeService(plannedIncomeRepo)
	plannedExpenseService := service.NewPlannedExpenseService(plannedExpenseRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo, creditCardRepo, ccTxRepo)
	transactionService := service.NewTransactionService(
		transactionRepo,
		bankAccountRepo,
		categoryRepo,
		plannedIncomeRepo,
		plannedExpenseRepo,
		invoiceRepo,
	)
	ccTxService := service.NewCreditCardTransactionService(ccTxRepo, creditCardRepo, categoryRepo, invoiceRepo)

	userHandler := handler.NewUserHandler(userService, jwtManager)
	bankAccountHandler := handler.NewBankAccountHandler(bankAccountService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	creditCardHandler := handler.NewCreditCardHandler(creditCardService)
	plannedIncomeHandler := handler.NewPlannedIncomeHandler(plannedIncomeService)
	plannedExpenseHandler := handler.NewPlannedExpenseHandler(plannedExpenseService)
	transactionHandler := handler.NewTransactionHandler(transactionService, bankAccountService)
	ccTxHandler := handler.NewCreditCardTransactionHandler(ccTxService, creditCardService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService, creditCardService)

	router := handler.NewRouter(
		userHandler,
		bankAccountHandler,
		categoryHandler,
		creditCardHandler,
		plannedIncomeHandler,
		plannedExpenseHandler,
		transactionHandler,
		ccTxHandler,
		invoiceHandler,
		handlermw.Auth(jwtManager),
	)

	log.Printf("Server running on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
