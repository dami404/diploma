package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dami404/diploma-parser/iternal/entity"
	"github.com/dami404/diploma-parser/iternal/usecase"
)

type ParserHandler struct {
	usecase usecase.ParserUsecase
}

func NewParserHandler(usecase usecase.ParserUsecase) *ParserHandler {
	return &ParserHandler{usecase: usecase}
}

func parseUrl(url *http.Request) *entity.Event {
	event := url.PathValue("event")
	city := url.PathValue("city")
	return &entity.Event{Name: event, City: city}
}

func (h *ParserHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller.SearchEvents")
	event := parseUrl(r)
	tickets := h.usecase.ProfitEvents(*event)
	w.WriteHeader(http.StatusAccepted)
	data, err := json.Marshal(tickets)
	if err != nil {
		log.Print("failed to marshal data")
		return
	}
	w.Write(data)

}
