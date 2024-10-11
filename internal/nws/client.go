package nws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// National Woeather Service API
// https://www.weather.gov/documentation/services-web-api

type Forecast struct {
	ShortForecast string  `json:"shortForecast"`
	Temperature   float64 `json:"temperature"`
}

const pointsBaseURL = "https://api.weather.gov/points"

var forecastCache sync.Map

func roundToTwoDecimalPlaces(num float64) float64 {
	return float64(int(num*100)) / 100
}

func GetGridInfo(lat, lon string) (string, error) {

	// Convert the latitude and longitude to float64 and round to two decimal places
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)
	latKey, lonKey := roundToTwoDecimalPlaces(latFloat), roundToTwoDecimalPlaces(lonFloat)

	// Cache the grid info for the given broader latitude and longitude
	cacheKey := fmt.Sprintf("%f,%f", latKey, lonKey)

	// Check if the forecast URL is already cached
	if forecastURL, ok := forecastCache.Load(cacheKey); ok {
		fmt.Printf("Retrieved forecast URL from cache for %s\n", cacheKey)
		return forecastURL.(string), nil
	}

	apiURL := fmt.Sprintf("%s/%s,%s", pointsBaseURL, lat, lon)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to make request to NWS API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NWS API returned non-200 status code: %d", resp.StatusCode)
	}

	// We need only the forecast URL from the response in this step
	var forecastResponse struct {
		Properties struct {
			ForecastURL string `json:"forecast"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&forecastResponse); err != nil {
		return "", fmt.Errorf("failed to decode NWS API response: %v", err)
	}

	// Cache the result
	forecastCache.Store(cacheKey, forecastResponse.Properties.ForecastURL)

	return forecastResponse.Properties.ForecastURL, nil
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
