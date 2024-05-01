package database

import (
	"github.com/sallescosta/crud-api/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) FindById(id string) (*entity.Product, error) {
	var product entity.Product

	err := p.DB.Find(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.FindById(product.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Save(product).Error
}

func (p *Product) Delete(product *entity.Product) error {
	_, err := p.FindById(product.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Delete(product).Error
}
