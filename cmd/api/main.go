package main

import (
	"log"
	"net/http"

	"github.com/thenopholo/my-money/internal/handler"
)

func main() {
	// TODO: conectar banco, criar repositories/services e instanciar handlers reais.
	router := handler.NewRouter(
		nil, // userHandler
		nil, // bankAccountHandler
		nil, // categoryHandler
		nil, // creditCardHandler
		nil, // transactionHandler
		nil, // creditCardTransactionHandler
		nil, // invoiceHandler
	)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
