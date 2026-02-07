package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/thenopholo/my-money/internal/auth"
	"github.com/thenopholo/my-money/internal/handler"
	handlermw "github.com/thenopholo/my-money/internal/handler/middleware"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}
	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)

	// TODO: conectar banco, criar repositories/services e instanciar handlers reais.
	router := handler.NewRouter(
		nil, // userHandler
		nil, // bankAccountHandler
		nil, // categoryHandler
		nil, // creditCardHandler
		nil, // transactionHandler
		nil, // creditCardTransactionHandler
		nil, // invoiceHandler
		handlermw.Auth(jwtManager),
	)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
