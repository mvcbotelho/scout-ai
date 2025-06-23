package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"github.com/mvcbotelho/scout-ai/models"
)

func CreatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&player).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create player"})
			return
		}

		c.JSON(http.StatusCreated, player)
	}
}

func GetPlayers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var players []models.Player

		if err := db.Find(&players).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve players"})
			return
		}

		c.JSON(http.StatusOK, players)
	}
}

func GetPlayerByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player
		id := c.Param("id")

		if err := db.First(&player, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}

		c.JSON(http.StatusOK, player)
	}
}

func UpdatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player
		id := c.Param("id")

		if err := db.First(&player, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}

		var input models.Player
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Model(&player).Updates(input)
		c.JSON(http.StatusOK, player)
	}
}

func DeletePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Delete(&models.Player{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete player"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "player deleted"})
	}
}
