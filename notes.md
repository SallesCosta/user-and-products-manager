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


