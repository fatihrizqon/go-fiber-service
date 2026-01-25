package handler

import (
	"strings"
	"time"

	"github.com/fatihrizqon/go-fiber-service/helper"
	"github.com/fatihrizqon/go-fiber-service/internal/entity"
	"github.com/fatihrizqon/go-fiber-service/internal/presenter/request"
	"github.com/fatihrizqon/go-fiber-service/internal/presenter/response"
	"github.com/fatihrizqon/go-fiber-service/internal/service"
	"github.com/fatihrizqon/go-fiber-service/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	IAuthService service.IAuthService
}

func NewAuthHandler(serv service.IAuthService) *AuthHandler {
	return &AuthHandler{IAuthService: serv}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return a JWT token in a cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.AuthJSON
// @Failure 400 {object} response.JSON "Invalid request format"
// @Failure 401 {object} response.JSON "Authentication failed"
// @Router /api/v1/auth/login [post]
func (handler *AuthHandler) Login(ctx *fiber.Ctx) error {
	log := logger.GetLogger()
	ip := ctx.IP()

	var req request.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.WithField("ip", ip).Error("failed to parse login request")
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request format")
	}

	log.WithField("ip", ip).Info("user login attempt: " + req.Email)

	result, err := handler.IAuthService.Login(req)
	if err != nil {
		log.WithField("ip", ip).Error("authentication failed: " + req.Email)
		return errorResponse(ctx, fiber.StatusUnauthorized, "authentication failed: "+err.Error())
	}

	accessToken, _ := helper.GenerateAccessToken(result.User)
	refreshToken, _ := helper.GenerateRefreshToken(result.User)

	setAuthCookies(ctx, accessToken, refreshToken)

	log.WithField("ip", ip).Info("user logged in: " + req.Email)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Message: "you are authenticated",
		Status:  fiber.StatusOK,
		User: response.UserInfo{
			Id:              result.User.Id,
			Username:        result.User.Username,
			Name:            result.User.Name,
			Email:           result.User.Email,
			Status:          result.User.Status,
			EmailVerifiedAt: result.User.EmailVerifiedAt.Format(time.RFC3339),
		},
		AccessToken: accessToken,
	})
}

// Refresh Token godoc
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token from HttpOnly cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.JSON "Access token refreshed"
// @Failure 400 {object} response.JSON "Invalid request format"
// @Failure 401 {object} response.JSON "Invalid or missing refresh token"
// @Router /api/v1/auth/refresh [post]
func (handler *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return errorResponse(ctx, fiber.StatusUnauthorized, "refresh token required")
	}

	claims, err := helper.ParseToken(refreshToken, true)
	if err != nil {
		return errorResponse(ctx, fiber.StatusUnauthorized, "invalid refresh token")
	}

	userID, err := parseUserIDFromClaims(claims)
	if err != nil {
		return errorResponse(ctx, fiber.StatusInternalServerError, "invalid user ID")
	}

	username, _ := claims["username"].(string)
	name, _ := claims["name"].(string)

	accessToken, _ := helper.GenerateAccessToken(entity.User{
		Id:       userID,
		Username: username,
		Name:     name,
	})

	setAuthCookies(ctx, accessToken, refreshToken)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Access token refreshed",
	})
}

// Get User Info godoc
// @Summary Get authenticated user info
// @Description Retrieve the current authenticated user's information using the access token stored in HttpOnly cookie or Authorization header.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.AuthJSON "User info retrieved"
// @Failure 401 {object} response.JSON "Invalid or missing access token"
// @Router /api/v1/auth/me [get]
func (handler *AuthHandler) Me(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies("access_token")

	if accessToken == "" {
		authHeader := ctx.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if accessToken == "" {
		return errorResponse(ctx, fiber.StatusUnauthorized, "access token required")
	}

	claims, err := helper.ParseToken(accessToken, false)
	if err != nil {
		return errorResponse(ctx, fiber.StatusUnauthorized, "invalid access token")
	}

	userID, err := parseUserIDFromClaims(claims)
	if err != nil {
		return errorResponse(ctx, fiber.StatusUnauthorized, "invalid user ID")
	}

	username, _ := claims["username"].(string)
	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)

	statusFloat, _ := claims["status"].(float64)
	status := int(statusFloat)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Message: "user info retrieved",
		Status:  fiber.StatusOK,
		User: response.UserInfo{
			Id:       userID,
			Username: username,
			Name:     name,
			Email:    email,
			Status:   status,
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user dengan menghapus access_token dan refresh_token dari cookie,
// @Description serta memasukkan refresh token ke blacklist.
// @Tags Auth
// @Success 200 {object} response.JSON "Successfully logged out"
// @Failure 401 {object} response.JSON "Unauthorized"
// @Router /api/v1/auth/logout [post]
func (handler *AuthHandler) Logout(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	helper.BlacklistToken(refreshToken)

	clearAuthCookies(ctx)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "successfully logged out",
	})
}

// setAuthCookies sets access & refresh tokens as HttpOnly cookies
func setAuthCookies(ctx *fiber.Ctx, accessToken, refreshToken string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}

func clearAuthCookies(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Expires: time.Now().Add(-1 * time.Hour),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Expires: time.Now().Add(-1 * time.Hour),
	})
}

func parseUserIDFromClaims(claims map[string]interface{}) (uuid.UUID, error) {
	idStr, ok := claims["id"].(string)
	if !ok {
		return uuid.UUID{}, fiber.ErrInternalServerError
	}
	return uuid.Parse(idStr)
}

func errorResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(response.JSON{
		Status:  status,
		Message: message,
	})
}
