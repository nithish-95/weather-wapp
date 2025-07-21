package models

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Name    string `json:"name"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"` // Added Pressure
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"` // Added Visibility
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
}

// ForecastResponse struct contains the 5-day forecast data.
type ForecastResponse struct {
	List []struct {
		Dt    int64  `json:"dt"`
		DtTxt string `json:"dt_txt"`
		Main  struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"` // Added TempMin for forecast
			TempMax   float64 `json:"temp_max"` // Added TempMax for forecast
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	} `json:"list"`
}

// IPResponse struct holds data from the IP geolocation API.
type IPResponse struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Zip  string  `json:"zip"`
}
