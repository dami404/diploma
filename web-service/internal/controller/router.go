package controller

import "net/http"

type H func(http.ResponseWriter, *http.Request)

func SetRouter(handler ...H) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handler[0])
	mux.HandleFunc("GET /results/{city}/{event}", handler[1])
	mux.HandleFunc("GET /last/{city}", handler[2])

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mux
}
