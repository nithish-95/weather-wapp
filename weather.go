package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func weatherReporter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Weather web app demo"))
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", weatherReporter)

	http.ListenAndServe(":3000", r)

}
