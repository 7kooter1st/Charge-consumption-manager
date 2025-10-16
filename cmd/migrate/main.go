package main

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.UCCP{},
		&ds.UseCase{},
		&ds.Consumption{},
		&ds.Users{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}

// package main

// import (
// 	"RIP/internal/app/ds"
// 	"RIP/internal/app/dsn"
// 	"log"

// 	"github.com/joho/godotenv"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func main() {
// 	_ = godotenv.Load()
// 	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	// Очистка базы данных (для разработки)
// 	sqlDB, err := db.DB()
// 	if err == nil {
// 		_, err = sqlDB.Exec("DROP TABLE IF EXISTS uccps, consumptions, use_cases, users CASCADE")
// 		if err != nil {
// 			log.Printf("Warning when dropping tables: %v", err)
// 		}
// 	}

// 	// Миграция в правильном порядке
// 	tables := []interface{}{
// 		&ds.Users{},
// 		&ds.UseCase{},
// 		&ds.Consumption{},
// 		&ds.UCCP{},
// 	}

// 	for _, table := range tables {
// 		err = db.AutoMigrate(table)
// 		if err != nil {
// 			log.Fatalf("Failed to migrate %T: %v", table, err)
// 		}
// 		log.Printf("Successfully migrated %T", table)
// 	}

// 	log.Println("Migration completed successfully")
// }
