package handler

import (
	"authn_authz/helpers"
	"authn_authz/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type UserHandler struct {
	Us       *service.UserService
	Validate *validator.Validate
}

func (u *UserHandler) RetrieveFiles(w http.ResponseWriter, r *http.Request) {
	// retrieve user id from context
	userId, ok := r.Context().Value(helpers.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "user id not found in the request context", http.StatusUnauthorized)
		return
	}

	// check the role from context
	role, ok := r.Context().Value(helpers.RoleKey).(string)
	if !ok {
		http.Error(w, "role not found in the request context", http.StatusUnauthorized)
		return
	}

	if role == "admin" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u.Us.Data)
		return
	}

	// call the service
	file, err := u.Us.RetrieveFiles(userId)
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
	json.NewEncoder(w).Encode(file)
}

func (u *UserHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	// decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	// Validate
	err = u.Validate.Var(req.Name, "required,min=1,max=6")
	if err != nil {
		http.Error(w, "invalid file name", http.StatusBadRequest)
		return
	}

	// retrieve user id from context
	userId, ok := r.Context().Value(helpers.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "user id not found in the request context", http.StatusUnauthorized)
		return
	}

	// call the service
	file, err := u.Us.UploadFile(userId, req.Name)
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
	json.NewEncoder(w).Encode(file)
}
