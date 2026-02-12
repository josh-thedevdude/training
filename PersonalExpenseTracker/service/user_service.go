package service

import (
	"context"
	"fmt"
	"personal_expense_tracker/domain"
	"personal_expense_tracker/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserSerivce interface {
	Register(ctx context.Context, email, password, role string) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository) UserSerivce {
	return &userService{
		repo: repo,
	}
}

type MyCustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (u *userService) Register(ctx context.Context, email, password, role string) (domain.User, error) {
	// hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	newUser := domain.NewUser(email, string(passwordHash), role)

	user, err := u.repo.Add(ctx, newUser)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userService) Login(ctx context.Context, email, password string) (string, error) {
	// check if user by email exists
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// generate jwt token
	claims := MyCustomClaims{
		user.Email,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// generate token string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
