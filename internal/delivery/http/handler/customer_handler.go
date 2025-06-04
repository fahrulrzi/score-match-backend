package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
	"github.com/fahrulrzi/score-match-backend/internal/usecase"
)

type CustomerHandler struct {
	customerUseCase usecase.CustomerUseCase
}

func NewCustomerHandler(customerUseCase usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
	}
}

func (c *CustomerHandler) GetFinalScore(w http.ResponseWriter, r *http.Request) {
	var req entity.CustomerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := c.customerUseCase.GetFinalScore(r.Context(), &req)

	err = c.customerUseCase.Create(r.Context(), customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res entity.CustomerScoreResponse
	res.Score = customer.Score
	res.Status = customer.Status
	res.Describe = customer.Describe

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
