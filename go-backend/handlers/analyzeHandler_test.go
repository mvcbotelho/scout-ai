package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mvcbotelho/scout-ai/models"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzePlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.GET("/analyze/players/:id", AnalyzePlayer(db))

	// Criar um jogador primeiro
	player := models.Player{
		Name:     "João Silva",
		Age:      25,
		Position: "Atacante",
		Team:     "Flamengo",
		Goals:    15,
		Tackles:  5,
		Passes:   120,
	}
	db.Create(&player)

	req, _ := http.NewRequest("GET", "/analyze/players/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response AnalysisResult
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, player.Name, response.PlayerName)
	assert.Equal(t, player.Position, response.Position)
	assert.True(t, response.Rating >= 1 && response.Rating <= 10)
	assert.NotEmpty(t, response.Analysis)
	assert.NotEmpty(t, response.Insights)
}

func TestAnalyzePlayerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.GET("/analyze/players/:id", AnalyzePlayer(db))

	req, _ := http.NewRequest("GET", "/analyze/players/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAnalyzeAllPlayers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.GET("/analyze/players", AnalyzeAllPlayers(db))

	// Criar alguns jogadores
	players := []models.Player{
		{
			Name:     "João Silva",
			Age:      25,
			Position: "Atacante",
			Team:     "Flamengo",
			Goals:    15,
			Tackles:  5,
			Passes:   120,
		},
		{
			Name:     "Pedro Santos",
			Age:      28,
			Position: "Meio-campo",
			Team:     "Palmeiras",
			Goals:    8,
			Tackles:  45,
			Passes:   350,
		},
	}

	for _, player := range players {
		db.Create(&player)
	}

	req, _ := http.NewRequest("GET", "/analyze/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "individual_analyses")
	assert.Contains(t, response, "comparative_analysis")

	individualAnalyses := response["individual_analyses"].([]interface{})
	assert.Len(t, individualAnalyses, 2)
}

func TestComparePlayers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.GET("/analyze/compare", ComparePlayers(db))

	// Criar alguns jogadores
	players := []models.Player{
		{
			Name:     "João Silva",
			Age:      25,
			Position: "Atacante",
			Team:     "Flamengo",
			Goals:    15,
			Tackles:  5,
			Passes:   120,
		},
		{
			Name:     "Pedro Santos",
			Age:      28,
			Position: "Meio-campo",
			Team:     "Palmeiras",
			Goals:    8,
			Tackles:  45,
			Passes:   350,
		},
	}

	for _, player := range players {
		db.Create(&player)
	}

	req, _ := http.NewRequest("GET", "/analyze/compare?ids=1&ids=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "players")
	assert.Contains(t, response, "comparison")

	playersList := response["players"].([]interface{})
	assert.Len(t, playersList, 2)
}

func TestComparePlayersInsufficient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.New()
	router.GET("/analyze/compare", ComparePlayers(db))

	req, _ := http.NewRequest("GET", "/analyze/compare?ids=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculatePlayerStats(t *testing.T) {
	player := models.Player{
		Name:     "João Silva",
		Age:      25,
		Position: "Atacante",
		Team:     "Flamengo",
		Goals:    15,
		Tackles:  5,
		Passes:   120,
	}

	stats := calculatePlayerStats(player)

	assert.Equal(t, player, stats.Player)
	assert.Greater(t, stats.Stats.Efficiency, 0.0)
	assert.NotEmpty(t, stats.Stats.PerformanceRank)
	assert.Greater(t, stats.Stats.GoalsPerGame, 0.0)
}

func TestGenerateInsights(t *testing.T) {
	player := models.Player{
		Name:     "João Silva",
		Age:      25,
		Position: "Atacante",
		Team:     "Flamengo",
		Goals:    15,
		Tackles:  5,
		Passes:   120,
	}

	stats := calculatePlayerStats(player)
	insights := generateInsights(player, stats)

	assert.NotEmpty(t, insights)
	assert.IsType(t, []string{}, insights)
}

func TestCalculateRating(t *testing.T) {
	player := models.Player{
		Name:     "João Silva",
		Age:      25,
		Position: "Atacante",
		Team:     "Flamengo",
		Goals:    15,
		Tackles:  5,
		Passes:   120,
	}

	stats := calculatePlayerStats(player)
	rating := calculateRating(player, stats)

	assert.GreaterOrEqual(t, rating, 1)
	assert.LessOrEqual(t, rating, 10)
}
