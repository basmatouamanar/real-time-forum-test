package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"forum/database"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	cookieValue := cookie.Value
	var userID int
	err = database.DataBase.QueryRow("SELECT id FROM users WHERE session = ?", cookieValue).Scan(&userID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid session"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	errParse := r.ParseForm()
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid form data"})
		return
	}

	postIDStr := r.FormValue("postId")
	reactionStr := r.FormValue("reaction")

	if postIDStr == "" || reactionStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "postId and reaction are required"})
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid post ID"})
		return
	}

	reaction, err := strconv.Atoi(reactionStr)
	if err != nil || (reaction != 1 && reaction != -1) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid reaction (must be 1 or -1)"})
		return
	}

	var postExists int
	errSelect := database.DataBase.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ?", postID).Scan(&postExists)
	if errSelect == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Post not found"})
		return
	} else if errSelect != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	if postExists == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Post not found"})
		return
	}

	var existingReaction int
	err = database.DataBase.QueryRow("SELECT reaction FROM postReactions WHERE userId = ? AND postId = ?", userID, postID).Scan(&existingReaction)

	switch err {
	case sql.ErrNoRows:
		// Nouvelle réaction
		_, err = database.DataBase.Exec("INSERT INTO postReactions (userId, postId, reaction) VALUES(?, ?, ?)", userID, postID, reaction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add reaction"})
			return
		}
	case nil:
		if existingReaction == reaction {
			// Supprimer la réaction (toggle off)
			_, err = database.DataBase.Exec("DELETE FROM postReactions WHERE userId = ? AND postId = ?", userID, postID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove reaction"})
				return
			}
		} else {
			// Changer la réaction
			_, err = database.DataBase.Exec("UPDATE postReactions SET reaction = ? WHERE userId = ? AND postId = ?", reaction, userID, postID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update reaction"})
				return
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// ========== PARTIE AJOUTÉE : Calculer les compteurs ==========
	var likesCount, dislikesCount int
	
	// Compter les likes
	err = database.DataBase.QueryRow("SELECT COUNT(*) FROM postReactions WHERE postId = ? AND reaction = 1", postID).Scan(&likesCount)
	if err != nil {
		likesCount = 0
	}
	
	// Compter les dislikes
	err = database.DataBase.QueryRow("SELECT COUNT(*) FROM postReactions WHERE postId = ? AND reaction = -1", postID).Scan(&dislikesCount)
	if err != nil {
		dislikesCount = 0
	}
	
	// Récupérer la réaction actuelle de l'utilisateur
	var userReaction int
	err = database.DataBase.QueryRow("SELECT reaction FROM postReactions WHERE userId = ? AND postId = ?", userID, postID).Scan(&userReaction)
	if err == sql.ErrNoRows {
		userReaction = 0 // Aucune réaction
	}
	// ============================================================

	// Réponse de succès avec les compteurs
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likesCount":    likesCount,
		"dislikesCount": dislikesCount,
		"userReaction":  userReaction,
		"message":       "Reaction processed successfully",
	})
}