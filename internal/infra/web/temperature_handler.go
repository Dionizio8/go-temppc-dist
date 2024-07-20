package web

import (
	"encoding/json"
	"net/http"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"github.com/Dionizio8/go-temppc-dist/internal/usecase"
	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WebTemperatureHandler struct {
	AddressRepository     entity.AddressRepositoryInterface
	TemperatureRepository entity.TemperatureRepositoryInterface
	OtelTracer            trace.Tracer
}

func NewWebTemperatureHandler(addressRepository entity.AddressRepositoryInterface, temperatureRepository entity.TemperatureRepositoryInterface, OtelTracer trace.Tracer) *WebTemperatureHandler {
	return &WebTemperatureHandler{
		AddressRepository:     addressRepository,
		TemperatureRepository: temperatureRepository,
		OtelTracer:            OtelTracer,
	}
}

func (t *WebTemperatureHandler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := t.OtelTracer.Start(ctx, "GetTemperature")
	defer span.End()

	zipCode := chi.URLParam(r, "zipCode")
	if zipCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	getTemperatureUserCase := usecase.NewGetTemperatureUseCase(t.AddressRepository, t.TemperatureRepository)
	temperature, err := getTemperatureUserCase.Execute(ctx, zipCode)
	if err != nil {
		msgErr := err.Error()
		if msgErr == entity.ErrAddressNotFoundMsg {
			w.WriteHeader(http.StatusNotFound)
		} else if msgErr == entity.ErrInvalidZipCodeMsg {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(msgErr))
		return
	}

	err = json.NewEncoder(w).Encode(temperature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
