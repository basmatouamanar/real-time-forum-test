package helpers

import (
	"encoding/json"
	"net/http"
)

func TransferToJSON(w http.ResponseWriter, comment string, postID string) {
	response := struct {
		PostID  string `json:"postId"`
		Comment string `json:"comment"`
		Message string `json:"message"`
	}{
		PostID:  postID,
		Comment: comment,
		Message: "Comment added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Ajoute ça pour être explicite
	
	json.NewEncoder(w).Encode(response)
}