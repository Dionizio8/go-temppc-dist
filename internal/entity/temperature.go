package entity

type Temperature struct {
	City       string
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NewTemperature(city string, celsius, fahrenheit float64) *Temperature {
	return &Temperature{
		City:       city,
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     celsius + 273,
	}
}
