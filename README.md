# Human-Readable Weather API

Weather API is a RESTful API written in Go using go-chi and the [openweathermap](https://openweathermap.org/) API, designed for human readability

## Running from source code
```bash
$ PORT=1234 API_ID=apikeystring go run cmd/*.go
```

## Docs

### `GET weather` current weather data for city
Access current weather data for any location and get an object representing weather in a human readable way

#### Parameters
Param | Type | Required | Description
------------ | ------------- | ------------- | -------------
city | string | true | City name
country | string | true | country code, use ISO 3166 country codes
forecast | number | false | The parameter value is between 0 and 6(0 is for today)

### Response
```json
{
    "location_name": "Querétaro, MX",
    "temperature": "22 °C",
    "wind": "1.03 m/s, east-northeast",
    "cloudiness": " few clouds",
    "pressure": "1015 hPa",
    "humidity": "42%",
    "sunrise": "07:07",
    "sunset": "17:59",
    "geo_coordinates": "[-99.92, 21.00]",
    "requested_time": "2020-12-12T11:33:10-06:00",
    "forecast": {...}
}
```

#### API call example 

```bash
$ curl -i http://localhost:3333/weather\?city\=queretaro\&country\=mx
```

#### API call example with forecast

```bash
$ curl -i http://localhost:3333/weather\?city\=queretaro\&country\=mx\&forecast\=0
```

## To Do
- [x] The data must be human-readable
- [X] Use environment variables for configuration
- [X] The response must include the content-type header (application/json)
- [ ] Functions must be tested
- [ ] Keep a cache of 2 minutes of the data. You can use a persistent layer for this.
- [ ] Add an additional GET parameter that allows getting the forecast of a specific day. The parameter value is between 0 and 6(0 is for today). TIP: you have to use an additional openweather endpoint.
- [ ] add docker
- [ ] request openweathermap API concurrently using go routines