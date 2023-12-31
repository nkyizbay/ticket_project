package database

import (
	"fmt"

	"github.com/nkyizbay/ticket_store/internal/notification"
	"github.com/nkyizbay/ticket_store/internal/ticket"
	"github.com/nkyizbay/ticket_store/internal/trip"
	"github.com/nkyizbay/ticket_store/internal/user"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_PORT"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_DB"))

	var err error

	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate() {
	if err := db.AutoMigrate(&user.User{}, &trip.Trip{}, &ticket.Ticket{}, &notification.Log{}); err != nil {
		panic(err)
	}
}
