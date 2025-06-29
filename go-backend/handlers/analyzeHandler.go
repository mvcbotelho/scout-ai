package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mvcbotelho/scout-ai/models"
	"gorm.io/gorm"
)

// AnalysisResult representa o resultado da análise
type AnalysisResult struct {
	PlayerID   uint     `json:"player_id"`
	PlayerName string   `json:"player_name"`
	Analysis   string   `json:"analysis"`
	Insights   []string `json:"insights"`
	Rating     int      `json:"rating"` // 1-10
	Position   string   `json:"position"`
	Team       string   `json:"team"`
	AIUsed     bool     `json:"ai_used"`
}

// PlayerStats representa estatísticas do jogador
type PlayerStats struct {
	Player models.Player `json:"player"`
	Stats  struct {
		GoalsPerGame    float64 `json:"goals_per_game"`
		TacklesPerGame  float64 `json:"tackles_per_game"`
		PassesPerGame   float64 `json:"passes_per_game"`
		Efficiency      float64 `json:"efficiency"`
		PerformanceRank string  `json:"performance_rank"`
	} `json:"stats"`
}

// AnalyzePlayer analisa um jogador específico
func AnalyzePlayer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Validação do ID
		if _, err := strconv.Atoi(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var player models.Player
		if err := db.First(&player, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Jogador não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogador: " + err.Error()})
			}
			return
		}

		// Verificar se deve usar Ollama
		useAI := c.Query("ai") == "true" || c.Query("ai") == "1"

		// Gerar análise
		var analysis AnalysisResult
		if useAI {
			analysis = generatePlayerAnalysisWithAI(player)
		} else {
			analysis = generatePlayerAnalysis(player)
		}

		c.JSON(http.StatusOK, analysis)
	}
}

// AnalyzeAllPlayers analisa todos os jogadores
func AnalyzeAllPlayers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var players []models.Player
		if err := db.Find(&players).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogadores: " + err.Error()})
			return
		}

		// Verificar se deve usar Ollama
		useAI := c.Query("ai") == "true" || c.Query("ai") == "1"

		var analyses []AnalysisResult
		for _, player := range players {
			var analysis AnalysisResult
			if useAI {
				analysis = generatePlayerAnalysisWithAI(player)
			} else {
				analysis = generatePlayerAnalysis(player)
			}
			analyses = append(analyses, analysis)
		}

		// Gerar análise comparativa
		comparativeAnalysis := generateComparativeAnalysis(players)

		// Se usar AI, gerar análise comparativa com Ollama
		if useAI {
			if comparativeText, err := generateComparativeOllamaAnalysis(players, DefaultOllamaConfig); err == nil {
				comparativeAnalysis["ai_comparative_analysis"] = comparativeText
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"individual_analyses":  analyses,
			"comparative_analysis": comparativeAnalysis,
		})
	}
}

// ComparePlayers compara dois ou mais jogadores
func ComparePlayers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ids := c.QueryArray("ids")
		if len(ids) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "É necessário pelo menos 2 IDs de jogadores para comparação"})
			return
		}

		var players []models.Player
		if err := db.Where("id IN ?", ids).Find(&players).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar jogadores: " + err.Error()})
			return
		}

		if len(players) != len(ids) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Alguns jogadores não foram encontrados"})
			return
		}

		// Verificar se deve usar Ollama
		useAI := c.Query("ai") == "true" || c.Query("ai") == "1"

		comparison := generatePlayerComparison(players)

		// Se usar AI, adicionar análise comparativa com Ollama
		if useAI {
			if comparativeText, err := generateComparativeOllamaAnalysis(players, DefaultOllamaConfig); err == nil {
				comparison["ai_comparative_analysis"] = comparativeText
			}
		}

		c.JSON(http.StatusOK, comparison)
	}
}

// generatePlayerAnalysisWithAI gera análise usando Ollama
func generatePlayerAnalysisWithAI(player models.Player) AnalysisResult {
	// Calcular estatísticas
	stats := calculatePlayerStats(player)

	// Gerar análise com Ollama
	analysis, insights, err := generateOllamaAnalysis(player, stats, DefaultOllamaConfig)

	// Se houver erro com Ollama, usar análise estática como fallback
	if err != nil {
		analysis = generateAnalysisText(player, stats, generateInsights(player, stats))
		insights = generateInsights(player, stats)
	}

	// Calcular rating (1-10)
	rating := calculateRating(player, stats)

	return AnalysisResult{
		PlayerID:   player.ID,
		PlayerName: player.Name,
		Analysis:   analysis,
		Insights:   insights,
		Rating:     rating,
		Position:   player.Position,
		Team:       player.Team,
		AIUsed:     err == nil, // true se Ollama funcionou
	}
}

