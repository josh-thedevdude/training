package repository

import (
	"context"
	"database/sql"
	"errors"
	"personal_expense_tracker/domain"
)

type ExpenseRepository interface {
	Add(ctx context.Context, newTransaction *domain.Transaction) (domain.Transaction, error)
	GetById(ctx context.Context, id string) (domain.Transaction, error)
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	UpdateById(ctx context.Context, id string, description string) (domain.Transaction, error)
	DeleteById(ctx context.Context, id string) (domain.Transaction, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{
		db: db,
	}
}

func (e *expenseRepository) Add(ctx context.Context, txn *domain.Transaction) (domain.Transaction, error) {
	query := `
		INSERT INTO expenses (id, amount, description)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at;
	`
	err := e.db.QueryRowContext(ctx, query, txn.ID, txn.Amount, txn.Description).
		Scan(&txn.CreatedAt, &txn.UpdatedAt)
	if err != nil {
		return domain.Transaction{}, err
	}
	return *txn, nil
}

func (e *expenseRepository) GetById(ctx context.Context, id string) (domain.Transaction, error) {
	query := `
		SELECT id, amount, description, created_at, updated_at
		FROM expenses
		WHERE id = $1;
	`
	var tx domain.Transaction
	err := e.db.QueryRowContext(ctx, query, id).
		Scan(&tx.ID, &tx.Amount, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Transaction{}, errors.New("invalid transaction id")
		}
		return domain.Transaction{}, err
	}
	return tx, nil
}

func (e *expenseRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	query := `
		SELECT id, amount, description, created_at, updated_at
		FROM expenses
		ORDER BY created_at DESC;
	`
	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []domain.Transaction

	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(&tx.ID, &tx.Amount, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (e *expenseRepository) UpdateById(ctx context.Context, id string, description string) (domain.Transaction, error) {
	query := `
		UPDATE expenses
		SET description = $2,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, amount, description, created_at, updated_at;
	`
	var tx domain.Transaction
	err := e.db.QueryRowContext(ctx, query, id, description).
		Scan(&tx.ID, &tx.Amount, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Transaction{}, errors.New("transaction not found")
		}
		return domain.Transaction{}, err
	}
	return tx, nil
}

func (e *expenseRepository) DeleteById(ctx context.Context, id string) (domain.Transaction, error) {
	query := `
		DELETE FROM expenses
		WHERE id = $1
		RETURNING id, amount, description, created_at, updated_at;
	`
	var tx domain.Transaction
	err := e.db.QueryRowContext(ctx, query, id).
		Scan(&tx.ID, &tx.Amount, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Transaction{}, errors.New("transaction not found")
		}
		return domain.Transaction{}, err
	}
	return tx, nil
}
