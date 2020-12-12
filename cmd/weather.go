package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type weatherResource struct{}

type weatherResponse struct {
	LocationName   string `json:"location_name"`
	Temperature    string `json:"temperature"`
	Wind           string `json:"wind"`
	Cloudiness     string `json:"cloudiness"`
	Pressure       string `json:"pressure"`
	Humidity       string `json:"humidity"`
	Sunrise        string `json:"sunrise"`
	Sunset         string `json:"sunset"`
	GeoCoordinates string `json:"geo_coordinates"`
	RequestedTime  string `json:"requested_time"`
}

func (wr weatherResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", wr.GetWeatherQueryParams)

	r.Route("/{city}/{country}", func(r chi.Router) {
		r.Get("/", wr.GetWeather)
	})

	r.Route("/{city}/{country}/{forecast}", func(r chi.Router) {
		r.Get("/", wr.GetWeatherAndForecast)
	})

	return r
}

func (wr weatherResource) GetWeatherQueryParams(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")
	weather, err := GetOpenWeather(city, country)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("upss something happend"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	hr := buildWeatherResponse(*weather)
	respose, _ := json.Marshal(hr)
	w.Write(respose)
}

func (wr weatherResource) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	w.Write([]byte(city))
}

func (wr weatherResource) GetWeatherAndForecast(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetWeatherAndForecast"))
}

func buildWeatherResponse(owr OpenWeatherResponse) weatherResponse {
	requestTime := time.Now()
	wr := weatherResponse{
		LocationName:   owr.getHumanReadableLocation(),
		Temperature:    owr.getHumanReadableTemperature(),
		Wind:           owr.getHumanReadableWind(),
		Cloudiness:     owr.getHumanReadableCloudiness(),
		Pressure:       owr.getHumanReadablePressure(),
		Humidity:       owr.getHumanReadableHumidity(),
		Sunrise:        owr.getHumanReadableSunrise(),
		Sunset:         owr.getHumanReadableSunset(),
		GeoCoordinates: owr.getHumanReadableGeoCoordinates(),
		RequestedTime:  requestTime.Format(time.RFC3339),
	}
	return wr
}
