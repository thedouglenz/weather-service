package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thedouglenz/weather-service/internal/nws"
)

func classifyTemperature(temp float64) string {
	switch {
	case temp < 32:
		return "Cold"
	case temp < 50:
		return "Cool"
	case temp < 70:
		return "Mild"
	case temp < 90:
		return "Fantastic"
	default:
		return "Hot"
	}
}

func createRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/weather/:coordinates", func(ctx *gin.Context) {
		coordinates := ctx.Param("coordinates")

		coords := strings.Split(coordinates, ",")
		if len(coords) != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid coordinates"})
		}
		lat, lon := coords[0], coords[1]

		forecastURL, err := nws.GetGridInfo(lat, lon)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		forecast, err := nws.GetForecast(forecastURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		temperature := forecast.Temperature
		temperatureClass := classifyTemperature(temperature)

		ctx.JSON(http.StatusOK, gin.H{
			"temperature":      temperature,
			"shortForecast":    forecast.ShortForecast,
			"temperatureClass": temperatureClass,
		})
	})
	return router
}

func main() {
	router := createRouter()
	router.Run(":8000")
}
