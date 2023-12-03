package config

//читает конфигурацию
import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//config структура конфигурации
//содержит все конфигурационные данные о сервисе
//автоподгружаются при изменении исходного файла

type Config struct {
	ServiceHost string
	ServicePort int
	JWT         JWTConfig
}
type JWTConfig struct {
	Token         string
	ExpiresIn     time.Duration
	SigningMethod jwt.SigningMethod
}

// создает новый объект конфигурации, загружая данные из файла конфигурации
func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg) //json->struct
	if err != nil {
		return nil, err
	}
	log.Info("config parsed")

	cfg.JWT.Token = "test"
	cfg.JWT.ExpiresIn = time.Hour
	cfg.JWT.SigningMethod = jwt.SigningMethodHS256

	return cfg, nil
}
