package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dami404/diploma-db/iternal/entity"
	"github.com/dami404/diploma-db/iternal/usecase"
)

type DBHandler struct {
	uc usecase.DBUsecase
}

func NewDBHandler(uc usecase.DBUsecase) *DBHandler {
	return &DBHandler{uc: uc}
}

func (h *DBHandler) Save(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Event entity.Event `json:"ticket"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	} else if err := h.uc.Save(data.Event); err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Data saved successfully"))
	}

}

func (h *DBHandler) LastEvents(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")

	events, err := h.uc.LastEvents(city)
	if err != nil {
		http.Error(w, "Failed to fetch popular events", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
