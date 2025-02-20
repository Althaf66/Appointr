package auth

// see youtube video
// https://www.youtube.com/watch?v=8yZImm2A1KE

// import (
// 	"context"
// 	"crypto/rand"
// 	"database/sql"
// 	"encoding/base64"
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"github.com/Althaf66/Appointr/internal/store"
// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// )

// // OAuth2 configuration
// var googleOauthConfig = &oauth2.Config{
// 	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
// 	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
// 	RedirectURL:  "http://localhost:8080/auth/google/callback",
// 	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
// 	Endpoint:     google.Endpoint,
// }

// // var jwtSecret = []byte("your-secret-key")
// // var db *sql.DB

// // User struct
// // type User struct {
// // 	ID        int       `db:"id"`
// // 	Email     string    `db:"email"`
// // 	Username  string    `db:"username"`
// // 	CreatedAt time.Time `db:"created_at"`
// // 	IsActive  bool      `db:"is_active"`
// // }

// // Generate a random state string
// func generateState() string {
// 	b := make([]byte, 16)
// 	_, _ = rand.Read(b)
// 	return base64.URLEncoding.EncodeToString(b)
// }

// // Google login handler
// func googleLogin(w http.ResponseWriter, r *http.Request) {
// 	state := generateState()
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "oauthstate",
// 		Value:    state,
// 		Expires:  time.Now().Add(1 * time.Hour),
// 		HttpOnly: true,
// 	})
// 	url := googleOauthConfig.AuthCodeURL(state)
// 	http.Redirect(w, r, url, http.StatusFound)
// }

// // Google callback handler
// func googleCallback(w http.ResponseWriter, r *http.Request) {
// 	// Verify state
// 	cookie, err := r.Cookie("oauthstate")
// 	if err != nil || r.URL.Query().Get("state") != cookie.Value {
// 		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
// 		return
// 	}

// 	code := r.URL.Query().Get("code")
// 	token, err := googleOauthConfig.Exchange(context.Background(), code)
// 	if err != nil {
// 		http.Error(w, "Code exchange failed", http.StatusUnauthorized)
// 		return
// 	}

// 	client := googleOauthConfig.Client(context.Background(), token)
// 	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
// 	if err != nil {
// 		http.Error(w, "Failed to get user info", http.StatusUnauthorized)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var googleUser struct {
// 		Email string `json:"email"`
// 		Name  string `json:"name"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
// 		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
// 		return
// 	}

// 	// Insert or update user in the database
// 	user, err := store.GauthCreate(googleUser.Email, googleUser.Name)
// 	if err != nil {
// 		http.Error(w, "Failed to save user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Generate JWT token
// 	jwtToken, err := generateJWT(user.Email)
// 	if err != nil {
// 		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
// }

// Insert or update user in the database
// func upsertUser(email, username string) (*User, error) {
// 	var user User
// 	err := db.QueryRow(`
// 		INSERT INTO users (email, username, is_active)
// 		VALUES ($1, $2, TRUE)
// 		ON CONFLICT (email) DO UPDATE
// 		SET username = EXCLUDED.username
// 		RETURNING id, email, username, created_at, is_active`,
// 		email, username,
// 	).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.IsActive)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// Generate JWT token
// func generateJWT(email string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"email": email,
// 		"exp":   time.Now().Add(time.Hour * 24).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(jwtSecret)
// }

// func main() {
// Connect to database
// r := chi.NewRouter()

// Routes
// r.Get("/auth/google/login", googleLogin)
// r.Get("/auth/google/callback", googleCallback)
// }
