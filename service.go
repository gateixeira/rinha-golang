package main

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type Service struct {
	prefix string
	store  Storage
}

func validTransactionPayload(t TransactionPayload) bool {
	return !(t.Amount <= 0 || (t.Type != "d" && t.Type != "c") || (len(t.Description) < 1 || len(t.Description) > 10))
}

func setValidResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request, store *Storage, id int) {
	var t TransactionPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	if !validTransactionPayload(t) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	result, err := store.SaveTransaction(context.Background(), id, t)
	if err != nil {
		if err.Error() == "limit exceeded" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err.Error() == "client not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setValidResponse(w, result)
}

func getStatementHandler(w http.ResponseWriter, r *http.Request, store *Storage, id int) {
	result, err := store.GetStatement(context.Background(), id)
	if err != nil {
		if err.Error() == "client not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setValidResponse(w, result)
}

func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, svc.prefix)
	if url == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	idStr, op := path.Split(url)
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost && op == "transacoes" {
		createTransactionHandler(w, r, &svc.store, id)
		return
	}

	if r.Method == http.MethodGet && op == "extrato" {
		getStatementHandler(w, r, &svc.store, id)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
