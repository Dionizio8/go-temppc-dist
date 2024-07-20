package entity

import "context"

type AddressRepositoryInterface interface {
	GetAddress(ctx context.Context, zipCode string) (Address, error)
}

type TemperatureRepositoryInterface interface {
	GetTemperature(ctx context.Context, city string) (Temperature, error)
}

type TemppcRepositoryInterface interface {
	GetTemperature(ctx context.Context, zipCode string) (Temperature, error)
}
