package main

import (
	"URL-Shortener/controller"
	"URL-Shortener/db"
	"URL-Shortener/service"
	"log"

	"github.com/gin-gonic/gin"
)

func chooseDB() *db.MongoRepository {
	mongourl := "<PASTE_URL_MONGODB_URL>"
	mongodb := "Shortener"
	timeout := 10

	repo, err := db.NewMongoRepo(mongourl, mongodb, timeout)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func main() {

	repo := chooseDB()
	service := service.NewRedirectService(repo)
	h := controller.NewHandler(service)

	router := gin.Default()

	router.GET("/url/:id", h.Find)
	router.POST("/url", h.Store)
	router.GET("/url", h.All)

	router.Run("localhost:8080")
}
