package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
	"github.com/fahrulrzi/score-match-backend/internal/usecase"
	"github.com/fahrulrzi/score-match-backend/pkg/jwt"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	jwtService  *jwt.JWTService
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, jwtService *jwt.JWTService) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entity.UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "Please provide all required fields", http.StatusBadRequest)
		return
	}

	user, token, err := h.authUseCase.Register(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.RegisterResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.LoginResponse{
		Token: token,
		User:  *user,
	})
}
