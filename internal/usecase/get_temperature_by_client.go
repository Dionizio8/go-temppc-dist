package usecase

import (
	"context"
	"errors"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
)

type TemperatureInputDTO struct {
	Cep string `json:"cep"`
}

type GetTemperatureByClientUseCase struct {
	TemppcRepository entity.TemppcRepositoryInterface
}

func NewGetTemperatureByClientUseCase(temppcRepository entity.TemppcRepositoryInterface) *GetTemperatureByClientUseCase {
	return &GetTemperatureByClientUseCase{
		TemppcRepository: temppcRepository,
	}
}

func (uc *GetTemperatureByClientUseCase) Execute(ctx context.Context, input TemperatureInputDTO) (TemperatureOutputDTO, error) {
	if !entity.ValidateZipCode(input.Cep) {
		return TemperatureOutputDTO{}, errors.New(entity.ErrInvalidZipCodeMsg)
	}
	temperature, err := uc.TemppcRepository.GetTemperature(ctx, input.Cep)
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