// generatePlayerAnalysis gera análise individual do jogador (versão estática)
func generatePlayerAnalysis(player models.Player) AnalysisResult {
	// Calcular estatísticas
	stats := calculatePlayerStats(player)

	// Gerar insights baseados nos dados
	insights := generateInsights(player, stats)

	// Gerar análise textual
	analysis := generateAnalysisText(player, stats, insights)

	// Calcular rating (1-10)
	rating := calculateRating(player, stats)

	return AnalysisResult{
		PlayerID:   player.ID,
		PlayerName: player.Name,
		Analysis:   analysis,
		Insights:   insights,
		Rating:     rating,
		Position:   player.Position,
		Team:       player.Team,
		AIUsed:     false,
	}
}

// calculatePlayerStats calcula estatísticas do jogador
func calculatePlayerStats(player models.Player) PlayerStats {
	stats := PlayerStats{Player: player}

	// Estatísticas básicas (assumindo 30 jogos por temporada)
	games := 30.0
	stats.Stats.GoalsPerGame = float64(player.Goals) / games
	stats.Stats.TacklesPerGame = float64(player.Tackles) / games
	stats.Stats.PassesPerGame = float64(player.Passes) / games

	// Eficiência baseada na posição
	switch strings.ToLower(player.Position) {
	case "atacante":
		stats.Stats.Efficiency = float64(player.Goals)*0.6 + float64(player.Passes)*0.4
	case "meio-campo", "meio-campista":
		stats.Stats.Efficiency = float64(player.Passes)*0.5 + float64(player.Tackles)*0.3 + float64(player.Goals)*0.2
	case "zagueiro":
		stats.Stats.Efficiency = float64(player.Tackles)*0.6 + float64(player.Passes)*0.4
	case "goleiro":
		stats.Stats.Efficiency = float64(player.Tackles)*0.8 + float64(player.Passes)*0.2
	default:
		stats.Stats.Efficiency = float64(player.Goals) + float64(player.Tackles) + float64(player.Passes)
	}

	// Ranking de performance
	if stats.Stats.Efficiency > 200 {
		stats.Stats.PerformanceRank = "Excelente"
	} else if stats.Stats.Efficiency > 150 {
		stats.Stats.PerformanceRank = "Muito Bom"
	} else if stats.Stats.Efficiency > 100 {
		stats.Stats.PerformanceRank = "Bom"
	} else if stats.Stats.Efficiency > 50 {
		stats.Stats.PerformanceRank = "Regular"
	} else {
		stats.Stats.PerformanceRank = "Abaixo da Média"
	}

	return stats
}

// generateInsights gera insights baseados nos dados (versão estática)
func generateInsights(player models.Player, stats PlayerStats) []string {
	var insights []string

	// Insights baseados na idade
	if player.Age < 23 {
		insights = append(insights, "Jogador jovem com potencial de desenvolvimento")
	} else if player.Age > 30 {
		insights = append(insights, "Jogador experiente, pode ser importante para liderança")
	}

	// Insights baseados na posição e estatísticas
	switch strings.ToLower(player.Position) {
	case "atacante":
		if player.Goals > 15 {
			insights = append(insights, "Artilheiro eficiente, excelente finalização")
		}
		if player.Passes > 100 {
			insights = append(insights, "Atacante que também participa da construção do jogo")
		}
	case "meio-campo", "meio-campista":
		if player.Passes > 200 {
			insights = append(insights, "Meio-campista com excelente visão de jogo")
		}
		if player.Tackles > 50 {
			insights = append(insights, "Meio-campista que também marca bem")
		}
	case "zagueiro":
		if player.Tackles > 80 {
			insights = append(insights, "Zagueiro com excelente marcação")
		}
		if player.Passes > 150 {
			insights = append(insights, "Zagueiro com boa saída de bola")
		}
	case "goleiro":
		if player.Tackles > 10 {
			insights = append(insights, "Goleiro com boa saída do gol")
		}
	}

	// Insights baseados no rating
	if stats.Stats.PerformanceRank == "Excelente" {
		insights = append(insights, "Performance excepcional na temporada")
	} else if stats.Stats.PerformanceRank == "Abaixo da Média" {
		insights = append(insights, "Necessita de melhoria no desempenho")
	}

	return insights
}

