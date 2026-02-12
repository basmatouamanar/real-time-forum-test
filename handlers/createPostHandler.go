package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"forum/database"
	"forum/helpers"
)

/*
CreatePostHandler handles the creation of a new post. It processes a
multipart/form-data request, validates the input, saves an optional image,
and inserts the post and its associated categories into the database.
Returns JSON response with the created post data.
*/
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// TOUJOURS renvoyer du JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	errParse := r.ParseMultipartForm(10 << 20)
	if errParse != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse form"})
		return
	}

	cookieValue := helpers.GetCookieValue(w, r)
	if cookieValue == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	userID := helpers.GetUserID(cookieValue)

	title, ok := r.PostForm["title"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title is required"})
		return
	}

	description, ok := r.PostForm["description"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Description is required"})
		return
	}

	description[0] = strings.ReplaceAll(description[0], "\r", "")
	description[0] = strings.TrimSpace(description[0])
	title[0] = strings.TrimSpace(title[0])

	if len(title[0]) == 0 || len(description[0]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title and description cannot be empty"})
		return
	}

	if len(title[0]) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title must be less than 100 characters"})
		return
	}

	if len(description[0]) > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Description must be less than 1000 characters"})
		return
	}

	categories, ok := r.PostForm["categories"]
	if !ok || len(categories) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please select at least one category"})
		return
	}

	categoriesID := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, catsID := range categories {
		if !slices.Contains(categoriesID, catsID) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category selected"})
			return
		}
	}

	imagePath := ""
	imageFile, fileHeader, err := r.FormFile("choose-file")
	if err != nil {
		if err != http.ErrMissingFile {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
			return
		}
		if len(r.MultipartForm.File) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid file upload"})
			return
		}
	} else {
		defer imageFile.Close()

		if !helpers.IsImage(fileHeader.Filename) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid image extension. Only jpg, jpeg, png, gif allowed"})
			return
		}

		const maxSize = 2 * 1024 * 1024
		if fileHeader.Size > maxSize {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Image size must be less than 2 MB"})
			return
		}

		query := `SELECT imageUrl FROM posts ORDER BY creationDate DESC LIMIT 1;`
		imgName := ""
		err := database.DataBase.QueryRow(query).Scan(&imgName)
		if err != nil && err != sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
			return
		}

		if imgName == "" {
			imgName = "image1" + strings.ToLower(filepath.Ext(fileHeader.Filename))
		} else {
			numWithExt := imgName[len("/static/upload/image"):]
			nb := numWithExt[:len(numWithExt)-len(filepath.Ext(numWithExt))]
			num, errConv := strconv.Atoi(nb)
			if errConv != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
				return
			}
			num++
			imgName = "image" + strconv.Itoa(num) + filepath.Ext(fileHeader.Filename)
		}

		file, errCreate := os.Create("./static/upload/" + imgName)
		if errCreate != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save image"})
			return
		}
		defer file.Close()

		_, errCopy := io.Copy(file, imageFile)
		if errCopy != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save image"})
			return
		}
		imagePath = "/static/upload/" + imgName
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	queryInsertPost := "INSERT INTO posts (title, post, imageUrl, userId, creationDate) VALUES (?, ?, ?, ?, ?)"
	errEx := database.ExecuteData(queryInsertPost, title[0], description[0], imagePath, userID, timeNow)
	if errEx != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create post"})
		return
	}

	lastPostID, err := database.SelectLastIdOfPosts("SELECT id FROM posts ORDER BY creationDate DESC LIMIT 1;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	queryInsertCategory := "INSERT INTO postCategories (postId, categoryId) VALUES (?, ?)"
	for _, catID := range categories {
		categoryID, err := strconv.Atoi(catID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category ID"})
			return
		}
		errExec := database.ExecuteData(queryInsertCategory, lastPostID, categoryID)
		if errExec != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add categories"})
			return
		}
	}

	// Récupérer le nom de l'utilisateur pour la réponse
	var username string
	database.DataBase.QueryRow("SELECT userName FROM users WHERE id = ?", userID).Scan(&username)

	// Renvoyer les données du post créé
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     "Post created successfully",
		"postId":      lastPostID,
		"title":       title[0],
		"description": description[0],
		"imageUrl":    imagePath,
		"author":      username,
		"date":        timeNow,
		"categories":  categories,
	})
}