package entity

import (
	"regexp"
)

const (
	ErrInvalidZipCodeMsg  = "invalid zipcode"
	ErrAddressNotFoundMsg = "can not find zipcode"
)

type Address struct {
	City  string
	State string
}

func NewAddress(city, state string) *Address {
	return &Address{
		City:  city,
		State: state,
	}
}

func ValidateZipCode(zipCode string) bool {
	regex := regexp.MustCompile(`^[0-9]{8}$`)
	return regex.MatchString(zipCode)
}
