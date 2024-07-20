package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"github.com/Dionizio8/go-temppc-dist/mocks"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestWebTemperatureHandler_GetTemperature_ErrorddressNotFound(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	addressRepository := mocks.NewMockAddressRepository(t)
	temperatureRepository := mocks.NewMockTemperatureRepository(t)
	webTemperatureHandler := NewWebTemperatureHandler(addressRepository, temperatureRepository, tracer)

	addressRepository.On("GetAddress", mock.Anything, "11111111").Return(entity.Address{}, errors.New(entity.ErrAddressNotFoundMsg))

	r := chi.NewRouter()
	r.Get("/temperature/{zipCode}", webTemperatureHandler.GetTemperature)

	req := httptest.NewRequest(http.MethodGet, "/temperature/11111111", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, entity.ErrAddressNotFoundMsg, w.Body.String())
}

func TestWebTemperatureHandler_GetTemperature_ErrorInvalidZipCode(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	addressRepository := mocks.NewMockAddressRepository(t)
	temperatureRepository := mocks.NewMockTemperatureRepository(t)
	webTemperatureHandler := NewWebTemperatureHandler(addressRepository, temperatureRepository, tracer)

	r := chi.NewRouter()
	r.Get("/temperature/{zipCode}", webTemperatureHandler.GetTemperature)

	req := httptest.NewRequest(http.MethodGet, "/temperature/11111111ABC", nil) // Invalid zip code
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, entity.ErrInvalidZipCodeMsg, w.Body.String())
}

func TestWebTemperatureHandler_GetTemperature_ErrorInternalServer(t *testing.T) {
	tracer := noop.NewTracerProvider().Tracer("microservice-tracer")
	addressRepository := mocks.NewMockAddressRepository(t)
	temperatureRepository := mocks.NewMockTemperatureRepository(t)
	webTemperatureHandler := NewWebTemperatureHandler(addressRepository, temperatureRepository, tracer)

	addressRepository.On("GetAddress", mock.Anything, "11111111").Return(entity.Address{City: "Limeira", State: "SP"}, nil).Once()
	temperatureRepository.On("GetTemperature", mock.Anything, "Limeira").Return(entity.Temperature{}, errors.New("Internal Server Error")).Once()

	r := chi.NewRouter()
	r.Get("/temperature/{zipCode}", webTemperatureHandler.GetTemperature)

	req := httptest.NewRequest(http.MethodGet, "/temperature/11111111", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}
