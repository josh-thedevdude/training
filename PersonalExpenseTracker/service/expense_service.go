package service

import (
	"context"
	"personal_expense_tracker/domain"
	"personal_expense_tracker/repository"
)

type ExpenseService interface {
	AddTransaction(ctx context.Context, amount int, description string) (domain.Transaction, error)
	GetTransactionById(ctx context.Context, id string) (domain.Transaction, error)
	GetTransactions(ctx context.Context) ([]domain.Transaction, error)
	UpdateTransactionById(ctx context.Context, id string, description string) (domain.Transaction, error)
	DeleteTransactionById(ctx context.Context, id string) (domain.Transaction, error)
}

type expenseService struct {
	repo repository.ExpenseRepository
}

func NewExpenseService(repo repository.ExpenseRepository) ExpenseService {
	return &expenseService{
		repo: repo,
	}
}

func (e *expenseService) AddTransaction(ctx context.Context, amount int, description string) (domain.Transaction, error) {
	newTransaction := domain.NewTransaction(amount, description)
	newTxn, err := e.repo.Add(ctx, newTransaction)
	if err != nil {
		return domain.Transaction{}, err
	}
	return newTxn, nil
}

func (e *expenseService) GetTransactionById(ctx context.Context, id string) (domain.Transaction, error) {
	txn, err := e.repo.GetById(ctx, id)
	if err != nil {
		return domain.Transaction{}, err
	}
	return txn, nil
}

func (e *expenseService) GetTransactions(ctx context.Context) ([]domain.Transaction, error) {
	txns, err := e.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return txns, nil
}

func (e *expenseService) UpdateTransactionById(ctx context.Context, id string, description string) (domain.Transaction, error) {
	updatedTxn, err := e.repo.UpdateById(ctx, id, description)
	if err != nil {
		return domain.Transaction{}, err
	}
	return updatedTxn, nil
}

func (e *expenseService) DeleteTransactionById(ctx context.Context, id string) (domain.Transaction, error) {
	deletedTxn, err := e.repo.DeleteById(ctx, id)
	if err != nil {
		return domain.Transaction{}, err
	}
	return deletedTxn, nil
}
