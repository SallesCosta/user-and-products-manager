# 1 - Create the configs package (config.go, .env)

## Cria arquivo de configuração
```sh
# root/configs/config.go

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver      string           `mapstructure:"DB_DRIVER"`
	DBHost        string           `mapstructure:"DB_HOST"`
	DBPort        string           `mapstructure:"DB_PORT"`
	DBPassword    string           `mapstrucutre:"DB_PASSWORD"`
	DBName        string           `mapstrucutre:"DB_NAME"`
	WebServerPort string           `mapstrucutre:"WEB_SERVER_PORT"`
	JWTSecret     string           `mapstrucutre:"JWT_SECRET"`
	SWTExpiresIn  int              `mapstrucutre:"SWT_EXPIRES_IN"`
	TokenAuth     *jwtauth.JWTAuth `mapstrucutre:"TOKEN_AUTH"`
}

func LoadConfig(path string) (*conf, error) {

	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, nil
}

```

## Carrega as configurações no arquivo main

```sh
# goot/cmd/server/main.go
package main

import (
	"fmt"

	"github.com/sallescosta/crud-api/configs"
)

func main() {
	config, _ := configs.LoadConfig(".")
	fmt.Println(config.DBDriver) //só para teste

}
```

# Create entity (user.go)
```sh
# root/internal/entity/user.go
package entity

import (
	"log"

	"github.com/sallescosta/crud-api/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

func NewUser(name, email, password string) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
```

# Test entity User

```sh
# root/internal/entity
package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "john@gmail.com", "123456")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@gmail.com", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("John Doe", "john@gmail.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("23456"))
	assert.NotEqual(t, "123456", user.Password)
}


# Para rodar, tem que navegar até a pasta onde esta este arquivo (root/internal/entity) e rodar `go test`
```
# Create Product entity and its test
```sh
# root/internal/entity/product.go
package entity

import (
	"errors"
	"time"

	"github.com/sallescosta/crud-api/pkg/entity"
)

var (
	ErrIDIsRequired    = errors.New("id is required")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidPrice    = errors.New("invalid price")
	ErrInvalidId       = errors.New("invalid id")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}

	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidId
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

func NewProduct(name string, price int) (*Product, error) {

	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil

}

```

## product entity unity test
```sh
# root/internal/entity/product_test.go
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
```
