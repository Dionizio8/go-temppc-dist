package usecase

import (
	"context"
	"errors"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

type TemperatureOutputDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTemperatureUseCase struct {
	AddressRepository     entity.AddressRepositoryInterface
	TemperatureRepository entity.TemperatureRepositoryInterface
}

func NewGetTemperatureUseCase(addressRepository entity.AddressRepositoryInterface, temperatureRepository entity.TemperatureRepositoryInterface) *GetTemperatureUseCase {
	return &GetTemperatureUseCase{
		AddressRepository:     addressRepository,
		TemperatureRepository: temperatureRepository,
	}
}

func (uc *GetTemperatureUseCase) Execute(ctx context.Context, zipCode string, tracer trace.Tracer) (TemperatureOutputDTO, error) {
	if !entity.ValidateZipCode(zipCode) {
		return TemperatureOutputDTO{}, errors.New(entity.ErrInvalidZipCodeMsg)
	}
	address, err := uc.AddressRepository.GetAddress(ctx, zipCode, tracer)
	if err != nil {
		return TemperatureOutputDTO{}, err
	}

	if address.City == "" {
		return TemperatureOutputDTO{}, errors.New(entity.ErrAddressNotFoundMsg)
	}

	temperature, err := uc.TemperatureRepository.GetTemperature(ctx, address.City, tracer)
	if err != nil {
		return TemperatureOutputDTO{}, err
	}

	return TemperatureOutputDTO{
		City:  address.City,
		TempC: temperature.Celsius,
		TempF: temperature.Fahrenheit,
		TempK: temperature.Kelvin,
	}, nil
}
