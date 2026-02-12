package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(userId uuid.UUID, role string) (string, error) {
	claims := CustomClaims{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", &AppError{
			Code:    500,
			Message: "failed to generate token",
			Tag:     "jwt_signing_error",
		}
	}
	return ss, nil
}

func VerifyJwtToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil {
		return nil, &AppError{
			Code:    500,
			Message: "failed to parse jwt claims",
			Tag:     "jwt_parse_error",
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, &AppError{
		Code:    500,
		Message: "unknown claims type",
		Tag:     "jwt_unknown_claims",
	}
}
