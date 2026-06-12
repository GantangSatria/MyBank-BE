package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/GantangSatria/MyBank-BE/internal/repository"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type AuthMiddleware struct {
	secret   string
	userRepo repository.UserRepository
}

func NewAuthMiddleware(secret string, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		secret:   secret,
		userRepo: userRepo,
	}
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (m *AuthMiddleware) Protect() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return apperrors.Unauthorized("missing token")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return apperrors.Unauthorized("invalid token format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.Unauthorized("invalid token method")
			}
			return []byte(m.secret), nil
		})

		if err != nil || !token.Valid {
			return apperrors.Unauthorized("invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

func (m *AuthMiddleware) RequirePIN() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := GetUserID(c)

		var body struct {
			PIN string `json:"pin"`
		}

		if err := c.Bind().JSON(&body); err != nil {
			return apperrors.BadRequest("PIN diperlukan")
		}

		if body.PIN == "" {
			return apperrors.BadRequest("PIN wajib diisi")
		}

		// ambil user dengan PIN
		user, err := m.userRepo.FindByIDWithPIN(c.Context(), userID)
		if err != nil {
			return apperrors.Unauthorized("user tidak valid")
		}

		// compare hash
		if err := bcrypt.CompareHashAndPassword([]byte(*user.PIN), []byte(body.PIN)); err != nil {
			return apperrors.ErrInvalidPIN
		}

		return c.Next()
	}
}

func GetUserID(c fiber.Ctx) uint64 {
	userID, ok := c.Locals("user_id").(uint64)
	if !ok {
		return 0
	}
	return userID
}