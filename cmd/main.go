package main

import (
	"fmt"
	"log"
	database "server/db"
	"server/internal/user"
	"server/router"
)

func main() {
	dbConn, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}
	defer dbConn.Close()
	fmt.Println("Database connection works!")

	err = dbConn.InitializeSchema()
	if err != nil {
		log.Fatalf("Wasnt able to create tables: %s", err)
	}
	fmt.Println("Tables created successfully")

	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	router.InitRouter(userHandler)
	router.Start(":8080")

}
