package apiserver

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// HandlerPods ...
func (s *APIServer) HandlerPods() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/pods/" {
			io.WriteString(w, "Pods help ")
		} else {
			vars := mux.Vars(r)
			io.WriteString(w, "Hello: "+vars["ns"])
		}

	}
}
