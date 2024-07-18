package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"time"
	"weather-cli/config"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	}
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch    int64   `json:"time_epoch"`
				TempC        float64 `json:"temp_c"`
				ChanceOfRain float64 `json:"chance_of_rain"`
				Condition    struct {
					Text string `json:"text"`
				} `json:"condition"`
			}
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	secretKey := config.SetUp().ApiKey
	city := "Seoul"
	if len(os.Args) >= 2 {
		fmt.Println(os.Args[0])
		city = os.Args[1]
	}

	url := "http://api.weatherapi.com/v1/forecast.json?key=" + secretKey + "&q=" + city + "&days=1&aqi=no&alerts=no"
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("not 200")
	}

	body, err := io.ReadAll(res.Body)

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%s - %.1fC\n", location.Name, current.TempC)
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		message := fmt.Sprintf("%s - %.0fC, %.0f, %s\n",
			date.Format("15:05"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if date.Before(time.Now()) {
			continue
		}

		if hour.ChanceOfRain < 50 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}

}
