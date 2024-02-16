package routers

import (
	"MyGram/controllers"
	"MyGram/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp(c controllers.UserRepository, p controllers.PhotoRepository, o controllers.CommentRepository, m controllers.SocialMediaRepository) *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.GET("/", c.GetAllUser)
		userRouter.POST("/register", c.UserRegister)
		userRouter.POST("/login", c.UserLogin)
		userRouter.PUT("/:userId", middlewares.Authentication(), c.UserUpdate)
		userRouter.DELETE("/", middlewares.Authentication(), c.UserDelete)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.GET("/", p.GetPhoto)
		photoRouter.POST("/", p.UploadPhoto)
		photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(), p.UpdatePhoto)
		photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(), p.DeletePhoto)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("/", o.GetComment)
		commentRouter.POST("/", o.UploadComment)
		commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), o.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.CommentAuthorization(), o.DeleteComment)
	}

	socialmediaRouter := r.Group("/socialmedias")
	{
		socialmediaRouter.Use(middlewares.Authentication())
		socialmediaRouter.GET("/", m.GetSocialMedia)
		socialmediaRouter.POST("/", m.UploadSocialMedia)
		socialmediaRouter.PUT("/:socialMediaId", middlewares.SocialMediaAuthorization(), m.UpdateSocialMedia)
		socialmediaRouter.DELETE("/:socialMediaId", middlewares.SocialMediaAuthorization(), m.DeleteSocialMedia)

	}

	return r
}
