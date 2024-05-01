package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("book", 15)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "book", product.Name)
	assert.Equal(t, 15, product.Price)
	assert.NotEmpty(t, product.ID)
	assert.NotEmpty(t, product.CreatedAt)
	assert.Nil(t, product.Validate())
}

func TestProductValidations_Name(t *testing.T) {
	product, err := NewProduct("", 15)
	assert.Nil(t, product)
	assert.Equal(t, err, ErrNameIsRequired)
}

func TestProductValidations_No_Price(t *testing.T) {
	product, err := NewProduct("book", 0)
	assert.Nil(t, product)
	assert.Equal(t, err, ErrPriceIsRequired)
}

func TestProductValidations_Invalid_Price(t *testing.T) {
	product, err := NewProduct("book", -15)
	assert.Nil(t, product)
	assert.Equal(t, err, ErrInvalidPrice)
}
