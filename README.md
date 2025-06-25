# Weather Application

A modern web application that provides real-time weather information and forecasts. Built with Go and featuring a beautiful, responsive UI, this application offers multiple ways to get weather data based on your location or search preferences.

## Live Demo
https://weather.app.nithish.net

## Key Features

1. **Smart Location Detection**:
   - Automatically detects user location using browser geolocation
   - Falls back to IP-based location if geolocation is denied
   - Manual search options for any location

2. **Weather Information**:
   - Real-time current weather conditions
   - Temperature (current, feels like, min/max)
   - Wind speed and direction
   - Humidity levels
   - Detailed weather descriptions with icons

3. **5-Day Weather Forecast**:
   - Future weather predictions
   - Temperature trends
   - Weather conditions with icons
   - Wind and humidity forecasts

4. **Search Options**:
   - Search by city name (works worldwide)
   - Search by ZIP code (US locations)
   - Automatic location detection

5. **User Interface**:
   - Clean, modern design using Tailwind CSS
   - Responsive layout for all devices
   - Interactive weather cards
   - Loading animations
   - Error handling with user-friendly messages

## Technology Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML5, Tailwind CSS, JavaScript
- **APIs**: OpenWeatherMap API, IP Geolocation API
- **Caching**: In-memory caching for better performance
- **Deployment**: Docker support

## Quick Start

### Local Setup

1. Clone the repository:
```bash
git clone git@github.com:nithish-95/weather-wapp.git
cd weather-wapp
```

2. Install dependencies:
```bash
go get -u github.com/go-chi/chi/v5
go mod tidy
```

3. Create `.env` file with your OpenWeather API key:
```
OPENWEATHER_API=your_api_key_here
PORT=3000
```

4. Run the application:
```bash
make
```

The application will be available at http://localhost:3000

### Docker Setup

1. Build the Docker image:
```bash
docker build --progress plain --no-cache -t weatherapp1 .
```

2. Run the container:
```bash
docker run -p 3000:3000 weatherapp1
```

Access the application at http://localhost:3000

## Project Structure

```
.
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   └── services/        # Weather and IP services
├── pkg/
│   └── cache/          # Caching implementation
├── templates/          # HTML templates
├── Dockerfile         # Docker configuration
└── Makefile          # Build automation
```