// generateAnalysisText gera análise textual do jogador (versão estática)
func generateAnalysisText(player models.Player, stats PlayerStats, insights []string) string {
	analysis := fmt.Sprintf("%s, %d anos, atua como %s no %s. ",
		player.Name, player.Age, player.Position, player.Team)

	analysis += fmt.Sprintf("Na temporada, marcou %d gols, realizou %d tackles e %d passes. ",
		player.Goals, player.Tackles, player.Passes)

	analysis += fmt.Sprintf("Sua eficiência geral é de %.1f pontos, classificando-o como %s. ",
		stats.Stats.Efficiency, stats.Stats.PerformanceRank)

	if len(insights) > 0 {
		analysis += "Principais características: " + strings.Join(insights, "; ") + ". "
	}

	// Recomendações baseadas na análise
	if stats.Stats.Efficiency > 150 {
		analysis += "Recomendado para times de alto nível."
	} else if stats.Stats.Efficiency > 100 {
		analysis += "Adequado para times de nível médio."
	} else {
		analysis += "Pode se beneficiar de mais tempo de desenvolvimento."
	}

	return analysis
}

// calculateRating calcula rating do jogador (1-10)
func calculateRating(player models.Player, stats PlayerStats) int {
	baseRating := 5

	// Ajustes baseados na eficiência
	if stats.Stats.Efficiency > 200 {
		baseRating += 3
	} else if stats.Stats.Efficiency > 150 {
		baseRating += 2
	} else if stats.Stats.Efficiency > 100 {
		baseRating += 1
	} else if stats.Stats.Efficiency < 50 {
		baseRating -= 1
	}

	// Ajustes baseados na idade
	if player.Age >= 25 && player.Age <= 30 {
		baseRating += 1 // Idade ideal
	} else if player.Age < 20 {
		baseRating -= 1 // Muito jovem
	}

	// Limitar entre 1 e 10
	if baseRating > 10 {
		baseRating = 10
	} else if baseRating < 1 {
		baseRating = 1
	}

	return baseRating
}

// generateComparativeAnalysis gera análise comparativa entre jogadores
func generateComparativeAnalysis(players []models.Player) map[string]interface{} {
	if len(players) == 0 {
		return map[string]interface{}{"message": "Nenhum jogador para análise"}
	}

	// Estatísticas gerais
	totalGoals := 0
	totalTackles := 0
	totalPasses := 0
	positions := make(map[string]int)

	for _, player := range players {
		totalGoals += player.Goals
		totalTackles += player.Tackles
		totalPasses += player.Passes
		positions[player.Position]++
	}

	// Encontrar melhores em cada categoria
	bestScorer := players[0]
	bestTackler := players[0]
	bestPasser := players[0]

	for _, player := range players {
		if player.Goals > bestScorer.Goals {
			bestScorer = player
		}
		if player.Tackles > bestTackler.Tackles {
			bestTackler = player
		}
		if player.Passes > bestPasser.Passes {
			bestPasser = player
		}
	}

	return map[string]interface{}{
		"total_players":          len(players),
		"total_goals":            totalGoals,
		"total_tackles":          totalTackles,
		"total_passes":           totalPasses,
		"positions_distribution": positions,
		"best_scorer": map[string]interface{}{
			"name":  bestScorer.Name,
			"goals": bestScorer.Goals,
			"team":  bestScorer.Team,
		},
		"best_tackler": map[string]interface{}{
			"name":    bestTackler.Name,
			"tackles": bestTackler.Tackles,
			"team":    bestTackler.Team,
		},
		"best_passer": map[string]interface{}{
			"name":   bestPasser.Name,
			"passes": bestPasser.Passes,
			"team":   bestPasser.Team,
		},
	}
}

// generatePlayerComparison compara jogadores específicos
func generatePlayerComparison(players []models.Player) map[string]interface{} {
	var analyses []AnalysisResult
	var comparison map[string]interface{}

	for _, player := range players {
		analysis := generatePlayerAnalysis(player)
		analyses = append(analyses, analysis)
	}

	// Comparação direta
	comparison = map[string]interface{}{
		"players": analyses,
		"comparison": map[string]interface{}{
			"highest_rating": analyses[0],
			"most_goals":     analyses[0],
			"most_tackles":   analyses[0],
			"most_passes":    analyses[0],
		},
	}

	// Encontrar melhores em cada categoria
	for _, analysis := range analyses {
		player := analysis.PlayerID
		for _, p := range players {
			if p.ID == player {
				if p.Goals > players[0].Goals {
					comparison["comparison"].(map[string]interface{})["most_goals"] = analysis
				}
				if p.Tackles > players[0].Tackles {
					comparison["comparison"].(map[string]interface{})["most_tackles"] = analysis
				}
				if p.Passes > players[0].Passes {
					comparison["comparison"].(map[string]interface{})["most_passes"] = analysis
				}
				break
			}
		}
	}

	return comparison
}
