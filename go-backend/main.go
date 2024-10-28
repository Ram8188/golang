package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

// User struct to represent user data
type User struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Username string `json:"username"` // Add this line for username

	Password string `json:"password"` // Storing password in plain text (not recommended)
}

// Response struct to represent API responses
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Store sessions in a cookie store using an environment variable for the secret key
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_KEY")))

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
}

func main() {
	// Check if the environment variable for the secret key is set
	if os.Getenv("SESSION_SECRET_KEY") == "" {
		log.Fatal("SESSION_SECRET_KEY environment variable is not set")
	}

	db, err := sql.Open("mysql", "u316144551_policy_point:Policy_point_123@tcp(193.203.184.6)/u316144551_policy_point")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping the database to confirm it's reachable
	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	http.HandleFunc("/api/register", registerUser(db))
	http.HandleFunc("/api/login", loginUser(db))
	http.HandleFunc("/api/verify-session", verifySession)

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)

		if r.Method == http.MethodOptions {
			return // Handle preflight request
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Save user to the database without hashing the password
		_, err = db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
		if err != nil {
			log.Printf("Error creating user: %v", err) // Log detailed error
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Successful registration response
		response := Response{
			Success: true,
			Message: "User registered successfully",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func loginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)

		if r.Method == http.MethodOptions {
			return // Handle preflight request
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Log the email being queried for debugging
		log.Printf("Attempting to log in with email: %s", user.Email)

		// Fetch the user from the database
		var dbUser User
		err = db.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Username, &dbUser.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}
			log.Printf("Database error: %v", err) // Log detailed error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Compare the plain passwords directly
		if dbUser.Password != user.Password {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Create a new session
		session, _ := store.Get(r, "session-name")
		session.Values["userID"] = dbUser.ID
		session.Values["email"] = dbUser.Email
		session.Values["username"] = dbUser.Username // Store username in session
		session.Save(r, w)

		// Successful login response
		response := Response{
			Success: true,
			Message: "Login successful",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response) // Encode and send the response
	}
}

func verifySession(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session:", err) // Log the error
		http.Error(w, "Session not valid, please log in again.", http.StatusUnauthorized)
		return
	}

	if session.Values["userID"] == nil {
		log.Println("Session is invalid: userID not found") // Log invalid session
		http.Error(w, "Session not valid, please log in again.", http.StatusUnauthorized)
		return
	}

	// If the session is valid, log the session details (for debugging)
	log.Println("Session is valid for userID:", session.Values["userID"])

	// Session is valid
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session is valid"))
}
