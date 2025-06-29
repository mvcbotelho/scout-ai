package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mvcbotelho/scout-ai/models"
)

// OllamaRequest representa a requisição para o Ollama
type OllamaRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Options struct {
		Temperature float64 `json:"temperature"`
		TopP        float64 `json:"top_p"`
	} `json:"options"`
}

// OllamaResponse representa a resposta do Ollama
type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// OllamaConfig configuração do Ollama
type OllamaConfig struct {
	BaseURL     string
	Model       string
	Temperature float64
	TopP        float64
}

// DefaultOllamaConfig configuração padrão
var DefaultOllamaConfig = OllamaConfig{
	BaseURL:     "http://localhost:11434",
	Model:       "llama3.2",
	Temperature: 0.7,
	TopP:        0.9,
}

// generateOllamaAnalysis gera análise usando Ollama
func generateOllamaAnalysis(player models.Player, stats PlayerStats, config OllamaConfig) (string, []string, error) {
	// Criar prompt para o Ollama
	prompt := createAnalysisPrompt(player, stats)

	// Fazer requisição para o Ollama
	analysis, err := callOllama(prompt, config)
	if err != nil {
		return "", nil, fmt.Errorf("erro ao chamar Ollama: %v", err)
	}

	// Extrair insights da análise
	insights := extractInsightsFromAnalysis(analysis)

	return analysis, insights, nil
}

// createAnalysisPrompt cria o prompt para o Ollama
func createAnalysisPrompt(player models.Player, stats PlayerStats) string {
	positionAnalysis := getPositionAnalysis(player.Position, player, stats)

	prompt := fmt.Sprintf(`Você é um scout de futebol experiente. Analise o seguinte jogador e forneça uma análise detalhada em português brasileiro.

DADOS DO JOGADOR:
- Nome: %s
- Idade: %d anos
- Posição: %s
- Time: %s
- Gols: %d
- Tackles: %d
- Passes: %d
- Eficiência: %.1f pontos
- Performance: %s

ANÁLISE ESPECÍFICA DA POSIÇÃO:
%s

INSTRUÇÕES:
1. Forneça uma análise completa e profissional
2. Use linguagem técnica de futebol
3. Identifique pontos fortes e áreas de melhoria
4. Sugira adequação para diferentes níveis de competição
5. Considere a idade do jogador no contexto
6. Seja específico sobre as estatísticas
7. Forneça recomendações práticas

Responda em português brasileiro com tom profissional de scout.`,
		player.Name, player.Age, player.Position, player.Team,
		player.Goals, player.Tackles, player.Passes,
		stats.Stats.Efficiency, stats.Stats.PerformanceRank,
		positionAnalysis)

	return prompt
}

// getPositionAnalysis retorna análise específica da posição
func getPositionAnalysis(position string, player models.Player, stats PlayerStats) string {
	switch strings.ToLower(position) {
	case "atacante":
		return fmt.Sprintf(`ANÁLISE DE ATAQUE:
- Gols por jogo: %.2f
- Participação no jogo: %d passes
- Eficiência de finalização: %.1f%%
- Comparação com padrões da posição: %s`,
			stats.Stats.GoalsPerGame, player.Passes,
			float64(player.Goals)/float64(player.Passes)*100,
			getPerformanceComparison(player.Goals, 15))

	case "meio-campo", "meio-campista":
		return fmt.Sprintf(`ANÁLISE DE MEIO-CAMPO:
- Passes por jogo: %.2f
- Marcação: %d tackles
- Participação ofensiva: %d gols
- Visão de jogo: %s`,
			stats.Stats.PassesPerGame, player.Tackles, player.Goals,
			getPassingQuality(player.Passes))

	case "zagueiro":
		return fmt.Sprintf(`ANÁLISE DEFENSIVA:
- Tackles por jogo: %.2f
- Saída de bola: %d passes
- Participação ofensiva: %d gols
- Eficiência defensiva: %s`,
			stats.Stats.TacklesPerGame, player.Passes, player.Goals,
			getDefensiveQuality(player.Tackles))

	case "goleiro":
		return fmt.Sprintf(`ANÁLISE DE GOLEIRO:
- Saída do gol: %d tackles
- Participação no jogo: %d passes
- Comando de área: %s`,
			player.Tackles, player.Passes,
			getGoalkeeperQuality(player.Tackles))

	default:
		return "Posição não especificada - análise geral baseada em estatísticas."
	}
}

