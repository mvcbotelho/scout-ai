package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mvcbotelho/scout-ai/handlers"
	"github.com/mvcbotelho/scout-ai/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// Conecta no banco usando variáveis de ambiente
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "scoutdb")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	// Faz AutoMigrate
	if err := db.AutoMigrate(&models.Player{}); err != nil {
		log.Fatal("Erro ao fazer migração:", err)
	}

	// Endpoints
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/players", handlers.CreatePlayer(db))
	r.GET("/players", handlers.GetPlayers(db))
	r.GET("/players/:id", handlers.GetPlayerByID(db))
	r.PUT("/players/:id", handlers.UpdatePlayer(db))
	r.DELETE("/players/:id", handlers.DeletePlayer(db))

	log.Println("Servidor iniciado na porta 8080")
	r.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
