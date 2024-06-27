package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const apiKey = "f34f18a4c3ca9bd80f6cb96488136858"

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

type ForecastResponse struct {
	ListF []struct {
		DtTxtF int64 `json:"dt"`
		MainF  struct {
			TempF      float64 `json:"temp"`
			FeelsLikeF float64 `json:"feels_like"`
			HumidityF  int     `json:"humidity"`
		} `json:"main"`
		WeatherF []struct {
			DescriptionF string `json:"description"`
			IconF        string `json:"icon"`
		} `json:"weather"`
		WindF struct {
			SpeedF float64 `json:"speed"`
		} `json:"wind"`
	} `json:"list"`
}

type IpResponse struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Zip  string  `json:"zip"`
}

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
	forcast, err := getForcastByCity(ipInfo.City)
	if err != nil {
		http.Error(w, "Could not get weather forcast info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, weather, forcast, "weather.html")
}

func getWeatherByZip(zipCode string) (*WeatherResponse, error) {
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
func getForcastByZip(zipCode string) (*ForecastResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?zip=%s&appid=%s&units=metric", zipCode, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var forecastResponse ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResponse); err != nil {
		return nil, err
	}

	return &forecastResponse, nil
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

func getForcastByCity(cityName string) (*ForecastResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", url.QueryEscape(cityName), apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var forecastResponse ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResponse); err != nil {
		return nil, err
	}

	return &forecastResponse, nil
}

func WeatherReport(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "Query not found", http.StatusNotFound)
		return
	}

	var weather *WeatherResponse
	var forcast *ForecastResponse
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

	if _, convErr := strconv.Atoi(query); convErr == nil {

		forcast, err = getForcastByZip(query)
	} else {

		forcast, err = getForcastByCity(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderTemplate(w, weather, forcast, "weather.html")
}

func getWeatherByLatLon(lat, lon string) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)
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
func getForcastByLatLon(lat, lon string) (*ForecastResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var forecastResponse ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResponse); err != nil {
		return nil, err
	}

	return &forecastResponse, nil
}

func latlonReporter(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if lat == "" || lon == "" {
		http.Error(w, "Latitude and Longitude not found", http.StatusBadRequest)
		return
	}
	weather, err := getWeatherByLatLon(lat, lon)
	if err != nil {
		log.Printf("Could not get weather info: %v", err)
		http.Error(w, "Could not get weather info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	forcast, err := getForcastByLatLon(lat, lon)
	if err != nil {
		log.Printf("Could not get forecast info: %v", err)
		http.Error(w, "Could not get forecast info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, weather, forcast, "weather.html")
}

func formatUnixTime(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	return t.Format("3:04 PM")
}
func formatUnixDay(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	return t.Format("02 Jan Monday")
}

func renderTemplate(w http.ResponseWriter, weather *WeatherResponse, forecast *ForecastResponse, file string) {
	tmpl := template.New("weather.html").Funcs(template.FuncMap{
		"formatUnixTime": func(unixTime int64) string {
			return formatUnixTime(unixTime)
		},
		"formatUnixDay": func(unixTime int64) string {
			return formatUnixDay(unixTime)
		},
	})
	tmpl, err := tmpl.ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Weather     *WeatherResponse
		Forecast    *ForecastResponse
		CurrentTime int64
	}{
		Weather:     weather,
		Forecast:    forecast,
		CurrentTime: time.Now().Unix(),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexpage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexpage)
	r.Get("/ip", IpReporter)
	r.Get("/weather", WeatherReport)
	r.Get("/weather/latlon", latlonReporter)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Print(err)
	}
}
