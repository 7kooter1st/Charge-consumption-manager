// package config

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/joho/godotenv"
// 	log "github.com/sirupsen/logrus" // Используем Logrus, как у вас
// 	"github.com/spf13/viper"
// )

// type JWTConfig struct {
// 	Secret           string
// 	ExpiresIn        time.Duration
// 	RefreshExpiresIn time.Duration // <-- Время жизни Refresh Token
// }

// type RedisConfig struct {
// 	Host        string
// 	Password    string
// 	Port        int
// 	User        string
// 	DialTimeout time.Duration
// 	ReadTimeout time.Duration
// }

// type Config struct {
// 	ServiceHost string
// 	ServicePort int
// 	JWT         JWTConfig
// 	Redis       RedisConfig // <-- Добавляем конфиг Redis

// 	// Поля БД будем загружать из .env
// 	DBHost     string
// 	DBPort     string
// 	DBUser     string
// 	DBPassword string
// 	DBName     string
// 	DBSSLMode  string
// }

// type MinioConfig struct {
// 	Host string
// 	Port string
// 	User string
// 	Pass string
// }

// const (
// 	envRedisHost = "REDIS_HOST"
// 	envRedisPort = "REDIS_PORT"
// 	envRedisUser = "REDIS_USER"
// 	envRedisPass = "REDIS_PASSWORD"
// )

// // func NewConfig() (*Config, error) {
// // 	var err error

// // 	configName := "config"
// // 	// Загружаем .env файл для настроек БД
// // 	err = godotenv.Load()
// // 	if err != nil {
// // 		log.Warn("No .env file found")
// // 	}

// // 	viper.SetConfigName(configName)
// // 	viper.SetConfigType("toml")
// // 	viper.AddConfigPath("config")
// // 	viper.AddConfigPath(".")
// // 	viper.WatchConfig()

// // 	err = viper.ReadInConfig()
// // 	if err != nil {
// // 		log.Warnf("Config file not found: %v", err)
// // 	}

// // 	cfg := &Config{}
// // 	err = viper.Unmarshal(cfg)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// Загружаем настройки БД из .env
// // 	cfg.DBHost = getEnv("DB_HOST", "localhost")
// // 	cfg.DBPort = getEnv("DB_PORT", "5432")
// // 	cfg.DBUser = getEnv("DB_USER", "postgres")
// // 	cfg.DBPassword = getEnv("DB_PASSWORD", "postgres1234")
// // 	cfg.DBName = getEnv("DB_NAME", "consumption-manager")
// // 	cfg.DBSSLMode = getEnv("DB_SSLMODE", "disable")
// // 	// cfg.Redis.Host = os.Getenv(envRedisHost)
// // 	// cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))

// // 	// if err != nil {
// // 	// 	return nil, fmt.Errorf("redis port must be int value: %w", err)
// // 	// }
// // 	// п
// // 	// 	cfg.Redis.Password = os.Getenv(envRedisPass)
// // 	// 	cfg.Redis.User = os.Getenv(envRedisUser)

// // 	// ОТЛАДКА
// // 	log.Infof("Service: %s:%d", cfg.ServiceHost, cfg.ServicePort)
// // 	log.Infof("Database: %s:%s, user: %s, db: %s",
// // 		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

// // 	log.Info("config parsed")

// // 	return cfg, nil
// // }

// func NewConfig() (*Config, error) {
// 	// Загружаем .env файл. Ошибки здесь нет, если файла нет - используются переменные ОС.
// 	if err := godotenv.Load(); err != nil {
// 		log.Warn("No .env file found, using OS environment variables")
// 	}

// 	viper.SetConfigName("config")
// 	viper.SetConfigType("toml")
// 	viper.AddConfigPath(".") // Искать config.toml в корне проекта
// 	if err := viper.ReadInConfig(); err != nil {
// 		return nil, fmt.Errorf("failed to read config.toml: %w", err)
// 	}

// 	var cfg Config
// 	if err := viper.Unmarshal(&cfg); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
// 	}

// 	// Читаем чувствительные данные из переменных окружения (.env или ОС)
// 	cfg.Redis.Host = os.Getenv("REDIS_HOST")
// 	redisPortStr := os.Getenv("REDIS_PORT")
// 	if redisPortStr != "" {
// 		port, err := strconv.Atoi(redisPortStr)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
// 		}
// 		cfg.Redis.Port = port
// 	}
// 	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

// 	cfg.Minio.Host = os.Getenv("MINIO_HOST")
// 	cfg.Minio.Port = os.Getenv("MINIO_PORT")
// 	cfg.Minio.User = os.Getenv("MINIO_USER")
// 	cfg.Minio.Pass = os.Getenv("MINIO_PASS")

// 	log.Info("Config parsed successfully")
// 	log.Infof("Service will run on: %s:%d", cfg.ServiceHost, cfg.ServicePort)

// 	return &cfg, nil
// }

// func getEnv(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus" // Используем Logrus, как у вас
	"github.com/spf13/viper"
)

type MinioConfig struct {
	Host      string
	Port      string
	User      string
	Pass      string
	Bucket    string // Добавили бакет
	PublicUrl string // Добавили публичный адрес
}

// RedisConfig хранит настройки для Redis.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

// JWTConfig хранит настройки для JWT.
type JWTConfig struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

// Config - главная структура, объединяющая все конфигурации.
type Config struct {
	ServiceHost string
	ServicePort int
	JWT         JWTConfig
	Redis       RedisConfig
	Minio       MinioConfig
}

// NewConfig - ЕДИНСТВЕННАЯ функция, которая читает все конфиги.
func NewConfig() (*Config, error) {
	// Загружаем .env файл. Ошибки здесь нет, если файла нет - используются переменные ОС.
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, using OS environment variables")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".") // Искать config.toml в корне проекта
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config.toml: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Читаем чувствительные данные из переменных окружения (.env или ОС)
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	redisPortStr := os.Getenv("REDIS_PORT")
	if redisPortStr != "" {
		port, err := strconv.Atoi(redisPortStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
		}
		cfg.Redis.Port = port
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	cfg.Minio.Host = os.Getenv("MINIO_HOST")
	cfg.Minio.Port = os.Getenv("MINIO_PORT")
	cfg.Minio.User = os.Getenv("MINIO_USER")
	cfg.Minio.Pass = os.Getenv("MINIO_PASS")

	log.Info("Config parsed successfully")
	log.Infof("Service will run on: %s:%d", cfg.ServiceHost, cfg.ServicePort)

	return &cfg, nil
}
