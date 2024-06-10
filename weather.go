package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

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
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int64   `json:"sunrise"`
		Sunset  int64   `json:"sunset"`
	} `json:"sys"`
}

type IpResponse struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Zip  string  `json:"zip"`
}

// http://ip-api.com/json/
func getIp(ip string) (*IpResponse, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ipResponse IpResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipResponse); err != nil {
		return nil, err
	}
	return &ipResponse, nil
}

func getWeatherByZip(zipCode string) (*WeatherResponse, error) {
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

func IpReporter(w http.ResponseWriter, r *http.Request) {
	userIP := r.Header.Get("X-Forwarded-For")

	ipInfo, err := getIp(userIP)
	if err != nil {
		http.Error(w, "Could not get IP info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := getWeatherByCity(ipInfo.City)
	if err != nil {
		http.Error(w, "Could not get weather info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, weather)
}

func WeatherReport(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, r.RequestURI)

	query := r.URL.Query().Get("query")

	for k, v := range r.Header {
		slog.Info("Header", "key", k, "value", v)
	}
	fmt.Println(r.Header.Get("Accept-Encoding"))
	if query == "" {
		http.Error(w, "Query not found", http.StatusNotFound)
		return
	}

	var weather *WeatherResponse
	var err error

	if _, convErr := strconv.Atoi(query); convErr == nil {
		weather, err = getWeatherByZip(query)
	} else {
		weather, err = getWeatherByCity(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderTemplate(w, weather)
}

func renderTemplate(w http.ResponseWriter, data *WeatherResponse) {
	tmpl, err := template.ParseFiles("weather.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Serve the form
	r.Get("/", IpReporter)
	r.Get("/weather", WeatherReport)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Print(err)
	}
}
