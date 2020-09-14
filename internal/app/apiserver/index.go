package apiserver

import (
	"html/template"
	"net/http"
)

// HandlerRoot ...
func (s *APIServer) HandlerRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		templ, _ := template.ParseFiles("templates/index.html")
		index := "index"
		templ.Execute(w, index)
	}
}
