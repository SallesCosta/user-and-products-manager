package configs

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
