package service

import (
	"authn_authz/db"
	"authn_authz/helpers"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Data []*db.User
}

func (a *AuthService) Register(email, password, role string) (*db.User, error) {
	// check if user already exists
	for _, u := range a.Data {
		if u.Email == email {
			return nil, &helpers.AppError{
				Code:    400,
				Message: "user already exists",
				Tag:     "duplicate_user",
			}
		}
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &helpers.AppError{
			Code:    500,
			Message: "failed to hash password",
			Tag:     "hash_error",
		}
	}

	// create a new user
	newUser := db.User{
		ID:       uuid.New(),
		Email:    email,
		Role:     role,
		Password: string(passwordHash),
	}
	a.Data = append(a.Data, &newUser)

	// reponse
	return &newUser, nil
}

func (a *AuthService) Login(email, password string) (string, error) {
	// check if user exists
	var user *db.User
	for i := range a.Data {
		if a.Data[i].Email == email {
			user = a.Data[i]
		}
	}
	if user == nil {
		return "", &helpers.AppError{
			Code:    404,
			Message: fmt.Sprintf("user with email %s not found", email),
			Tag:     "user_not_found",
		}
	}

	// check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", &helpers.AppError{
			Code:    401,
			Message: "authentication failed",
			Tag:     "invalid_password",
		}
	}

	// generate token
	token, err := helpers.GenerateJwtToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	// return token
	return token, nil
}
