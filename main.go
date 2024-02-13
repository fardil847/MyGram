package main

import (
	"MyGram/controllers"
	"MyGram/database"
	"MyGram/routers"
)

func main() {
	db := database.InitDB()

	userRepository := controllers.UserRepository{DB: db}

	r := routers.StartApp(userRepository)
	r.Run(":8080")
}
