package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/dami404/diploma-web/internal/usecase"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

type Page struct {
	Title string
}

func (h *Handler) StartPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Show start page")
	ts, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to parse html", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func (h *Handler) Results(w http.ResponseWriter, r *http.Request) {
	event := r.PathValue("event")
	city := r.PathValue("city")

	h.uc.SaveEvent(event, city)

	tickets := h.uc.ProfitEvents(event, city)

	data, err := json.Marshal(tickets)
	if err != nil {
		log.Println("failed to marshal data")
		return
	}

	w.Write(data)
}

func (h *Handler) LastEvents(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")
	tickets := h.uc.LastEvents(city)

	data, err := json.Marshal(tickets)
	if err != nil {
		log.Print("failed to marshal data")
		return
	}

	w.Write(data)
}
