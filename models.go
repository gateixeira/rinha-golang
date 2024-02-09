package main

import "time"

type TransactionPayload struct {
	Amount      int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type TransactionResponse struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type Transaction struct {
	Amount      int       `json:"valor"`
	Op          string    `json:"tipo"`
	Description string    `json:"descricao"`
	CompletedAt time.Time `json:"realizada_em"`
}

type Balance struct {
	Total int       `json:"total"`
	Date  time.Time `json:"data_extrato"`
	Limit int       `json:"limite"`
}

type Statement struct {
	Transactions []Transaction `json:"ultimas_transacoes"`
	Balance      Balance       `json:"saldo"`
}
