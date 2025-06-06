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

func writeJSON(w http.ResponseWriter, statusCode int, status string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"code":    statusCode,
		"message": message,
		"data":    data,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entity.UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		writeJSON(w, http.StatusBadRequest, "error", "Please provide all required fields", nil)
		return
	}

	user, token, err := h.authUseCase.Register(r.Context(), &req)
	if err != nil {
		writeJSON(w, http.StatusConflict, "error", err.Error(), nil)
		return
	}

	resp := entity.RegisterResponse{
		Token: token,
		User:  *user,
	}
	writeJSON(w, http.StatusCreated, "success", "User registered successfully", resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, "error", "Email and password are required", nil)
		return
	}

	user, token, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, "error", err.Error(), nil)
		return
	}

	resp := entity.LoginResponse{
		Token: token,
		User:  *user,
	}
	writeJSON(w, http.StatusOK, "success", "Login successful", resp)
}
