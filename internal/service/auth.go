package service

import (
	"errors"
	"fmt"

	"github.com/fatihrizqon/go-fiber-service/internal/delivery/http/request"
	"github.com/fatihrizqon/go-fiber-service/internal/delivery/http/response"
	"github.com/fatihrizqon/go-fiber-service/internal/repository"
	"github.com/fatihrizqon/go-fiber-service/internal/util"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
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

	token, err := util.CreateToken(result)
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

func (e *AuthService) Verify(req *request.VerifyUserRequest) (response.AuthResponse, error) {
	var res response.AuthResponse

	claims, err := util.ParseToken(req.Token)
	if err != nil {
		return res, fmt.Errorf("failed to parse token: %w", err)
	}

	user, err := e.IAuthRepository.Login(claims.Email)
	if err != nil {
		return res, fmt.Errorf("failed to get user: %w", err)
	}

	return response.AuthResponse{
		ID: user.Id,
	}, nil
}
