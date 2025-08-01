package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"nithish-95/weather-wapp.git/internal/handlers"
	"nithish-95/weather-wapp.git/internal/services"
	"nithish-95/weather-wapp.git/pkg/cache"
)

var templates *template.Template

// toFloat is a helper function to convert various numeric types to float64.
func toFloat(v interface{}) float64 {
	// Using reflect to handle different possible numeric types (int, float64, etc.)
	val := reflect.ValueOf(v)
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return float64(val.Int())
	case float32, float64:
		return val.Float()
	default:
		return 0
	}
}

// initTemplates initializes the HTML templates and registers custom functions.
func initTemplates() error {
	tmpl := template.New("").Funcs(template.FuncMap{
		// NEWLY ADDED function to fix the error
		"float": toFloat,

		// Your existing functions
		"formatUnixTime": func(unixTime int64) string {
			return time.Unix(unixTime, 0).Format("3:04 PM")
		},
		"formatUnixDay": func(unixTime int64) string {
			return time.Unix(unixTime, 0).Format("Mon")
		},
		"div": func(a, b interface{}) float64 {
			fa := toFloat(a)
			fb := toFloat(b)
			if fb == 0 {
				return 0
			}
			return fa / fb
		},
		"sub": func(a, b interface{}) float64 {
			return toFloat(a) - toFloat(b)
		},
		"mul": func(a, b interface{}) float64 {
			return toFloat(a) * toFloat(b)
		},
		"formatUnixDate": func(unixTime int64) string {
			return time.Unix(unixTime, 0).Format("Jan 2")
		},
		"sunPosition": func(sunrise, sunset int64) float64 {
			now := time.Now().Unix()
			if now < sunrise {
				return 0
			}
			if now > sunset {
				return 100
			}
			duration := float64(sunset - sunrise)
			progress := float64(now - sunrise)
			if duration == 0 {
				return 50 // Avoid division by zero
			}
			return (progress / duration) * 100
		},
		"currentYear": func() int {
			return time.Now().Year()
		},
		"currentTime": func() string {
			return time.Now().Format("15:04")
		},
		// Mock values for now - these would ideally come from another API
		"uvIndex":       func(lat, lon float64) int { return 5 },   // Placeholder
		"moonPhase":     func() string { return "Waxing Gibbous" }, // Placeholder
		"moonPhaseName": func() string { return "Waxing Gibbous" }, // Placeholder
	})

	var err error
	templates, err = tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
	}
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

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
