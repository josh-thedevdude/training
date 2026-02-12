package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

var (
	// Google OAuth2 Configuration - loaded from environment variables
	googleClientID     string
	googleClientSecret string
	googleRedirectURI  string

	// Session store - uses secure cookies for session management
	// SECURITY: In production, use Redis or database-backed sessions for scalability
	store *sessions.CookieStore
)

const (
	sessionName     = "oauth-session"
	sessionUserKey  = "user"
	sessionStateKey = "oauth_state"
	sessionMaxAge   = 86400 // 24 hours in seconds

	// Google OAuth2 endpoints
	googleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// ============================================================================
// DATA STRUCTURES
// ============================================================================

// User represents an authenticated user in our system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
}

// GoogleUserInfo represents the user profile data from Google
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GoogleTokenResponse represents Google's OAuth token response
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// ============================================================================
// IN-MEMORY USER STORE
// ============================================================================
// NOTE: In production, replace this with a proper database (PostgreSQL, MySQL, etc.)

var users = make(map[string]*User)

// findOrCreateUser finds an existing user or creates a new one
func findOrCreateUser(profile *GoogleUserInfo) *User {
	if user, exists := users[profile.ID]; exists {
		log.Printf("Found existing user: %s (%s)", user.Name, user.Email)
		return user
	}

	user := &User{
		ID:        profile.ID,
		Email:     profile.Email,
		Name:      profile.Name,
		Picture:   profile.Picture,
		CreatedAt: time.Now(),
	}

	users[profile.ID] = user
	log.Printf("Created new user: %s (%s)", user.Name, user.Email)
	return user
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// generateSecureState creates a cryptographically secure random state value
// This is used for CSRF protection in the OAuth flow
func generateSecureState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// getSession retrieves or creates a session for the request
func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionName)
}

// ============================================================================
// OAUTH HANDLERS
// ============================================================================

/**
 * STEP 1: Start OAuth Flow
 * GET /auth/google/start
 *
 * This endpoint initiates the OAuth 2.0 Authorization Code Flow:
 * 1. Generates a cryptographically secure random "state" value for CSRF protection
 * 2. Stores the state in the server-side session
 * 3. Constructs the Google OAuth authorization URL with required parameters
 * 4. Redirects the user to Google's consent screen
 *
 * SECURITY:
 * - State parameter prevents CSRF attacks
 * - Server-side session storage ensures state cannot be tampered with
 * - Client secret is NEVER exposed to the frontend
 */
func handleGoogleStart(w http.ResponseWriter, r *http.Request) {
	// Get or create session
	session, err := getSession(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		log.Printf("Session error: %v", err)
		return
	}

	// Generate cryptographically secure random state
	state, err := generateSecureState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		log.Printf("State generation error: %v", err)
		return
	}

	// Store state in server-side session for validation on callback
	session.Values[sessionStateKey] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		log.Printf("Session save error: %v", err)
		return
	}

	log.Printf("🔐 Starting OAuth flow with state: %s", state[:10]+"...")

	// Build Google OAuth authorization URL
	params := url.Values{
		"client_id":     {googleClientID},
		"redirect_uri":  {googleRedirectURI},
		"response_type": {"code"},                 // Authorization Code Flow
		"scope":         {"openid email profile"}, // Request user's email and profile
		"state":         {state},                  // CSRF protection
		"access_type":   {"offline"},              // Optional: get refresh token
		"prompt":        {"consent"},              // Optional: force consent screen
	}

	authURL := fmt.Sprintf("%s?%s", googleAuthURL, params.Encode())

	// Redirect user to Google's consent screen
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

