package controller

import "net/http"

type H func(w http.ResponseWriter, r *http.Request)

func SetRouter(searchHandler H, saveHandler H) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /save", saveHandler)
	mux.HandleFunc("GET /last/{city}", searchHandler)
	return mux
}
