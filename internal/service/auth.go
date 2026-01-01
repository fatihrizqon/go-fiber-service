package service

import (
	"errors"
	"fmt"

	"github.com/fatihrizqon/go-fiber-service/helper"
	"github.com/fatihrizqon/go-fiber-service/internal/presenter/request"
	"github.com/fatihrizqon/go-fiber-service/internal/presenter/response"
	"github.com/fatihrizqon/go-fiber-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(req request.RegisterRequest) (response.RegisterResponse, error)
	Login(req request.LoginRequest) (response.LoginResponse, error)
}
type AuthService struct {
	IAuthRepository repository.IAuthRepository
	validate        *validator.Validate
}

func NewAuthService(repo repository.IAuthRepository, validate *validator.Validate) IAuthService {
	return &AuthService{
		IAuthRepository: repo,
		validate:        validate,
	}
}

// Register implements IAuthService.
func (e *AuthService) Register(req request.RegisterRequest) (response.RegisterResponse, error) {
	panic("unimplemented")
}

// Login implements IAuthService.
func (e *AuthService) Login(req request.LoginRequest) (response.LoginResponse, error) {
	var res response.LoginResponse

	result, err := e.IAuthRepository.Login(req.Email)
	if err != nil {
		return res, err
	}

	err = ValidatePassword(req.Password, result.Password)
	if err != nil {
		return res, errors.New("credentials does not matches our record")
	}

	token, err := helper.GenerateRefreshToken(result)
	if err != nil {
		return res, fmt.Errorf("failed to generate token: %w", err)
	}

	return response.LoginResponse{
		Token: token,
		User:  result,
	}, nil
}

// ValidatePassword compares a plain password with a hashed password
func ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
