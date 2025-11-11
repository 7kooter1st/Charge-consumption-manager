// package config

// import (
// 	"os"

// 	"github.com/joho/godotenv"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	ServiceHost string
// 	ServicePort int
// }

// func NewConfig() (*Config, error) {
// 	var err error

// 	configName := "config"
// 	_ = godotenv.Load()
// 	if os.Getenv("CONFIG_NAME") != "" {
// 		configName = os.Getenv("CONFIG_NAME")
// 	}

// 	viper.SetConfigName(configName)
// 	viper.SetConfigType("toml")
// 	viper.AddConfigPath("config")
// 	viper.AddConfigPath(".")
// 	viper.WatchConfig()

// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	cfg := &Config{}           // создаем объект конфига
// 	err = viper.Unmarshal(cfg) // читаем информацию из файла,
// 	// конвертируем и затем кладем в нашу переменную cfg
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Info("config parsed")

// 	return cfg, nil
// }

package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string `mapstructure:"ServiceHost"`
	ServicePort int    `mapstructure:"ServicePort"`

	// Поля БД будем загружать из .env
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	// Загружаем .env файл для настроек БД
	err = godotenv.Load()
	if err != nil {
		log.Warn("No .env file found")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		log.Warnf("Config file not found: %v", err)
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	// Загружаем настройки БД из .env
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "5432")
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPassword = getEnv("DB_PASSWORD", "postgres1234")
	cfg.DBName = getEnv("DB_NAME", "consumption-manager")
	cfg.DBSSLMode = getEnv("DB_SSLMODE", "disable")

	// ОТЛАДКА
	log.Infof("Service: %s:%d", cfg.ServiceHost, cfg.ServicePort)
	log.Infof("Database: %s:%s, user: %s, db: %s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

	log.Info("config parsed")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
