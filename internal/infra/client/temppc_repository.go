package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TemppctDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type TemppcRepository struct {
	TemppcAPIClientURL string
}

func NewTemppcRepository(temppcAPIClientURL string) *TemppcRepository {
	return &TemppcRepository{
		TemppcAPIClientURL: temppcAPIClientURL,
	}
}

func (r *TemppcRepository) GetTemperature(ctx context.Context, zipCode string, tracer trace.Tracer) (entity.Temperature, error) {
	ctx, span := tracer.Start(ctx, "get-cep-temp")
	defer span.End()

	url := fmt.Sprintf("%s/temppc/temperature/%s", r.TemppcAPIClientURL, zipCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		switch resp.StatusCode {
		case http.StatusNotFound:
			return entity.Temperature{}, errors.New(entity.ErrAddressNotFoundMsg)
		case http.StatusUnprocessableEntity:
			return entity.Temperature{}, errors.New(entity.ErrInvalidZipCodeMsg)
		default:
			return entity.Temperature{}, errors.New("error fetching temperature")
		}
	}

	var temppct TemppctDTO
	err = json.NewDecoder(resp.Body).Decode(&temppct)
	if err != nil {
		return entity.Temperature{}, err
	}

	return entity.Temperature{
		City:       temppct.City,
		Celsius:    temppct.TempC,
		Fahrenheit: temppct.TempF,
		Kelvin:     temppct.TempK,
	}, nil
}
