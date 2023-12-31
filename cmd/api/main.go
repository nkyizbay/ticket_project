package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/nkyizbay/ticket_store/internal/auth"
	"github.com/nkyizbay/ticket_store/internal/notification"
	"github.com/nkyizbay/ticket_store/internal/ticket"
	"github.com/nkyizbay/ticket_store/internal/trip"
	"github.com/nkyizbay/ticket_store/internal/user"
	"github.com/nkyizbay/ticket_store/pkg/database"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()
	e.Use(auth.TokenMiddleware)

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	fmt.Println(viper.Get("POSTGRES_DB"))

	jwtSecretKey := viper.GetString("ONLINE_TICKET_GO_JWTKEY")

	connectionPool, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}

	database.Migrate()

	// NOTIFICATION
	notificationRepository := notification.NewNotificationRepository(connectionPool)
	notificationService := notification.NewService(notificationRepository)

	// USER
	userRepository := user.NewRepository(connectionPool)
	userService := user.NewUserService(userRepository)
	user.NewHandler(e, userService, notificationService, jwtSecretKey)

	// TRİP
	tripRepo := trip.NewTripRepository(connectionPool)
	tripService := trip.NewTripService(tripRepo)
	trip.Handler(e, tripService)

	// TICKET
	ticketRepo := ticket.NewTicketRepository(connectionPool)
	service := ticket.NewService(ticketRepo, notificationService, tripRepo)
	ticket.NewHandler(e, service)

	e.Logger.Fatal(e.Start(":8080"))
}
