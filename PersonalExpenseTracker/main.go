package main

import (
	"fmt"
	"net/http"
	"personal_expense_tracker/db"
	"personal_expense_tracker/handler"
	"personal_expense_tracker/repository"
	"personal_expense_tracker/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbConn, err := db.NewPgDatabase()
	if err != nil {
		panic(fmt.Errorf("Error connecting to database %s", err.Error()))
	}
	defer dbConn.Close()

	// Initialize layers
	expenseRepo := repository.NewExpenseRepository(dbConn)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	// Setup router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/expenses", expenseHandler.AddTransaction).Methods(http.MethodPost)
	r.HandleFunc("/expenses", expenseHandler.GetTransactions).Methods(http.MethodGet)
	r.HandleFunc("/expenses/{id}", expenseHandler.GetTransactionById).Methods(http.MethodGet)
	r.HandleFunc("/expenses/{id}", expenseHandler.UpdateTransactionById).Methods(http.MethodPatch)
	r.HandleFunc("/expenses/{id}", expenseHandler.DeleteTransactionById).Methods(http.MethodDelete)

	// Start server
	fmt.Println("🚀 Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(fmt.Errorf("❌ Server failed: %v", err))
	}
}
