package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddress_Ok(t *testing.T) {
	address := NewAddress("Limeira", "SP")

	assert.Equal(t, "Limeira", address.City)
	assert.Equal(t, "SP", address.State)
}
