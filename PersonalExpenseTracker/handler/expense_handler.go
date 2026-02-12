package handler

import (
	"encoding/json"
	"net/http"
	"personal_expense_tracker/service"

	"github.com/gorilla/mux"
)

type ExpenseHandler struct {
	service service.ExpenseService
}

func NewExpenseHandler(service service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

type CreateExpenseRequest struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type UpdateExpenseRequest struct {
	Description string `json:"description"`
}

func (h *ExpenseHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	txn, err := h.service.AddTransaction(
		r.Context(),
		req.Amount,
		req.Description,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(txn)
}

func (h *ExpenseHandler) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	txn, err := h.service.GetTransactionById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(txn)
}

func (h *ExpenseHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	txns, err := h.service.GetTransactions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(txns)
}

func (h *ExpenseHandler) UpdateTransactionById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedTxn, err := h.service.UpdateTransactionById(
		r.Context(),
		id,
		req.Description,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedTxn)
}

func (h *ExpenseHandler) DeleteTransactionById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	deletedTxn, err := h.service.DeleteTransactionById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(deletedTxn)
}
