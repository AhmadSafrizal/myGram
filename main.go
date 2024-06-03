package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AhmadSafrizal/myGram/handler"
	"github.com/AhmadSafrizal/myGram/middleware"
	"github.com/AhmadSafrizal/myGram/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	engine := gin.New()

	// Add logger and recovery middleware
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// Add timeout middleware
	engine.Use(middleware.Timeout(60 * time.Second))

	engine.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]any{
			"message": "This is a test",
		})
	})

	myDb := "host=localhost user=postgres password=postgre dbname=mygram port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(myDb), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	socialMediaRepo := &repository.SocialMediaRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}
	commentRepo := &repository.CommentRepository{DB: db}
	photoRepo := &repository.PhotoRepository{DB: db}

	socialMediaRepo.Migrate()
	userRepo.Migrate()
	commentRepo.Migrate()
	photoRepo.Migrate()

	socialMediaHandler := handler.NewSocialMediaHandler(socialMediaRepo)
	userHandler := handler.NewUserHandler(userRepo)
	photoHandler := handler.NewPhotoHandler(photoRepo)
	commentHandler := handler.NewCommentHandler(commentRepo)

	userGroup := engine.Group("/users")
	{
		userGroup.GET("", userHandler.GetGorm)
		userGroup.POST("/register", userHandler.CreateGorm)
		userGroup.POST("/login", userHandler.Login)

		userGroup.Use(middleware.Authotization())
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}

	photoGroup := engine.Group("/photos")
	{
		photoGroup.Use(middleware.Authotization())
		photoGroup.GET("", photoHandler.GetPhotos)
		photoGroup.POST("", photoHandler.AddPhoto)
		photoGroup.PUT("/:id", photoHandler.UpdatePhoto)
		photoGroup.DELETE("/:id", photoHandler.DeletePhoto)
	}

	socialMediaGroup := engine.Group("/socialmedias")
	{
		socialMediaGroup.Use(middleware.Authotization())
		socialMediaGroup.GET("", socialMediaHandler.GetSocialMedias)
		socialMediaGroup.POST("", socialMediaHandler.CreateSocialMedia)
		socialMediaGroup.PUT("/:id", socialMediaHandler.UpdateSocialMedia)
		socialMediaGroup.DELETE("/:id", socialMediaHandler.DeleteSocialMedia)
	}

	commentGroup := engine.Group("/comments")
	{
		commentGroup.Use(middleware.Authotization())
		commentGroup.GET("", commentHandler.GetComments)
		commentGroup.POST("", commentHandler.CreateComment)
		commentGroup.PUT("/:id", commentHandler.UpdateComment)
		commentGroup.DELETE("/:id", commentHandler.DeleteComment)
	}

	err = engine.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
