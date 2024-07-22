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

func TestGetTemperatureUseCase_Ok(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	addressRepository := mocks.NewMockAddressRepository(t)
	temperatureRepository := mocks.NewMockTemperatureRepository(t)
	usecase := NewGetTemperatureUseCase(addressRepository, temperatureRepository)

	addressRepository.On("GetAddress", ctx, "11111111", tracer).Return(entity.Address{City: "Limeira", State: "SP"}, nil).Once()
	temperatureRepository.On("GetTemperature", ctx, "Limeira", tracer).Return(entity.Temperature{Celsius: 20, Fahrenheit: 68, Kelvin: 300}, nil).Once()

	temp, err := usecase.Execute(ctx, "11111111", tracer)

	assert.Nil(t, err)
	assert.Equal(t, 20.0, temp.TempC)
	assert.Equal(t, 68.0, temp.TempF)
	assert.Equal(t, 300.0, temp.TempK)
}

func TestGetTemperatureUseCase_InvalidZipCode(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	usecase := NewGetTemperatureUseCase(nil, nil)

	_, err := usecase.Execute(ctx, "11111111ABC", tracer)

	assert.NotNil(t, err)
	assert.Equal(t, entity.ErrInvalidZipCodeMsg, err.Error())
}

func TestGetTemperatureUseCase_AddressNotFound(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	addressRepository := mocks.NewMockAddressRepository(t)
	usecase := NewGetTemperatureUseCase(addressRepository, nil)

	addressRepository.On("GetAddress", ctx, "11111111", tracer).Return(entity.Address{}, errors.New(entity.ErrAddressNotFoundMsg)).Once()

	_, err := usecase.Execute(ctx, "11111111", tracer)

	assert.NotNil(t, err)
	assert.Equal(t, entity.ErrAddressNotFoundMsg, err.Error())
}

func TestGetTemperatureUseCase_TemperatureNotFound(t *testing.T) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	addressRepository := mocks.NewMockAddressRepository(t)
	temperatureRepository := mocks.NewMockTemperatureRepository(t)
	usecase := NewGetTemperatureUseCase(addressRepository, temperatureRepository)

	addressRepository.On("GetAddress", ctx, "11111111", tracer).Return(entity.Address{City: "Limeira", State: "SP"}, nil).Once()
	temperatureRepository.On("GetTemperature", ctx, "Limeira", tracer).Return(entity.Temperature{}, errors.New("error")).Once()

	_, err := usecase.Execute(ctx, "11111111", tracer)

	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Error())
}
