package nws

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Forecast struct {
	ShortForecast string  `json:"shortForecast"`
	Temperature   float64 `json:"temperature"`
}

const pointsBaseURL = "https://api.weather.gov/points"

func GetGridInfo(lat, lon string) (string, error) {
	apiURL := fmt.Sprintf("%s/%s,%s", pointsBaseURL, lat, lon)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to make request to NWS API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NWS API returned non-200 status code: %d", resp.StatusCode)
	}

	var forecastResponse struct {
		Properties struct {
			ForecastURL string `json:"forecast"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&forecastResponse); err != nil {
		return "", fmt.Errorf("failed to decode NWS API response: %v", err)
	}

	forecastURL := forecastResponse.Properties.ForecastURL
	return forecastURL, nil
}

func GetForecast(forecastURL string) (*Forecast, error) {
	resp, err := http.Get(forecastURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to NWS Forecast API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NWS API returned non-200 status code: %d", resp.StatusCode)
	}

	var forecastResponse struct {
		Properties struct {
			Periods []Forecast `json:"periods"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&forecastResponse); err != nil {
		return nil, fmt.Errorf("failed to decode NWS Forecast API response: %v", err)
	}

	if len(forecastResponse.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast data found")
	}

	forecast := &Forecast{
		ShortForecast: forecastResponse.Properties.Periods[0].ShortForecast,
		Temperature:   forecastResponse.Properties.Periods[0].Temperature,
	}

	return forecast, nil
}
