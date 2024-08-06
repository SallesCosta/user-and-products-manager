package database

import (
	"github.com/sallescosta/user-and-products-manager/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := u.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetAllUsers() ([]entity.User, error) {
	var usersList []entity.User

	err := u.DB.Find(&usersList).Error
	if err != nil {
		return nil, err
	}

	return usersList, nil
}
