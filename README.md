# Weather Service

Coding exercise. Proxy to NWS with temperature classification.

### REST API: 

Routes:

* `GET /weather/:lat,long`: Example:  `{ shortForecast: "Sunny", temperature: 75, temperatureClass: "Fantastic" }`

### Run locally

1. `go run main.go`
2. `curl 'localhost:8000/weather/40.040943,-85.9520099'`