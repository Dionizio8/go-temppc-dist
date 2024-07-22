package entity

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type AddressRepositoryInterface interface {
	GetAddress(ctx context.Context, zipCode string, tracer trace.Tracer) (Address, error)
}

type TemperatureRepositoryInterface interface {
	GetTemperature(ctx context.Context, city string, tracer trace.Tracer) (Temperature, error)
}

type TemppcRepositoryInterface interface {
	GetTemperature(ctx context.Context, zipCode string, tracer trace.Tracer) (Temperature, error)
}
