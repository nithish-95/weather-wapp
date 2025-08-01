package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"nithish-95/weather-wapp.git/internal/models"
	"nithish-95/weather-wapp.git/internal/services"
)

type Handler struct {
	weatherService *services.WeatherService
	ipService      *services.IPService
	templates      *template.Template
}

func NewHandler(ws *services.WeatherService, ips *services.IPService, t *template.Template) *Handler {
	return &Handler{
		weatherService: ws,
		ipService:      ips,
		templates:      t,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error executing index template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) IPReport(w http.ResponseWriter, r *http.Request) {
	ipInfo, err := h.ipService.GetIPInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := h.weatherService.GetWeatherByCity(ipInfo.City)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	forecast, err := h.weatherService.GetForecastByCity(ipInfo.City)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderWeatherTemplate(w, weather, forecast)
}

func (h *Handler) WeatherReport(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	var weather *models.WeatherResponse
	var forecast *models.ForecastResponse
	var err error

	// Check if query is a ZIP code or a city name
	if _, err = strconv.Atoi(query); err == nil {
		weather, err = h.weatherService.GetWeatherByZip(query)
		if err == nil {
			forecast, err = h.weatherService.GetForecastByZip(query)
		}
	} else {
		weather, err = h.weatherService.GetWeatherByCity(query)
		if err == nil {
			forecast, err = h.weatherService.GetForecastByCity(query)
		}
	}

	if err != nil {
		log.Printf("Error fetching weather data for query '%s': %v", query, err)
		http.Error(w, "Could not retrieve weather data. Please try again.", http.StatusInternalServerError)
		return
	}

	h.renderWeatherTemplate(w, weather, forecast)
}

func (h *Handler) LatLonReport(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	if lat == "" || lon == "" {
		http.Error(w, "lat and lon parameters are required", http.StatusBadRequest)
		return
	}

	weather, err := h.weatherService.GetWeatherByLatLon(lat, lon)
	if err != nil {
		log.Printf("Error fetching weather by lat/lon: %v", err)
		http.Error(w, "Error fetching weather data for your location.", http.StatusInternalServerError)
		return
	}

	forecast, err := h.weatherService.GetForecastByLatLon(lat, lon)
	if err != nil {
		log.Printf("Error fetching forecast by lat/lon: %v", err)
		http.Error(w, "Error fetching forecast data for your location.", http.StatusInternalServerError)
		return
	}

	h.renderWeatherTemplate(w, weather, forecast)
}

// renderWeatherTemplate now includes checks for nil data before executing
func (h *Handler) renderWeatherTemplate(w http.ResponseWriter, weather *models.WeatherResponse, forecast *models.ForecastResponse) {
	// Crucial check: ensure weather and forecast data are not nil before rendering.
	if weather == nil || forecast == nil {
		log.Println("Render error: weather or forecast data is nil.")
		http.Error(w, "Could not retrieve complete weather data. The location may not be supported.", http.StatusNotFound)
		return
	}

	data := struct {
		Weather  *models.WeatherResponse
		Forecast *models.ForecastResponse
	}{
		Weather:  weather,
		Forecast: forecast,
	}

	err := h.templates.ExecuteTemplate(w, "weather.html", data)
	if err != nil {
		// This will catch errors within the template itself (e.g., calling a non-existent field)
		log.Printf("Error executing weather template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
