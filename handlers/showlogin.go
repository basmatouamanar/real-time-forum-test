package handlers

import (
	"net/http"
	"text/template"

	"forum/helpers"
)

// show login page with render func

func Showloginhandler(w http.ResponseWriter, r *http.Request) {
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
