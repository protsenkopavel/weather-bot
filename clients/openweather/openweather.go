package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
	}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	if city == "" {
		return Coordinates{}, fmt.Errorf("city is empty")
	}

	endpoint := "http://api.openweathermap.org/geo/1.0/direct"
	u := fmt.Sprintf("%s?q=%s&limit=5&appid=%s", endpoint, url.QueryEscape(city), o.apiKey)

	resp, err := http.Get(u)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error get coordinates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("fail get coordinates: status=%d", resp.StatusCode)
	}

	var coordinatesResponse []CoordinatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&coordinatesResponse); err != nil {
		return Coordinates{}, fmt.Errorf("error decode response: %w", err)
	}

	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("no coordinates found for city=%q", city)
	}

	return Coordinates{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}, nil
}
