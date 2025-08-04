package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   Server   `mapstructure:"server" validate:"required"`
		Database Database `mapstructure:"database" validate:"required"`
		Redis    Redis    `mapstructure:"redis" validate:"required"`
		Jwt      Jwt      `mapstructure:"jwt" validate:"required"`
	}

	Server struct {
		Port         int           `mapstructure:"port" validate:"required"`
		AllowOrigins []string      `mapstructure:"allowOrigins" validate:"required"`
		BodyLimit    string        `mapstructure:"bodyLimit" validate:"required"`
		Timeout      time.Duration `mapstructure:"timeout" validate:"required"`
	}

	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     string `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		DBName   string `mapstructure:"dbname" validate:"required"`
	}

	Jwt struct {
		Secret string        `mapstructure:"secret" validate:"required"`
		TTL    time.Duration `mapstructure:"ttl" validate:"required"`
	}

	Redis struct {
		NameSpace string        `mapstructure:"nameSpace" validate:"required"`
		Read      RedisInstance `mapstructure:"read" validate:"required"`
		Write     RedisInstance `mapstructure:"write" validate:"required"`
		Pub       RedisInstance `mapstructure:"pub" validate:"required"`
	}

	RedisInstance struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     string `mapstructure:"port" validate:"required"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}
)

var (
	configInstance *Config
	once           sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  No .env file found (this is OK if running in container or CI)")
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("error reading config: %v", err)
		}

		for _, key := range viper.AllKeys() {
			val := viper.GetString(key)
			viper.Set(key, os.ExpandEnv(val))
		}

		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatalf("unable to unmarshal config: %v", err)
		}

		validate := validator.New()
		if err := validate.Struct(cfg); err != nil {
			log.Fatalf("config validation failed: %v", err)
		}

		configInstance = &cfg
		fmt.Println("✅ Config loaded successfully")
	})
	return configInstance
}