// getPerformanceComparison compara performance com padrões
func getPerformanceComparison(goals, threshold int) string {
	if goals >= threshold {
		return "Acima da média para a posição"
	} else if goals >= threshold/2 {
		return "Na média para a posição"
	} else {
		return "Abaixo da média para a posição"
	}
}

// getPassingQuality avalia qualidade dos passes
func getPassingQuality(passes int) string {
	if passes > 300 {
		return "Excelente visão de jogo e distribuição"
	} else if passes > 200 {
		return "Boa capacidade de distribuição"
	} else if passes > 100 {
		return "Capacidade de distribuição adequada"
	} else {
		return "Necessita melhorar distribuição"
	}
}

// getDefensiveQuality avalia qualidade defensiva
func getDefensiveQuality(tackles int) string {
	if tackles > 100 {
		return "Excelente marcação e recuperação"
	} else if tackles > 70 {
		return "Boa capacidade defensiva"
	} else if tackles > 40 {
		return "Capacidade defensiva adequada"
	} else {
		return "Necessita melhorar marcação"
	}
}

// getGoalkeeperQuality avalia qualidade do goleiro
func getGoalkeeperQuality(tackles int) string {
	if tackles > 20 {
		return "Excelente saída do gol e comando"
	} else if tackles > 10 {
		return "Boa saída do gol"
	} else {
		return "Necessita melhorar saída do gol"
	}
}

// callOllama faz a chamada para o Ollama
func callOllama(prompt string, config OllamaConfig) (string, error) {
	requestBody := OllamaRequest{
		Model:  config.Model,
		Prompt: prompt,
		Stream: false,
	}
	requestBody.Options.Temperature = config.Temperature
	requestBody.Options.TopP = config.TopP

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar requisição: %v", err)
	}

	url := fmt.Sprintf("%s/api/generate", config.BaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro na requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro do servidor Ollama (status %d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	return ollamaResp.Response, nil
}

// extractInsightsFromAnalysis extrai insights da análise do Ollama
func extractInsightsFromAnalysis(analysis string) []string {
	var insights []string

	// Dividir a análise em frases
	sentences := strings.Split(analysis, ". ")

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// Identificar frases que contêm insights importantes
		lowerSentence := strings.ToLower(sentence)
		if strings.Contains(lowerSentence, "excelente") ||
			strings.Contains(lowerSentence, "forte") ||
			strings.Contains(lowerSentence, "bom") ||
			strings.Contains(lowerSentence, "melhorar") ||
			strings.Contains(lowerSentence, "potencial") ||
			strings.Contains(lowerSentence, "recomendado") ||
			strings.Contains(lowerSentence, "adequado") {
			insights = append(insights, sentence)
		}
	}

	// Limitar a 5 insights mais relevantes
	if len(insights) > 5 {
		insights = insights[:5]
	}

	return insights
}

// generateComparativeOllamaAnalysis gera análise comparativa usando Ollama
func generateComparativeOllamaAnalysis(players []models.Player, config OllamaConfig) (string, error) {
	if len(players) == 0 {
		return "Nenhum jogador para análise comparativa.", nil
	}

	// Criar resumo dos jogadores
	var playerSummaries []string
	for i, player := range players {
		stats := calculatePlayerStats(player)
		summary := fmt.Sprintf("%d. %s (%s, %s): %d gols, %d tackles, %d passes, eficiência %.1f",
			i+1, player.Name, player.Position, player.Team,
			player.Goals, player.Tackles, player.Passes, stats.Stats.Efficiency)
		playerSummaries = append(playerSummaries, summary)
	}

	prompt := fmt.Sprintf(`Você é um scout de futebol experiente. Analise e compare os seguintes jogadores, fornecendo uma análise comparativa detalhada em português brasileiro.

JOGADORES PARA COMPARAÇÃO:
%s

INSTRUÇÕES:
1. Compare os jogadores em diferentes aspectos (técnico, físico, mental)
2. Identifique o melhor em cada categoria (ataque, meio-campo, defesa)
3. Sugira qual jogador seria mais adequado para diferentes tipos de time
4. Considere a idade e experiência de cada jogador
5. Forneça recomendações específicas para cada jogador
6. Use linguagem técnica de futebol
7. Seja objetivo e baseado nos dados fornecidos

Responda em português brasileiro com tom profissional de scout.`, strings.Join(playerSummaries, "\n"))

	return callOllama(prompt, config)
}
