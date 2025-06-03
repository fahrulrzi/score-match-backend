package usecase

import (
	"context"
	"errors"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
	"github.com/fahrulrzi/score-match-backend/internal/repository"
	"github.com/fahrulrzi/score-match-backend/pkg/hash"
	"github.com/fahrulrzi/score-match-backend/pkg/jwt"
)

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtService *jwt.JWTService
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	jwtService *jwt.JWTService,
) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register implements AuthUseCase.
func (a *authUseCase) Register(ctx context.Context, userRegister *entity.UserRegisterRequest) (*entity.User, string, error) {
	existingUser, err := a.userRepo.GetUserByEmail(ctx, userRegister.Email)
	if err != nil {
		return nil, "", err
	}

	if existingUser != nil {
		return nil, "", errors.New("user with this email already exists")
	}

	hashedPassword, err := hash.HashPassword(userRegister.Password)
	if err != nil {
		return nil, "", err
	}

	user := &entity.User{
		Email:    userRegister.Email,
		Password: hashedPassword,
		Username: userRegister.Username,
	}

	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login implements AuthUseCase.
func (a *authUseCase) Login(ctx context.Context, userLogin *entity.UserLoginRequest) (*entity.User, string, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, userLogin.Email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errors.New("Invalid email or password")
	}

	if !hash.CheckPasswordHash(userLogin.Password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
