package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const apiKey = "b3e15c4737290fe97fed1c359bd47d12"

type WeatherResponse struct {
	Name string `json: name`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func getWeather(zipCode string) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s&units=metric", zipCode, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}

func weatherReporter(w http.ResponseWriter, r *http.Request) {
	zipcode := chi.URLParam(r, "zipcode")

	if zipcode == "" {
		http.Error(w, "Zip code is required", http.StatusBadRequest)
		return
	}

	weather, err := getWeather(zipcode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := fmt.Sprintf("The current temperature in %s with zip code %s is %v°C", weather.Name, zipcode, weather.Main.Temp)
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
