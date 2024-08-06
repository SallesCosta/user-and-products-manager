package database

import "github.com/sallescosta/user-and-products-manager/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	GetAllUsers() ([]entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindById(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
	FindAll(page, limit int, sort string) (ProductResponse, error)
}
