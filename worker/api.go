package worker

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
}

type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
}

func (a *Api) Start() {
	a.initRouter()
	addressPort := fmt.Sprintf("%s:%d", a.Address, a.Port)
	fmt.Println("Listening on: ", addressPort)
	http.ListenAndServe(fmt.Sprintf(addressPort), a.Router)
}
