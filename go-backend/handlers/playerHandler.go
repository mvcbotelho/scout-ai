package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mvcbotelho/scout-ai/models"
	"gorm.io/gorm"
)

func CreatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
			return
		}

		// Validação básica
		if player.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nome é obrigatório"})
			return
		}

		if player.Age <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idade deve ser maior que zero"})
			return
		}

		if err := db.Create(&player).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar jogador: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, player)
	}
}

func GetPlayers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var players []models.Player

		if err := db.Find(&players).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogadores: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, players)
	}
}

func GetPlayerByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player
		id := c.Param("id")

		// Validação do ID
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		if err := db.First(&player, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Jogador não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogador: " + err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, player)
	}
}

func UpdatePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player
		id := c.Param("id")

		// Validação do ID
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// Verifica se o jogador existe
		if err := db.First(&player, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Jogador não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogador: " + err.Error()})
			}
			return
		}

		var input models.Player
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
			return
		}

		// Validação básica
		if input.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nome é obrigatório"})
			return
		}

		if input.Age <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idade deve ser maior que zero"})
			return
		}

		// Atualiza apenas os campos fornecidos
		if err := db.Model(&player).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar jogador: " + err.Error()})
			return
		}

		// Busca o jogador atualizado
		db.First(&player, id)
		c.JSON(http.StatusOK, player)
	}
}

func DeletePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Validação do ID
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// Verifica se o jogador existe antes de deletar
		var player models.Player
		if err := db.First(&player, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Jogador não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogador: " + err.Error()})
			}
			return
		}

		if err := db.Delete(&player).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar jogador: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Jogador deletado com sucesso"})
	}
}
