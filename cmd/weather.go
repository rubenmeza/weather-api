package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type weatherResource struct{}

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
	weather := GetOpenWeather(city, country)
	fmt.Println(weather.Coord.Lon)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respose, _ := json.Marshal(weather)
	w.Write(respose)
}

func (wr weatherResource) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	w.Write([]byte(city))
}

func (wr weatherResource) GetWeatherAndForecast(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetWeatherAndForecast"))
}
