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
	JWTExpiresIn  int              `mapstrucutre:"JWT_EXPIRES_IN"`
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
# root/cmd/server/main.go
package main

import (
	"fmt"

	"github.com/sallescosta/user-and-products-manager/configs"
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

	"github.com/sallescosta/user-and-products-manager/pkg/entity"
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

	"github.com/sallescosta/user-and-products-manager/pkg/entity"
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
	Price     float64   `json:"price"`
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

# Criação das entities no db

```sh
# root/internal/infra/database/interface.go
package database

import "github.com/sallescosta/user-and-products-manager/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindById(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
}

```

```sh
# root/internal/infra/database/user_db.go
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
```

```sh
# root/internal/infra/database/user_db_test.go
package database

import (
	"testing"

	"github.com/sallescosta/user-and-products-manager/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "john@gmail.com", "123456")

	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, userFound.ID, user.ID)
	assert.Equal(t, userFound.Name, user.Name)
	assert.Equal(t, userFound.Email, user.Email)
	assert.NotNil(t, userFound.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "john@gmail.com", "123456")

	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.Equal(t, userFound.ID, user.ID)
	assert.Equal(t, userFound.Name, user.Name)
	assert.Equal(t, userFound.Email, user.Email)
	assert.NotNil(t, userFound.Password)
}
```

```sh
# root/internal/infra/database/product_db.go
package database

import (
	"github.com/sallescosta/user-and-products-manager/internal/entity"
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

func (p *Product) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		err = p.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	} else {
		err = p.DB.Order("created_at " + sort).Find(&products).Error
	}

	return products, err
}
```

## Integration tests for product_db

```sh
# root/internal/infra/database/product_db_test.go

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/sallescosta/user-and-products-manager/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	name    = "Product 1"
	price   = 10.34
	perPage = 10
)

func NewTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	if err := db.AutoMigrate(&entity.Product{}); err != nil {
		t.Error(err)
	}

	return db
}

func TestDeleteProduct(t *testing.T) {
	db := NewTestDB(t)

	product, err := entity.NewProduct(name, price)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	_, err = productDB.FindById(product.ID.String())
	assert.Error(t, err)
}

func TestCreateNewProduct(t *testing.T) {
	db := NewTestDB(t)

	product, err := entity.NewProduct(name, price)
	assert.Nil(t, err)

	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T) {
	db := NewTestDB(t)

	var precoSorteado = rand.Float64() * 100

	for i := 1; i < 33; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Produto %d", i), precoSorteado)
		assert.NoError(t, err)
		db.Create(product)
	}

	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, perPage, "asc")
	assert.NoError(t, err)

	assert.Len(t, products, perPage)
	assert.Equal(t, "Produto 1", products[0].Name)
	assert.Equal(t, "Produto 10", products[9].Name)

	products, err = productDB.FindAll(2, perPage, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, perPage)
	assert.Equal(t, "Produto 11", products[0].Name)
	assert.Equal(t, "Produto 20", products[9].Name)

	products, err = productDB.FindAll(3, perPage, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, perPage)
	assert.Equal(t, "Produto 21", products[0].Name)
	assert.Equal(t, "Produto 30", products[9].Name)

	products, err = productDB.FindAll(4, perPage, "asc")
	println(products)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Produto 31", products[0].Name)
	assert.Equal(t, "Produto 32", products[1].Name)
}

func TestFindProductByID(t *testing.T) {
	db := NewTestDB(t)

	product, err := entity.NewProduct(name, price)
	assert.NoError(t, err)

	db.Create(product)
	productDB := NewProduct(db)

	product, err = productDB.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, name, product.Name)
	assert.Equal(t, price, product.Price)

}

func TestUpdateProduct(t *testing.T) {
	db := NewTestDB(t)

	product, err := entity.NewProduct(name, price)
	assert.NoError(t, err)

	db.Create(product)
	productDB := NewProduct(db)

	product.Name = "Produto 2"
	err = productDB.Update(product)
	assert.NoError(t, err)

	product, err = productDB.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Produto 2", product.Name)
}

```

# Criando Handlers (Produtos)

- em `main.go` já deve haver a chamada para o método 
`LoadConfig` e a criação do banco de
dados (com ORM fazendo a migração das entities):

```sh
func main() {
	config, _ := configs.LoadConfig(".")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})
}

```

## Criação de DTOs (Não tem dependencia nenhuma..)

```sh
# root/internal/dto/dto.go

package dto

type CreateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetJWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}

type GetUsersOutput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

```


```sh
# root/internal/infra/webserver/handlers/product_handlers.go
type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)  // isso normalmente não é bom. Não é normal que o handler saiba como criar uma entidade.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
// para testar, independente de roteador ou framework, é só usar um:
// http.handleFunc("/products", h.CreateProduct)
// http.ListenAndServe(":8000", nil)
// e chmar num POST
// para verificar no db (sqlite) se está ok, abre um outro terminal
// e roda `sqlite3 cmd/server/test.db` e depois `select * from products;`
```

# Implementação do roteador (Chi)

- com o roteador e os handlers, o código do `main.go` fica assim:

```sh
func main() {
	config, _ := configs.LoadConfig(".")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productHandler := handlers.NewProductHandler(database.NewProduct(db))
	userHandler := handlers.NewUserHandler(database.NewUser(db))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", config.JWTExpiresIn))

	r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/", productHandler.GetProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)
	r.Get("/users", userHandler.AllUsers)

	r.Route("/users", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
	})

	http.HandleFunc("/products", productHandler.CreateProduct)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	log.Fatal(http.ListenAndServe(":8000", r))
}
```

# Gerando JWT

- os jwt estão associados aos usuários, então deve ter jwt na interface do usuário
```sh
# root/internal/infra/webserver/handlers/user_handlers.go

import "github.com/go-chi/jwtauth"

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

```

- Package publico de geração de JWT

```sh
# root/pkg/jwt/jwt.go
package entity

import (
	"github.com/google/uuid"
)

type ID = uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}
```

- Criar o handler GetJwt
// tem que ter o dto pronto (dto.GetJWTInput)
// 

```sh
# root/internal/infra/webserver/handlers/user_handlers.go
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

// serializa o body
	var user dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verifica se o usuário existe
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
	}

	// se chegou aqui, o usuário existe, então verifica a senha
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m := map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	}
	_, tokenString, _ := jwt.Encode(m)

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)
	if err != nil {
		return
	}
}



```
