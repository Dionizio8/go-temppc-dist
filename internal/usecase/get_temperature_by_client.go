package usecase

import (
	"context"
	"errors"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
)

type GetTemperatureByClientUseCase struct {
	TemppcRepository entity.TemppcRepositoryInterface
}

func NewGetTemperatureByClientUseCase(temppcRepository entity.TemppcRepositoryInterface) *GetTemperatureByClientUseCase {
	return &GetTemperatureByClientUseCase{
		TemppcRepository: temppcRepository,
	}
}

func (uc *GetTemperatureByClientUseCase) Execute(ctx context.Context, zipCode string) (TemperatureOutputDTO, error) {
	if !entity.ValidateZipCode(zipCode) {
		return TemperatureOutputDTO{}, errors.New(entity.ErrInvalidZipCodeMsg)
	}
	temperature, err := uc.TemppcRepository.GetTemperature(ctx, zipCode)
	if err != nil {
		return TemperatureOutputDTO{}, err
	}

	return TemperatureOutputDTO{
		City:  temperature.City,
		TempC: temperature.Celsius,
		TempF: temperature.Fahrenheit,
		TempK: temperature.Kelvin,
	}, nil
}
