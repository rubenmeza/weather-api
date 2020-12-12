package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type weatherResource struct{}

type weatherRequest struct {
	City     string `json:"city" validate:"required"`
	Country  string `json:"country" validate:"required,len=2"`
	Forecast int    `json:"forecast"`
}
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
	wrs := weatherRequest{City: city, Country: country}
	v := validator.New()
	err := v.Struct(wrs)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	weather, err := GetOpenWeather(wrs.City, wrs.Country)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
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

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
