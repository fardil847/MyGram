package main

import (
	"MyGram/controllers"
	"MyGram/database"
	"MyGram/routers"
)

func main() {
	db := database.InitDB()

	userRepository := controllers.UserRepository{DB: db}
	photoRepository := controllers.PhotoRepository{DB: db}
	commentRepository := controllers.CommentRepository{DB: db}
	socialmediaRepository := controllers.SocialMediaRepository{DB: db}

	r := routers.StartApp(userRepository, photoRepository, commentRepository, socialmediaRepository)
	r.Run(":8080")
}
