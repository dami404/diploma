package entity

type Event struct {
	Name    string   `json:"name"`
	City    string   `json:"city"`
	Tickets []Ticket `json:"tickets"`
}

type Ticket struct {
	Price    int    `json:"price"`
	Location string `json:"location"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Url      string `json:"url"`
}