/**
 * STEP 2: OAuth Callback Handler
 * GET /auth/google/callback
 *
 * Google redirects back to this endpoint after user grants/denies permission:
 * 1. Validates the state parameter to prevent CSRF attacks
 * 2. Exchanges the authorization code for access tokens
 * 3. Uses the access token to fetch user profile from Google
 * 4. Creates or finds the user in our database
 * 5. Stores the authenticated user in the session
 * 6. Redirects to the dashboard (frontend)
 *
 * SECURITY:
 * - State validation prevents CSRF
 * - Token exchange happens server-side (client secret never exposed)
 * - User data stored in HTTP-only session cookies
 */
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, err := getSession(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		log.Printf("Session error: %v", err)
		return
	}

	// STEP 1: Validate State (CSRF Protection)
	// =========================================
	queryState := r.URL.Query().Get("state")
	sessionState, ok := session.Values[sessionStateKey].(string)

	if !ok || sessionState == "" {
		http.Error(w, "Invalid session state", http.StatusBadRequest)
		log.Println("❌ Session state not found")
		return
	}

	if queryState != sessionState {
		http.Error(w, "State mismatch - possible CSRF attack", http.StatusBadRequest)
		log.Printf("❌ State mismatch! Expected: %s, Got: %s", sessionState[:10]+"...", queryState[:10]+"...")
		return
	}

	log.Println("✅ State validated successfully")

	// Clear the state from session (one-time use)
	delete(session.Values, sessionStateKey)

	// STEP 2: Exchange Authorization Code for Tokens
	// ==============================================
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		log.Println("❌ No authorization code in callback")
		return
	}

	log.Printf("📝 Exchanging code for tokens...")

	// Prepare token exchange request
	tokenData := url.Values{
		"code":          {code},
		"client_id":     {googleClientID},
		"client_secret": {googleClientSecret}, // Server-side only - NEVER exposed to frontend
		"redirect_uri":  {googleRedirectURI},
		"grant_type":    {"authorization_code"},
	}

	// Make POST request to Google's token endpoint
	resp, err := http.PostForm(googleTokenURL, tokenData)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		log.Printf("❌ Token exchange error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		log.Printf("❌ Token exchange failed: %s", string(body))
		return
	}

	// Parse token response
	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		log.Printf("❌ Token parse error: %v", err)
		return
	}

	log.Println("✅ Successfully obtained access token")

	// STEP 3: Fetch User Profile from Google
	// ======================================
	log.Println("👤 Fetching user profile...")

	// Create HTTP request with Authorization header
	req, err := http.NewRequest("GET", googleUserInfoURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Printf("❌ Request creation error: %v", err)
		return
	}

	// Add access token to Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenResp.AccessToken))

	// Make request to get user info
	client := &http.Client{Timeout: 10 * time.Second}
	userResp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		log.Printf("❌ User info fetch error: %v", err)
		return
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(userResp.Body)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		log.Printf("❌ User info failed: %s", string(body))
		return
	}

	// Parse user profile
	var googleUser GoogleUserInfo
	if err := json.NewDecoder(userResp.Body).Decode(&googleUser); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		log.Printf("❌ User info parse error: %v", err)
		return
	}

	log.Printf("✅ Retrieved profile for: %s (%s)", googleUser.Name, googleUser.Email)

	// STEP 4: Create or Find User in Database
	// =======================================
	user := findOrCreateUser(&googleUser)

	// STEP 5: Store User in Session
	// =============================
	// Convert user to JSON for session storage
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to serialize user", http.StatusInternalServerError)
		log.Printf("❌ User serialization error: %v", err)
		return
	}

	session.Values[sessionUserKey] = string(userJSON)
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		log.Printf("❌ Session save error: %v", err)
		return
	}

	log.Printf("✅ User authenticated and session created")

	// STEP 6: Redirect to Dashboard
	// =============================
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

// ============================================================================
// API ENDPOINTS
// ============================================================================

/**
 * GET /api/user
 * Returns the currently authenticated user's information
 * Returns 401 if user is not authenticated
 */
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Check if user is in session
	userJSON, ok := session.Values[sessionUserKey].(string)
	if !ok || userJSON == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not authenticated"})
		return
	}

	// Parse user from session
	var user User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		http.Error(w, "Failed to parse user", http.StatusInternalServerError)
		return
	}

	// Return user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

/**
 * POST /api/logout
 * Destroys the user's session and logs them out
 */
func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Clear session data
	session.Values[sessionUserKey] = ""
	session.Options.MaxAge = -1 // Delete the cookie

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		log.Printf("Logout error: %v", err)
		return
	}

	log.Println("👋 User logged out")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

/**
 * CORS Middleware
 * Adds CORS headers to allow frontend requests
 */
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/**
 * Logging Middleware
 * Logs all HTTP requests
 */
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

// ============================================================================
// INITIALIZATION AND MAIN
// ============================================================================

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load Google OAuth credentials
	googleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI = os.Getenv("GOOGLE_REDIRECT_URI")

	// Validate required configuration
	if googleClientID == "" || googleClientSecret == "" {
		log.Fatal("❌ Missing required environment variables: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET")
	}

	// Set default redirect URI if not specified
	if googleRedirectURI == "" {
		googleRedirectURI = "http://localhost:8080/auth/google/callback"
		log.Printf("⚠️  Using default redirect URI: %s", googleRedirectURI)
	}

	// Initialize session store with secure settings
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "development-secret-change-in-production"
		log.Println("⚠️  Using default session secret (CHANGE IN PRODUCTION!)")
	}

	store = sessions.NewCookieStore([]byte(sessionSecret))

	// Configure session options for security
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: true,                             // Prevents JavaScript access (XSS protection)
		Secure:   os.Getenv("ENV") == "production", // HTTPS only in production
		SameSite: http.SameSiteLaxMode,             // CSRF protection
	}

	log.Println("✅ Configuration loaded successfully")
}

func main() {
	// Create HTTP router
	mux := http.NewServeMux()

	// OAuth endpoints
	mux.HandleFunc("/auth/google/start", handleGoogleStart)
	mux.HandleFunc("/auth/google/callback", handleGoogleCallback)

	// API endpoints
	mux.HandleFunc("/api/user", handleGetUser)
	mux.HandleFunc("/api/logout", handleLogout)

	// Serve static files (frontend)
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	// Wrap with middleware
	handler := loggingMiddleware(corsMiddleware(mux))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on http://localhost:%s", port)
	log.Printf("📝 Redirect URI configured: %s", googleRedirectURI)
	log.Println("👉 Visit http://localhost:" + port + " to test the application")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
