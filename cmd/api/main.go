package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/nkyizbay/ticket_store/internal/user"
	"github.com/nkyizbay/ticket_store/pkg/database"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()

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

	// USER
	userRepository := user.NewRepository(connectionPool)
	userService := user.NewUserService(userRepository)
	user.NewHandler(e, userService, jwtSecretKey)

	e.Logger.Fatal(e.Start(":8080"))
}