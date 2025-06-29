package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mvcbotelho/scout-ai/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Player{})
	return db
}

func TestCreatePlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.POST("/players", CreatePlayer(db))

	player := models.Player{
		Name:     "João Silva",
		Age:      25,
		Position: "Atacante",
		Team:     "Flamengo",
		Goals:    15,
		Tackles:  5,
		Passes:   120,
	}

	jsonData, _ := json.Marshal(player)
	req, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Player
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, player.Name, response.Name)
	assert.Equal(t, player.Age, response.Age)
}

func TestCreatePlayerInvalidData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.POST("/players", CreatePlayer(db))

	// Teste com dados inválidos
	player := models.Player{
		Name: "", // Nome vazio
		Age:  25,
	}

	jsonData, _ := json.Marshal(player)
	req, _ := http.NewRequest("POST", "/players", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlayers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.GET("/players", GetPlayers(db))

	// Criar um jogador primeiro
	player := models.Player{
		Name:     "João Silva",
		Age:      25,
		Position: "Atacante",
		Team:     "Flamengo",
	}
	db.Create(&player)

	req, _ := http.NewRequest("GET", "/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var players []models.Player
	json.Unmarshal(w.Body.Bytes(), &players)
	assert.Len(t, players, 1)
	assert.Equal(t, player.Name, players[0].Name)
}
