package handlers

import (
	"net/http"
	"text/template"

	"forum/helpers"
)

// Showregister displays the registration page to the user.
func Showregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.Errorhandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
