package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"forum/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler processes user login by validating credentials, creating a session, and setting a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Always return JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username and password are required"})
		return
	}

	if len(username) > 50 || len(username) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username must be between 4 and 50 characters"})
		return
	}

	if len(password) > 20 || len(password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password must be between 6 and 20 characters"})
		return
	}

	stmt := `SELECT id, password, userName FROM users WHERE userName = ? OR email = ?`
	row := database.DataBase.QueryRow(stmt, username, username)

	var hashPass, dbUsername string
	var userID int
	err := row.Scan(&userID, &hashPass, &dbUsername)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
		return
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	strsessionID := sessionID.String()
	expireTime := time.Now().Add(1 * time.Hour)
	stmt2 := `UPDATE users SET dateexpired = ?, session = ? WHERE id = ?`
	_, err = database.DataBase.Exec(stmt2, expireTime, strsessionID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create session"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    strsessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})

	// Return success response with user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "Login successful",
		"userId":   userID,
		"username": dbUsername,
	})
}