package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

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
	h.templates.ExecuteTemplate(w, "index.html", nil)
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

	if _, err = strconv.Atoi(query); err == nil {
		weather, err = h.weatherService.GetWeatherByZip(query)
		forecast, err = h.weatherService.GetForecastByZip(query)
	} else {
		weather, err = h.weatherService.GetWeatherByCity(query)
		forecast, err = h.weatherService.GetForecastByCity(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	forecast, err := h.weatherService.GetForecastByLatLon(lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderWeatherTemplate(w, weather, forecast)
}

func (h *Handler) renderWeatherTemplate(w http.ResponseWriter, weather *models.WeatherResponse, forecast *models.ForecastResponse) {
	data := struct {
		Weather     *models.WeatherResponse
		Forecast    *models.ForecastResponse
		CurrentTime int64
	}{
		Weather:     weather,
		Forecast:    forecast,
		CurrentTime: time.Now().Unix(),
	}

	if err := h.templates.ExecuteTemplate(w, "weather.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
