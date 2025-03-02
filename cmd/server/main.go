package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"nithish-95/weather-wapp.git/internal/handlers"
	"nithish-95/weather-wapp.git/internal/services"
	"nithish-95/weather-wapp.git/pkg/cache"
)

var templates *template.Template

func initTemplates() error {
	tmpl := template.New("").Funcs(template.FuncMap{
		"formatUnixTime": func(unixTime int64) string {
			return time.Unix(unixTime, 0).Format("3:04 PM")
		},
		"formatUnixDay": func(unixTime int64) string {
			return time.Unix(unixTime, 0).Format("02 Jan Monday")
		},
	})

	var err error
	templates, err = tmpl.ParseGlob("templates/*.html")
	return err
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("OPENWEATHER_API")
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API must be set")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	weatherCache := cache.NewWeatherCache(5 * time.Minute)
	weatherService := services.NewWeatherService(apiKey, client, weatherCache)
	ipService := services.NewIPService(client)

	if err := initTemplates(); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(securityHeaders)

	handler := handlers.NewHandler(weatherService, ipService, templates)

	r.Get("/", handler.Index)
	r.Get("/ip", handler.IPReport)
	r.Get("/weather", handler.WeatherReport)
	r.Get("/weather/latlon", handler.LatLonReport)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
