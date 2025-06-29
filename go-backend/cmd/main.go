package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	// Configurar Ollama
	configureOllama()

	// Endpoints básicos
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Endpoints de jogadores
	r.POST("/players", handlers.CreatePlayer(db))
	r.GET("/players", handlers.GetPlayers(db))
	r.GET("/players/:id", handlers.GetPlayerByID(db))
	r.PUT("/players/:id", handlers.UpdatePlayer(db))
	r.DELETE("/players/:id", handlers.DeletePlayer(db))

	// Endpoints de análise
	r.GET("/analyze/players/:id", handlers.AnalyzePlayer(db))
	r.GET("/analyze/players", handlers.AnalyzeAllPlayers(db))
	r.GET("/analyze/compare", handlers.ComparePlayers(db))

	log.Println("Servidor iniciado na porta 8080")
	log.Println("Ollama configurado:", handlers.DefaultOllamaConfig.BaseURL)
	r.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func configureOllama() {
	// Configurar Ollama com variáveis de ambiente
	handlers.DefaultOllamaConfig.BaseURL = getEnv("OLLAMA_BASE_URL", "http://localhost:11434")
	handlers.DefaultOllamaConfig.Model = getEnv("OLLAMA_MODEL", "llama3.2")

	if temp := getEnv("OLLAMA_TEMPERATURE", "0.7"); temp != "" {
		if temperature, err := strconv.ParseFloat(temp, 64); err == nil {
			handlers.DefaultOllamaConfig.Temperature = temperature
		}
	}

	if topP := getEnv("OLLAMA_TOP_P", "0.9"); topP != "" {
		if topPValue, err := strconv.ParseFloat(topP, 64); err == nil {
			handlers.DefaultOllamaConfig.TopP = topPValue
		}
	}
}
