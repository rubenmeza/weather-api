package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenWeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat int     `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type OpenWeatherError struct {
	Cod     string `json:"cod"`
	Message string `json:"message"`
}

// GetOpenWeather get weather by city and country
func GetOpenWeather(city string, country string) (*OpenWeatherResponse, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?q=" + city + "," + country + "&units=metric&appid=1508a9a4840a5574c822d70ca2132032"
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	w := new(OpenWeatherResponse)
	err = json.NewDecoder(res.Body).Decode(w)

	if err != nil {
		return nil, err
	}

	return w, nil
}

func (owr *OpenWeatherResponse) getHumanReadableLocation() string {
	return fmt.Sprintf("%s, %s", owr.Name, owr.Sys.Country)
}

func (owr *OpenWeatherResponse) getHumanReadableTemperature() string {
	return fmt.Sprintf("%.0f °C", owr.Main.Temp)
}

func (owr *OpenWeatherResponse) getHumanReadableWind() string {
	return fmt.Sprintf("%s, %s", owr.Wind.Speed, owr.Wind.Deg)
}

func (owr *OpenWeatherResponse) getHumanReadableCloudiness() string {
	// necesario iterar sobre owr.Weather
	return fmt.Sprintf("%s", owr.Clouds.All)
}

func (owr *OpenWeatherResponse) getHumanReadablePressure() string {
	return fmt.Sprintf("%s", owr.Main.Pressure)
}

func (owr *OpenWeatherResponse) getHumanReadableHumidity() string {
	return fmt.Sprintf("%s", owr.Main.Humidity)
}

func (owr *OpenWeatherResponse) getHumanReadableSunrise() string {
	return fmt.Sprintf("%s", owr.Sys.Sunrise)
}

func (owr *OpenWeatherResponse) getHumanReadableSunset() string {
	return fmt.Sprintf("%s", owr.Sys.Sunset)
}

func (owr *OpenWeatherResponse) getHumanReadableGeoCoordinates() string {
	return fmt.Sprintf("[%f, %i]", owr.Coord.Lon, owr.Coord.Lat)
}
