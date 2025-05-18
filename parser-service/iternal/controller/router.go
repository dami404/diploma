package controller

import "net/http"

type H func(http.ResponseWriter, *http.Request)

func SetRouter(handler H) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search/{city}/{event}", handler)
	return mux
}
