package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"nithish-95/weather-wapp.git/internal/models"
	"nithish-95/weather-wapp.git/pkg/cache"
)

type WeatherService struct {
	apiKey string
	client *http.Client
	cache  *cache.WeatherCache
}

func NewWeatherService(apiKey string, client *http.Client, cache *cache.WeatherCache) *WeatherService {
	return &WeatherService{
		apiKey: apiKey,
		client: client,
		cache:  cache,
	}
}

func (ws *WeatherService) GetWeatherByCity(city string) (*models.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("weather:city:%s", city)
	if cached, found := ws.cache.Get(cacheKey); found {
		return cached.(*models.WeatherResponse), nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", url.QueryEscape(city), ws.apiKey)
	resp, err := ws.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	ws.cache.Set(cacheKey, &weather)
	return &weather, nil
}

func (ws *WeatherService) GetForecastByCity(city string) (*models.ForecastResponse, error) {
	cacheKey := fmt.Sprintf("forecast:city:%s", city)
	if cached, found := ws.cache.Get(cacheKey); found {
		return cached.(*models.ForecastResponse), nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", url.QueryEscape(city), ws.apiKey)
	resp, err := ws.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status: %d", resp.StatusCode)
	}

	var forecast models.ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	ws.cache.Set(cacheKey, &forecast)
	return &forecast, nil
}

func (ws *WeatherService) GetWeatherByZip(zip string) (*models.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("weather:zip:%s", zip)
	if cached, found := ws.cache.Get(cacheKey); found {
		return cached.(*models.WeatherResponse), nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?zip=%s&appid=%s&units=metric", url.QueryEscape(zip), ws.apiKey)
	resp, err := ws.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	ws.cache.Set(cacheKey, &weather)
	return &weather, nil
}

func (ws *WeatherService) GetForecastByZip(zip string) (*models.ForecastResponse, error) {
	cacheKey := fmt.Sprintf("forecast:zip:%s", zip)
	if cached, found := ws.cache.Get(cacheKey); found {
		return cached.(*models.ForecastResponse), nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?zip=%s&appid=%s&units=metric", url.QueryEscape(zip), ws.apiKey)
	resp, err := ws.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status: %d", resp.StatusCode)
	}

	var forecast models.ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	ws.cache.Set(cacheKey, &forecast)
	return &forecast, nil
}

func (ws *WeatherService) GetWeatherByLatLon(lat, lon string) (*models.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("weather:latlon:%s,%s", lat, lon)
	if cached, found := ws.cache.Get(cacheKey); found {
		return cached.(*models.WeatherResponse), nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=metric", url.QueryEscape(lat), url.QueryEscape(lon), ws.apiKey)
	resp, err := ws.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	ws.cache.Set(cacheKey, &weather)
	return &weather, nil
}

func (ws *WeatherService) GetForecastByLatLon(lat, lon string) (*models.ForecastResponse, error) {
	cacheKey := fmt.Sprintf("forecast:latlon:%s,%s", lat, lon)
	if cached, found := ws.cache.Get(cacheKey); found {
		return cached.(*models.ForecastResponse), nil
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&appid=%s&units=metric", url.QueryEscape(lat), url.QueryEscape(lon), ws.apiKey)
	resp, err := ws.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status: %d", resp.StatusCode)
	}

	var forecast models.ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	ws.cache.Set(cacheKey, &forecast)
	return &forecast, nil
}

// Similar methods for other endpoints (GetForecastByCity, GetWeatherByZip, etc.)
