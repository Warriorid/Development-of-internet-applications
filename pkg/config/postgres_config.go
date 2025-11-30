package config

import (
	"DIA/internal/model"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresConfDB() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func MigrateBD() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(PostgresConfDB()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.Migrator().DropTable(
		&model.CalculationMaterial{},
		&model.PitsCalculation{}, 
		&model.Material{},
		&model.Users{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}