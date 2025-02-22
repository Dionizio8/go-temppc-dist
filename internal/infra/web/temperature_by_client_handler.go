package web

import (
	"encoding/json"
	"net/http"

	"github.com/Dionizio8/go-temppc-dist/internal/entity"
	"github.com/Dionizio8/go-temppc-dist/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WebTemperatureByClientHandler struct {
	TemppcRepository entity.TemppcRepositoryInterface
	OtelTracer       trace.Tracer
}

func NewWebTemperatureByClientHandler(temppcRepository entity.TemppcRepositoryInterface, OtelTracer trace.Tracer) *WebTemperatureByClientHandler {
	return &WebTemperatureByClientHandler{
		TemppcRepository: temppcRepository,
		OtelTracer:       OtelTracer,
	}
}

func (t *WebTemperatureByClientHandler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	var dto usecase.TemperatureInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getTemperatureByClientUserCase := usecase.NewGetTemperatureByClientUseCase(t.TemppcRepository)
	temperature, err := getTemperatureByClientUserCase.Execute(ctx, dto, t.OtelTracer)
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
