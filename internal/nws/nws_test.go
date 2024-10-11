package nws

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetForecast(t *testing.T) {
	mockResponse := `{
        "properties": {
            "periods": [
                {
                    "shortForecast": "Partly Cloudy",
                    "temperature": 75
                }
            ]
        }
    }`

	// Create a mock server that returns the mockResponse
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	forecast, err := GetForecast(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if forecast.Temperature != 75 {
		t.Errorf("Expected temperature to be 75, got %v", forecast.Temperature)
	}
}
