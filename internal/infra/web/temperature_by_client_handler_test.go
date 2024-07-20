package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"github.com/Dionizio8/go-temppc-dist/internal/usecase"
	"github.com/Dionizio8/go-temppc-dist/mocks"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestWebTemperatureByClientHandler_GetTemperature_ErrorInvalidZipCode(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	temppcRepository := mocks.NewMockTemppcRepository(t)
	webTemperatureByClientHandler := NewWebTemperatureByClientHandler(temppcRepository, tracer)

	r := chi.NewRouter()
	r.Post("/temperature", webTemperatureByClientHandler.GetTemperature)

	invalidZipCodeDTO := usecase.TemperatureInputDTO{
		Cep: "11111111ABC", // Invalid zip code
	}
	jsonInvalidZipCodeDTO, _ := json.Marshal(invalidZipCodeDTO)

	req := httptest.NewRequest(http.MethodPost, "/temperature", bytes.NewBuffer(jsonInvalidZipCodeDTO))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, entity.ErrInvalidZipCodeMsg, w.Body.String())
}

func TestWebTemperatureByClientHandler_GetTemperature_ErrorInternalServer(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	temppcRepository := mocks.NewMockTemppcRepository(t)
	webTemperatureByClientHandler := NewWebTemperatureByClientHandler(temppcRepository, tracer)

	temppcRepository.On("GetTemperature", mock.Anything, "11111111").Return(entity.Temperature{}, errors.New("Internal Server Error")).Once()

	r := chi.NewRouter()
	r.Post("/temperature", webTemperatureByClientHandler.GetTemperature)

	zipCodeDTO := usecase.TemperatureInputDTO{
		Cep: "11111111",
	}
	jsonZipCodeDTO, _ := json.Marshal(zipCodeDTO)

	req := httptest.NewRequest(http.MethodPost, "/temperature", bytes.NewBuffer(jsonZipCodeDTO))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestWebTemperatureByClientHandler_GetTemperature_Success(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	temppcRepository := mocks.NewMockTemppcRepository(t)
	webTemperatureByClientHandler := NewWebTemperatureByClientHandler(temppcRepository, tracer)

	temperature := entity.Temperature{
		City:       "Limeira",
		Celsius:    25.5,
		Fahrenheit: 77.9,
		Kelvin:     100.00,
	}

	temppcRepository.On("GetTemperature", mock.Anything, "11111111").Return(temperature, nil).Once()

	r := chi.NewRouter()
	r.Post("/temperature", webTemperatureByClientHandler.GetTemperature)

	zipCodeDTO := usecase.TemperatureInputDTO{
		Cep: "11111111",
	}
	jsonZipCodeDTO, _ := json.Marshal(zipCodeDTO)

	req := httptest.NewRequest(http.MethodPost, "/temperature", bytes.NewBuffer(jsonZipCodeDTO))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"city":"Limeira","temp_C":25.5,"temp_F":77.9,"temp_K":100}`, w.Body.String())
}

func TestWebTemperatureByClientHandler_GetTemperature_ErrorAddressNotFound(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	temppcRepository := mocks.NewMockTemppcRepository(t)
	webTemperatureByClientHandler := NewWebTemperatureByClientHandler(temppcRepository, tracer)

	temppcRepository.On("GetTemperature", mock.Anything, "11111111").Return(entity.Temperature{}, errors.New(entity.ErrAddressNotFoundMsg)).Once()

	r := chi.NewRouter()
	r.Post("/temperature", webTemperatureByClientHandler.GetTemperature)

	zipCodeDTO := usecase.TemperatureInputDTO{
		Cep: "11111111",
	}
	jsonZipCodeDTO, _ := json.Marshal(zipCodeDTO)

	req := httptest.NewRequest(http.MethodPost, "/temperature", bytes.NewBuffer(jsonZipCodeDTO))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, entity.ErrAddressNotFoundMsg, w.Body.String())
}
