package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type app struct {
	Route   *mux.Router
	ListApi map[string]listApi
}

type listApi struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Method   string `json:"method"`
	Response string `json:"response,omitempty"`
}

type App interface {
	Start()
}

func New() App {
	return app{
		ListApi: make(map[string]listApi),
		Route:   mux.NewRouter(),
	}
}

func (a app) Start() {

	a.Route.HandleFunc("/api", a.handleAPI).Methods("GET", "POST")
	handler := cors.Default().Handler(a.Route)

	log.Println("Running on port 8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal(err)
	}
}
