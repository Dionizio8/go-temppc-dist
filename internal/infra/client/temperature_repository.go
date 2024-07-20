package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TemperatureWeatherDTO struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type TemperatureRepository struct {
	WeatherAPIClientURL    string
	WeatherAPIClientAPIKey string
}

func NewTemperatureRepository(weatherAPIClientURL string, weatherAPIClientAPIKey string) *TemperatureRepository {
	return &TemperatureRepository{
		WeatherAPIClientURL:    weatherAPIClientURL,
		WeatherAPIClientAPIKey: weatherAPIClientAPIKey,
	}
}

func (r *TemperatureRepository) GetTemperature(ctx context.Context, city string) (entity.Temperature, error) {
	urlW, err := url.Parse(fmt.Sprintf("%s/v1/current.json", r.WeatherAPIClientURL))
	if err != nil {
		return entity.Temperature{}, err
	}
	q := urlW.Query()
	q.Set("q", city)
	q.Set("lang", "pt")
	q.Set("key", r.WeatherAPIClientAPIKey)
	urlW.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlW.String(), nil)
	if err != nil {
		return entity.Temperature{}, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.Temperature{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return entity.Temperature{}, errors.New(entity.ErrAddressNotFoundMsg)
		}
		return entity.Temperature{}, errors.New("error fetching temperature")
	}

	var temperature TemperatureWeatherDTO
	err = json.NewDecoder(resp.Body).Decode(&temperature)
	if err != nil {
		return entity.Temperature{}, err
	}

	return *entity.NewTemperature(city, temperature.Current.TempC, temperature.Current.TempF), nil
}
