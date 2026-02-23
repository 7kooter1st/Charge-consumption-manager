package dsn

import (
	"fmt"
	"os"
)

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// FromEnv возвращает DSN для PostgreSQL.
// Значения по умолчанию соответствуют docker-compose (consumption_db, postgrespassword).
func FromEnv() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASS", "postgrespassword")
	dbname := getEnv("DB_NAME", "consumption_db")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
