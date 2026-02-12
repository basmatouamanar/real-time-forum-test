package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"forum/database"
	"forum/helpers"
	"forum/tools"
)

// HanldlerShowHome handles the home page display, checking user session and rendering posts, categories, and reactions.
func HanldlerShowHome(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	var userID int

	// Détecte si c'est une requête API
	isAPI := r.URL.Path == "/api/posts"

	if !isAPI && r.URL.Path != "/" {
		helpers.Errorhandler(w, "Status Not Found!", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, errSession := r.Cookie("session")
	if errSession == nil && cookie.Value != "" {
		var userExists bool
		var expiredTime time.Time
		err := database.DataBase.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE session = ?)", cookie.Value).Scan(&userExists)
		if err == nil && userExists {
			loggedIn = true
			err = database.DataBase.QueryRow("SELECT id, dateexpired FROM users WHERE session = ?", cookie.Value).Scan(&userID, &expiredTime)
			if err == nil {
				if expiredTime.After(time.Now()) {
					loggedIn = true
				} else {
					_, err = database.DataBase.Exec(
						"UPDATE users SET session = NULL, dateexpired = NULL WHERE session = ?", cookie.Value)
					if err != nil {
						helpers.Errorhandler(w, "internal server error", http.StatusInternalServerError)
						return
					}

					expiredCookie := &http.Cookie{
						Name:     "session",
						Value:    "",
						Path:     "/",
						MaxAge:   -1,
						Expires:  time.Now().Add(-1 * time.Hour),
						HttpOnly: true,
					}
					http.SetCookie(w, expiredCookie)
					loggedIn = false
				}
			} else if err == sql.ErrNoRows {
				if isAPI {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "Session not found"})
					return
				}
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	dataIsLogin := tools.IsLogin{LoggedIn: loggedIn, UserID: userID}

	posts := helpers.GetAllPosts(w)

	// Si c'est /api/posts, retourne JUSTE les posts en JSON


	// Sinon, continue avec le reste pour le HTML
	categories := helpers.GetAllCategories(w)
	reactionStats, err := helpers.GetAllReactionStats()
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userReactions, err := helpers.GetUserPostReactions(userID)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	comments := helpers.GetAllComments(w)
	connectUserName := helpers.GetConnectUserName(w, userID)

	commentReactionStats, err := helpers.GetAllCommentReactionStats()
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userCommentReactions, err := helpers.GetUserCommentReactions(userID)
	if err != nil {
		helpers.Errorhandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var pageData tools.PageData
	pageData.Posts = posts
	pageData.Categories = categories
	pageData.IsLogin = dataIsLogin
	pageData.ReactionStats = reactionStats
	pageData.UserReactions = userReactions
	pageData.Comment = comments
	pageData.ConnectUserName = connectUserName
	pageData.CommentReactionStats = commentReactionStats
	pageData.UserCommentReactions = userCommentReactions
	if isAPI {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pageData)
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		return
	}
}
