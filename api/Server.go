package api

import (
	"Onepiece/router"
	"Onepiece/routes"
	"fmt"
	"net/http"
)

type Server struct {
	listenAddress string
}

func NewServer(listenAddress int) *Server {
	var Addr string = fmt.Sprintf(":%d", listenAddress)
	return &Server{
		listenAddress: Addr,
	}
}

const (
	POST   = "POST"
	GET    = "GET"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func (s *Server) Start() error {
	Router := router.NewRouter()

	Router.Insert("/api/home", routes.Home, GET)

	http.HandleFunc("/", Router.Router)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	return http.ListenAndServe(s.listenAddress, nil)
}
