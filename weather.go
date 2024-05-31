package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const apiKey = "f34f18a4c3ca9bd80f6cb96488136858"

// http://api.openweathermap.org/data/2.5/weather?q=Dearborn&appid=f34f18a4c3ca9bd80f6cb96488136858&units=metric
type WeatherResponse struct {
	Name    string `json:"name"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func getWeather(zipCode string) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s&units=metric", zipCode, apiKey)
	// url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", url.QueryEscape(zipCode), apiKey)
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

func getWeatherByCity(cityName string) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", url.QueryEscape(cityName), apiKey)
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

// :3000/zipcode/{Zipcode}
func ZipcodeReport(w http.ResponseWriter, r *http.Request) {
	// zipcode := chi.URLParam(r, "zipcode")
	zipcode := r.URL.Query().Get("zipcode")

	if zipcode == "" {
		http.Error(w, "Zipcode not found", http.StatusNotFound)
		return
	}

	weather, err := getWeather(zipcode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("weather.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, weather)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

// http://localhost:3000/weather?name={CityName}
func CityNameReport(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("name")

	if cityName == "" {
		http.Error(w, "City name not found", http.StatusNotFound)
		return
	}

	weather, err := getWeatherByCity(cityName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("weather.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, weather)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func mainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Serve the form
	r.Get("/", mainPage)

	// r.Get("/zipcode/{zipcode}", ZipcodeReport)
	r.Get("/zipcode", ZipcodeReport)
	r.Get("/weather", CityNameReport)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Print(err)
	}
}
