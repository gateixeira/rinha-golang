package main

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func (s *Storage) SaveTransaction(ctx context.Context, clientId int, t TransactionPayload) (*TransactionResponse, error) {
	var balance, limit int

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "SELECT balance, acc_limit FROM client WHERE id = $1 FOR UPDATE", clientId).Scan(&balance, &limit)
	if err != nil {
		return nil, errors.New("client not found")
	}

	op := 1
	if t.Type == "d" {
		op = -1
	}

	newBalance := balance + (t.Amount * op)
	if newBalance+limit < 0 {
		return nil, errors.New("limit exceeded")
	}

	tx.Exec(ctx, "INSERT INTO transaction (client_id, amount, op, transaction_description) VALUES ($1, $2, $3, $4)", clientId, t.Amount, t.Type, t.Description)

	tx.Exec(ctx, "UPDATE client SET balance = $1 WHERE id = $2", newBalance, clientId)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &TransactionResponse{
		Limit:   limit,
		Balance: newBalance,
	}, nil
}

func (s *Storage) GetStatement(ctx context.Context, clientId int) (*Statement, error) {
	var balance, limit int
	var transactions []Transaction

	err := s.DB.QueryRow(ctx, "SELECT balance, acc_limit FROM client WHERE id = $1", clientId).Scan(&balance, &limit)
	if err != nil {
		return nil, errors.New("client not found")
	}

	rows, err := s.DB.Query(ctx, "SELECT amount, op, transaction_description, completed_at FROM transaction WHERE client_id = $1 ORDER BY id DESC LIMIT 10", clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.Amount, &t.Op, &t.Description, &t.CompletedAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return &Statement{
		Transactions: transactions,
		Balance: Balance{
			Total: balance,
			Date:  time.Now(),
			Limit: limit,
		},
	}, nil
}
