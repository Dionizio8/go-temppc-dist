package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"github.com/Dionizio8/go-temppc-dist/mocks"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestGetTemperatureByClientUseCase_Ok(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	temppcRepository := mocks.NewMockTemppcRepository(t)
	usecase := NewGetTemperatureByClientUseCase(temppcRepository)

	temppcRepository.On("GetTemperature", ctx, "11111111", tracer).Return(entity.Temperature{Celsius: 20, Fahrenheit: 68, Kelvin: 300}, nil)

	temp, err := usecase.Execute(ctx, TemperatureInputDTO{Cep: "11111111"}, tracer)

	assert.Nil(t, err)
	assert.Equal(t, 20.0, temp.TempC)
	assert.Equal(t, 68.0, temp.TempF)
	assert.Equal(t, 300.0, temp.TempK)
}

func TestGetTemperatureByClientUseCase_NotFound(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	temppcRepository := mocks.NewMockTemppcRepository(t)
	usecase := NewGetTemperatureByClientUseCase(temppcRepository)

	temppcRepository.On("GetTemperature", ctx, "11111111", tracer).Return(entity.Temperature{}, errors.New(entity.ErrAddressNotFoundMsg))

	temp, err := usecase.Execute(ctx, TemperatureInputDTO{Cep: "11111111"}, tracer)

	assert.NotNil(t, err)
	assert.Equal(t, entity.ErrAddressNotFoundMsg, err.Error())
	assert.Equal(t, "", temp.City)
}

func TestGetTemperatureByClientUseCase_InvalidZipCode(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	usecase := NewGetTemperatureByClientUseCase(nil)

	_, err := usecase.Execute(ctx, TemperatureInputDTO{Cep: "124ABC"}, tracer)

	assert.NotNil(t, err)
	assert.Equal(t, entity.ErrInvalidZipCodeMsg, err.Error())
}
