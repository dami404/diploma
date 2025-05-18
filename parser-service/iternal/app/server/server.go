package server

import (
	"net/http"
)

type Server struct {
	notify chan error
	Router *http.ServeMux
	Address
}

type Address struct {
	Host string
	Port string
}

func NewHttpServer(host string, port string) *Server {
	return &Server{
		notify: make(chan error, 1),
		Router: http.NewServeMux(),
		Address: Address{
			Port: port,
			Host: host,
		},
	}
}

func (s *Server) Start() {
	go func() {
		s.notify <- http.ListenAndServe(s.Address.Host+":"+s.Address.Port, s.Router)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}
