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

type AddressViaCepDTO struct {
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type AddressRepository struct {
	ViaCEPClientURL string
}

func NewAddressRepository(viaCEPClientURL string) *AddressRepository {
	return &AddressRepository{
		ViaCEPClientURL: viaCEPClientURL,
	}
}

func (r *AddressRepository) GetAddress(ctx context.Context, zipCode string, tracer trace.Tracer) (entity.Address, error) {
	ctx, span := tracer.Start(ctx, "get-city-name")
	defer span.End()

	url := fmt.Sprintf("%s/ws/%s/json", r.ViaCEPClientURL, zipCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entity.Address{}, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return entity.Address{}, errors.New(entity.ErrAddressNotFoundMsg)
		}
		return entity.Address{}, err
	}

	var addressViaCep AddressViaCepDTO
	err = json.NewDecoder(resp.Body).Decode(&addressViaCep)
	if err != nil {
		return entity.Address{}, err
	}

	return entity.Address{
		City:  addressViaCep.Localidade,
		State: addressViaCep.Uf,
	}, nil
}
