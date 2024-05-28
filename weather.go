package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const apiKey = "b3e15c4737290fe97fed1c359bd47d12"

func weatherReporter(w http.ResponseWriter, r *http.Request) {
	zipcode := chi.URLParam(r, "zipcode")

	if zipcode == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := fmt.Sprintf("Weather at the %v is -----", zipcode)
	w.Write([]byte(resp))
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{zipcode}", weatherReporter)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Print(err)
	}

}
