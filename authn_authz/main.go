package main

import (
	"authn_authz/db"
	"authn_authz/handler"
	"authn_authz/middleware"
	"authn_authz/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initialize
	user := db.Users
	file := db.Files

	validate := validator.New()

	// auth
	authService := service.AuthService{
		Data: user,
	}
	authHandler := handler.AuthHandler{
		As:       &authService,
		Validate: validate,
	}

	userService := service.UserService{
		Data: file,
	}
	userHandler := handler.UserHandler{
		Us:       &userService,
		Validate: validate,
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "server is up",
		})
	})
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/files", middleware.Authenticate(middleware.RequireRole("admin", "user")(userHandler.RetrieveFiles))).Methods("GET")
	router.HandleFunc("/upload", middleware.Authenticate(userHandler.UploadFile)).Methods("POST")

	http.ListenAndServe(":8080", router)
}
