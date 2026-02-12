package handler

import (
	"authn_authz/helpers"
	"authn_authz/service"
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type AuthHandler struct {
	As       *service.AuthService
	Validate *validator.Validate
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role     string `json:"role"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// parse to json
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	// Validate
	err = a.Validate.Var(req.Email, "required,email")
	if err != nil {
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return
	}
	err = a.Validate.Var(req.Password, "required,min=8,max=15")
	if err != nil {
		http.Error(w, "invalid password format", http.StatusBadRequest)
		return
	}
	err = a.Validate.Var(req.Role, "required,oneof=user admin")
	if err != nil {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	// call the service
	user, err := a.As.Register(req.Email, req.Password, req.Role)
	if err != nil {
		if appErr, ok := err.(*helpers.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// parse to json
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	// Validate
	err = a.Validate.Var(req.Email, "required,email")
	if err != nil {
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return
	}
	err = a.Validate.Var(req.Password, "required,min=8,max=15")
	if err != nil {
		http.Error(w, "invalid password format", http.StatusBadRequest)
		return
	}

	// call the service
	user, err := a.As.Login(req.Email, req.Password)
	if err != nil {
		if appErr, ok := err.(*helpers.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
