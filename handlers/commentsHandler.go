package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/database"
	"forum/helpers"
	//"forum/helpers"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// TOUJOURS renvoyer du JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	// Vérifier la session MANUELLEMENT
	cookie, err := r.Cookie("session") // Change "session" par le nom de ton cookie si différent
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	cookieValue := cookie.Value
	userID, errSelect := database.SelectUserID("SELECT id FROM users WHERE session = ?", cookieValue)
	if errSelect == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid session"})
		return
	} else if errSelect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	errParse := r.ParseForm()
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid form data"})
		return
	}

	commentText := strings.TrimSpace(r.FormValue("comment"))
	postIDStr := r.FormValue("postId")

	if commentText == "" || postIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comment and postId are required"})
		return
	}

	if len(commentText) > 200 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comment too long (max 200 characters)"})
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid post ID"})
		return
	}

	query := `INSERT INTO comments (comment, postId, userId, creationDate) 
	          VALUES (?, ?, ?, ?)`

	creationDate := time.Now().Format("2006-01-02 15:04:05")
	errExec := database.ExecuteData(query, commentText, postID, userID, creationDate)
	if errExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save comment"})
		return
	}

	username := helpers.GetConnectUserName(w, userID)

	// Réponse de succès
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"postId":   postIDStr,
		"comment":  commentText,
		"userName": username, 

		"message": "Comment added successfully",
	})
}

/*
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]string{
    "postId":  postIDStr,
    "comment": commentText,
    "userName": username, // <-- voilà ce qu'il manquait
    "message": "Comment added successfully",
})

*/
