package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type OpenWeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
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

	if res.StatusCode != 200 {
		e := new(OpenWeatherError)
		err = json.NewDecoder(res.Body).Decode(e)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(e.Cod + " " + e.Message)
	}

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
	direction := getWindCardinalDirection(owr.Wind.Deg)
	return fmt.Sprintf("%0.2f m/s, %s", owr.Wind.Speed, direction)
}

func (owr *OpenWeatherResponse) getHumanReadableCloudiness() string {
	cloudines := getCloudines(*owr)
	return fmt.Sprintf("%s", cloudines)
}

func (owr *OpenWeatherResponse) getHumanReadablePressure() string {
	return fmt.Sprintf("%d hPa", owr.Main.Pressure)
}

func (owr *OpenWeatherResponse) getHumanReadableHumidity() string {
	return fmt.Sprintf("%d%%", owr.Main.Humidity)
}

func (owr *OpenWeatherResponse) getHumanReadableSunrise() string {
	sunrise := getHourMinutes(owr.Sys.Sunrise)
	return fmt.Sprintf("%s", sunrise)
}

func (owr *OpenWeatherResponse) getHumanReadableSunset() string {
	sunset := getHourMinutes(owr.Sys.Sunset)
	return fmt.Sprintf("%s", sunset)
}

func (owr *OpenWeatherResponse) getHumanReadableGeoCoordinates() string {
	return fmt.Sprintf("[%0.2f, %0.2f]", owr.Coord.Lon, owr.Coord.Lat)
}

func getWindCardinalDirection(deg int) string {
	cardinalDirections := [16]string{"north", "north-northeast", "northeast", "east-northeast", "east", "east-southeast", "southeast", "south-southeast", "south", "south-southwest", "southwest", "west-southwest", "west", "west-northwest", "northwest", "north-northwest"}
	cardinalIndex := int((float64(deg) / 22.5) + 0.5)
	return cardinalDirections[cardinalIndex]
}

func getCloudines(owr OpenWeatherResponse) string {
	var cloudines string
	for _, w := range owr.Weather {
		if w.ID >= 800 && w.ID < 900 {
			cloudines = cloudines + " " + w.Description
		}
	}
	return cloudines
}

func getHourMinutes(unixTimestamp int) string {
	t := time.Unix(int64(unixTimestamp), 0)
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}
