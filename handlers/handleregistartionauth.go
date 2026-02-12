package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"forum/database"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler processes user registration by validating inputs, creating user account, and setting up session.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Always return JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	// Get form values
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	firstname := strings.TrimSpace(r.FormValue("firstname"))
	lastname := strings.TrimSpace(r.FormValue("lastname"))
	ageStr := strings.TrimSpace(r.FormValue("age"))
	gender := strings.TrimSpace(r.FormValue("gender"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("passwordre"))

	// Validate required fields
	if nickname == "" || firstname == "" || lastname == "" || ageStr == "" || gender == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tous les champs sont obligatoires"})
		return
	}

	// Validate nickname length
	if len(nickname) < 4 || len(nickname) > 50 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Le pseudo doit contenir entre 4 et 50 caractères"})
		return
	}

	// Validate firstname length
	if len(firstname) < 2 || len(firstname) > 50 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Le prénom doit contenir entre 2 et 50 caractères"})
		return
	}

	// Validate lastname length
	if len(lastname) < 2 || len(lastname) > 50 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Le nom doit contenir entre 2 et 50 caractères"})
		return
	}

	// Validate age
	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 13 || age > 120 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "L'âge doit être entre 13 et 120 ans"})
		return
	}

	// Validate gender
	if gender != "male" && gender != "female" && gender != "other" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Genre invalide"})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format d'e-mail invalide"})
		return
	}

	// Validate password length
	if len(password) < 6 || len(password) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Le mot de passe doit contenir entre 6 et 20 caractères"})
		return
	}

	// Check if nickname already exists
	var existsNickname bool
	err = database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE userName = ?)", nickname).Scan(&existsNickname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	if existsNickname {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ce pseudo est déjà utilisé"})
		return
	}

	// Check if email already exists
	var existsEmail bool
	err = database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&existsEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	if existsEmail {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cet e-mail est déjà utilisé"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Generate session
	sessionID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	strsessionID := sessionID.String()
	expireTime := time.Now().Add(1 * time.Hour)

	// Insert user into database
	stmt := `INSERT INTO users (userName, firstname, lastname, age, gender, email, password, session, dateexpired) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DataBase.Exec(stmt, nickname, firstname, lastname, age, gender, email, string(hashedPassword), strsessionID, expireTime)
	if err != nil {
		fmt.Println("ERROR Database Insert:", err)
		fmt.Println("Values:", nickname, firstname, lastname, age, gender, email)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create account"})
		return
	}

	// Get the inserted user ID
	userID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Set session cookie
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
		"message":  "Inscription réussie",
		"userId":   userID,
		"username": nickname,
	})
}
