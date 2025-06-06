package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
	"github.com/fahrulrzi/score-match-backend/internal/repository"
	"github.com/fahrulrzi/score-match-backend/pkg/jwt"
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
	tokenRepo  repository.TokenRepository
}

func NewAuthMiddleware(jwtService *jwt.JWTService, tokenRepo repository.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		tokenRepo:  tokenRepo,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeErrorJSON(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			writeErrorJSON(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := bearerToken[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			writeErrorJSON(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Cek blacklist token
		isBlacklisted, err := m.tokenRepo.Exists(r.Context(), token, entity.Blacklisted)
		if err != nil {
			writeErrorJSON(w, http.StatusUnauthorized, "Authentication failed")
			return
		}
		if isBlacklisted {
			writeErrorJSON(w, http.StatusUnauthorized, "Token has been revoked")
			return
		}

		// Tambah user_id ke context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeErrorJSON(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"code":    statusCode,
		"message": message,
	})
}
