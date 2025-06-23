package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mvcbotelho/scout-ai/models"
	"github.com/mvcbotelho/scout-ai/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	r := gin.Default()

	// Conecta no banco
	dsn := "host=db user=postgres password=postgres dbname=scoutdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	// Faz AutoMigrate
	db.AutoMigrate(&models.Player{})

	// Endpoints
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/players", handlers.CreatePlayer(db))
	r.GET("/players", handlers.GetPlayers(db))
	r.GET("/players/:id", handlers.GetPlayerByID(db))
	r.PUT("/players/:id", handlers.UpdatePlayer(db))
	r.DELETE("/players/:id", handlers.DeletePlayer(db))


	r.Run(":8080")
}
