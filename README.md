# Scout AI

Um sistema de scouting de jogadores desenvolvido em Go que utiliza o framework Gin e PostgreSQL para gerenciar dados de atletas de futebol.

## 📋 Descrição

O Scout AI é uma aplicação backend desenvolvida em Go que fornece uma API REST completa para gerenciamento de dados de jogadores de futebol. O sistema permite criar, consultar, atualizar e deletar informações de atletas, incluindo estatísticas como gols, tackles e passes. O projeto está estruturado para ser facilmente containerizado e deployado usando Docker.

## 🚀 Tecnologias Utilizadas

- **Go 1.22** - Linguagem de programação principal
- **Gin Framework** - Framework web para Go
- **GORM** - ORM para Go
- **PostgreSQL** - Banco de dados relacional
- **Docker** - Containerização da aplicação
- **Docker Compose** - Orquestração de containers
- **Testify** - Framework de testes
- **Ollama3** - Modelo de linguagem local para análise inteligente

## 📁 Estrutura do Projeto

```
scout-ai/
├── docker-compose.yml          # Configuração do Docker Compose
├── Makefile                    # Comandos de automação
├── README.md                   # Este arquivo
├── .gitignore                  # Arquivos ignorados pelo Git
└── go-backend/                 # Código fonte do backend
    ├── cmd/
    │   └── main.go            # Ponto de entrada da aplicação
    ├── handlers/
    │   ├── playerHandler.go   # Handlers para endpoints de jogadores
    │   ├── playerHandler_test.go # Testes dos handlers de jogadores
    │   ├── analyzeHandler.go  # Handlers para análise com IA
    │   ├── analyzeHandler_test.go # Testes dos handlers de análise
    │   └── ollamaHandler.go   # Integração com Ollama3
    ├── models/
    │   └── player.go          # Modelo de dados do jogador
    ├── Dockerfile             # Configuração do container Docker
    ├── go.mod                 # Dependências do Go
    └── go.sum                 # Checksums das dependências
```

## 🛠️ Pré-requisitos

- [Go 1.22](https://golang.org/dl/) ou superior
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## ⚙️ Instalação e Execução

### Opção 1: Execução com Docker (Recomendado)

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd scout-ai
```

2. Execute com Docker Compose:
```bash
docker-compose up --build
```

3. **Configure o Ollama (após os serviços iniciarem):**
```bash
# Aguarde alguns segundos e execute:
make setup-ollama
```

**OU use o comando completo:**
```bash
make deploy-with-ollama
```

A aplicação estará disponível em `http://localhost:8080`
O banco PostgreSQL estará disponível na porta `5432`
O Ollama estará disponível na porta `11434`

### Opção 2: Execução Local

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd scout-ai
```

2. Configure o banco PostgreSQL localmente ou use Docker:
```bash
docker run --name postgres-scout -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=scoutdb -p 5432:5432 -d postgres:15
```

3. Navegue para o diretório do backend:
```bash
cd go-backend
```

4. Baixe as dependências:
```bash
go mod download
```

5. Execute a aplicação:
```bash
go run cmd/main.go
```

### Opção 3: Usando Makefile

O projeto inclui um Makefile com comandos úteis:

```bash
# Ver todos os comandos disponíveis
make help

# Executar em modo desenvolvimento
make dev

# Construir e executar com Docker
make docker-build
make docker-up

# Executar testes
make test

# Verificar status dos serviços
make check
```

## 🔧 Endpoints da API

### Health Check
- **GET** `/ping`
  - **Descrição**: Endpoint de verificação de saúde da aplicação
  - **Resposta**: `{"message": "pong"}`
  - **Status**: 200 OK

### Jogadores (Players)

#### Criar Jogador
- **POST** `/players`
  - **Descrição**: Cria um novo jogador
  - **Validações**: Nome obrigatório, idade > 0
  - **Body**: 
    ```json
    {
      "name": "João Silva",
      "age": 25,
      "position": "Atacante",
      "team": "Flamengo",
      "goals": 15,
      "tackles": 5,
      "passes": 120
    }
    ```
  - **Status**: 201 Created
  - **Erro**: 400 Bad Request (dados inválidos)

#### Listar Jogadores
- **GET** `/players`
  - **Descrição**: Retorna todos os jogadores cadastrados
  - **Status**: 200 OK

#### Buscar Jogador por ID
- **GET** `/players/:id`
  - **Descrição**: Retorna um jogador específico pelo ID
  - **Validações**: ID deve ser um número válido
  - **Status**: 200 OK
  - **Erro**: 404 Not Found (jogador não encontrado)

#### Atualizar Jogador
- **PUT** `/players/:id`
  - **Descrição**: Atualiza dados de um jogador
  - **Validações**: ID válido, nome obrigatório, idade > 0
  - **Body**: Mesmo formato do POST
  - **Status**: 200 OK
  - **Erro**: 404 Not Found (jogador não encontrado)

#### Deletar Jogador
- **DELETE** `/players/:id`
  - **Descrição**: Remove um jogador do sistema
  - **Validações**: ID deve ser um número válido
  - **Status**: 200 OK
  - **Erro**: 404 Not Found (jogador não encontrado)

### Análise de Jogadores (AI-Powered Analysis)

#### Analisar Jogador Específico
- **GET** `/analyze/players/:id`
  - **Descrição**: Gera análise completa de um jogador usando IA
  - **Parâmetros**: 
    - `ai=true` - Usar Ollama3 para análise (opcional)
  - **Validações**: ID deve ser um número válido
  - **Resposta**:
    ```json
    {
      "player_id": 1,
      "player_name": "João Silva",
      "analysis": "João Silva demonstra características de um atacante moderno e versátil. Com 25 anos, está na idade ideal para desenvolvimento e consolidação de carreira. Sua média de 0.5 gols por jogo é impressionante, demonstrando excelente eficiência de finalização. A participação de 120 passes na temporada mostra que não é apenas um finalizador, mas também contribui para a construção do jogo. Sua eficiência de 123.0 pontos o classifica como um jogador de nível Bom, adequado para competições de nível médio a alto. Recomendo atenção especial ao seu potencial de evolução, considerando sua idade e versatilidade técnica.",
      "insights": [
        "Excelente eficiência de finalização com 0.5 gols por jogo",
        "Participação ativa na construção do jogo com 120 passes",
        "Idade ideal para desenvolvimento e consolidação",
        "Versatilidade técnica como atacante moderno"
      ],
      "rating": 8,
      "position": "Atacante",
      "team": "Flamengo",
      "ai_used": true
    }
    ```
  - **Status**: 200 OK
  - **Erro**: 404 Not Found (jogador não encontrado)

#### Analisar Todos os Jogadores
- **GET** `/analyze/players`
  - **Descrição**: Gera análise individual e comparativa de todos os jogadores
  - **Parâmetros**: 
    - `ai=true` - Usar Ollama3 para análise (opcional)
  - **Resposta**:
    ```json
    {
      "individual_analyses": [
        {
          "player_id": 1,
          "player_name": "João Silva",
          "analysis": "...",
          "insights": [...],
          "rating": 8,
          "position": "Atacante",
          "team": "Flamengo",
          "ai_used": true
        }
      ],
      "comparative_analysis": {
        "total_players": 4,
        "total_goals": 25,
        "total_tackles": 185,
        "total_passes": 735,
        "positions_distribution": {
          "Atacante": 1,
          "Meio-campo": 1,
          "Zagueiro": 1,
          "Goleiro": 1
        },
        "best_scorer": {
          "name": "João Silva",
          "goals": 15,
          "team": "Flamengo"
        },
        "best_tackler": {
          "name": "Carlos Oliveira",
          "tackles": 120,
          "team": "São Paulo"
        },
        "best_passer": {
          "name": "Pedro Santos",
          "passes": 350,
          "team": "Palmeiras"
        },
        "ai_comparative_analysis": "Análise comparativa gerada pelo Ollama3..."
      }
    }
    ```
  - **Status**: 200 OK

#### Comparar Jogadores
- **GET** `/analyze/compare?ids=1&ids=2&ids=3`
  - **Descrição**: Compara dois ou mais jogadores especificados
  - **Parâmetros**: 
    - `ids` - IDs dos jogadores (mínimo 2)
    - `ai=true` - Usar Ollama3 para análise (opcional)
  - **Validações**: Pelo menos 2 IDs válidos
  - **Resposta**:
    ```json
    {
      "players": [
        {
          "player_id": 1,
          "player_name": "João Silva",
          "analysis": "...",
          "insights": [...],
          "rating": 8,
          "position": "Atacante",
          "team": "Flamengo",
          "ai_used": true
        }
      ],
      "comparison": {
        "highest_rating": {...},
        "most_goals": {...},
        "most_tackles": {...},
        "most_passes": {...}
      },
      "ai_comparative_analysis": "Análise comparativa detalhada gerada pelo Ollama3..."
    }
    ```
  - **Status**: 200 OK
  - **Erro**: 400 Bad Request (IDs insuficientes), 404 Not Found (jogadores não encontrados)

## 📝 Exemplos de Uso

### Exemplos para Postman/Insomnia

#### Exemplo 1 - Atacante
```json
{
    "name": "João Silva",
    "age": 25,
    "position": "Atacante",
    "team": "Flamengo",
    "goals": 15,
    "tackles": 5,
    "passes": 120
}
```

#### Exemplo 2 - Meio-campista
```json
{
    "name": "Pedro Santos",
    "age": 28,
    "position": "Meio-campo",
    "team": "Palmeiras",
    "goals": 8,
    "tackles": 45,
    "passes": 350
}
```

#### Exemplo 3 - Zagueiro
```json
{
    "name": "Carlos Oliveira",
    "age": 32,
    "position": "Zagueiro",
    "team": "São Paulo",
    "goals": 2,
    "tackles": 120,
    "passes": 180
}
```

#### Exemplo 4 - Goleiro
```json
{
    "name": "Rafael Costa",
    "age": 29,
    "position": "Goleiro",
    "team": "Corinthians",
    "goals": 0,
    "tackles": 15,
    "passes": 85
}
```

### Exemplos de Análise com IA

#### Analisar Jogador Específico
**Método:** GET  
**URL:** `http://localhost:8080/analyze/players/1?ai=true`  
**Resposta esperada:**
```json
{
    "player_id": 1,
    "player_name": "João Silva",
    "analysis": "João Silva demonstra características de um atacante moderno e versátil. Com 25 anos, está na idade ideal para desenvolvimento e consolidação de carreira. Sua média de 0.5 gols por jogo é impressionante, demonstrando excelente eficiência de finalização. A participação de 120 passes na temporada mostra que não é apenas um finalizador, mas também contribui para a construção do jogo. Sua eficiência de 123.0 pontos o classifica como um jogador de nível Bom, adequado para competições de nível médio a alto. Recomendo atenção especial ao seu potencial de evolução, considerando sua idade e versatilidade técnica.",
    "insights": [
        "Excelente eficiência de finalização com 0.5 gols por jogo",
        "Participação ativa na construção do jogo com 120 passes",
        "Idade ideal para desenvolvimento e consolidação",
        "Versatilidade técnica como atacante moderno"
    ],
    "rating": 8,
    "position": "Atacante",
    "team": "Flamengo",
    "ai_used": true
}
```

#### Analisar Todos os Jogadores
**Método:** GET  
**URL:** `http://localhost:8080/analyze/players`  
**Resposta:** Análise individual de cada jogador + análise comparativa geral

#### Comparar Jogadores
**Método:** GET  
**URL:** `http://localhost:8080/analyze/compare?ids=1&ids=2&ids=3`  
**Resposta:** Comparação detalhada entre os jogadores especificados

### Exemplos com cURL
```bash
# Health check
curl http://localhost:8080/ping

# Criar jogador
curl -X POST http://localhost:8080/players \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","age":25,"position":"Atacante","team":"Flamengo","goals":15,"tackles":5,"passes":120}'

# Listar jogadores
curl http://localhost:8080/players

# Buscar jogador por ID
curl http://localhost:8080/players/1

# Atualizar jogador
curl -X PUT http://localhost:8080/players/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","age":26,"position":"Atacante","team":"Flamengo","goals":18,"tackles":5,"passes":125}'

# Deletar jogador
curl -X DELETE http://localhost:8080/players/1

# Analisar jogador específico
curl http://localhost:8080/analyze/players/1

# Analisar jogador com IA (Ollama3)
curl "http://localhost:8080/analyze/players/1?ai=true"

# Analisar todos os jogadores
curl http://localhost:8080/analyze/players

# Analisar todos os jogadores com IA
curl "http://localhost:8080/analyze/players?ai=true"

# Comparar jogadores
curl "http://localhost:8080/analyze/compare?ids=1&ids=2&ids=3"

# Comparar jogadores com IA
curl "http://localhost:8080/analyze/compare?ids=1&ids=2&ids=3&ai=true"
```

### Exemplos com PowerShell
```powershell
# Criar jogador
$body = @{
    name     = "João Silva"
    age      = 25
    position = "Atacante"
    team     = "Flamengo"
    goals    = 15
    tackles  = 5
    passes   = 120
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/players" -Method Post -Body $body -ContentType "application/json"

# Analisar jogador
Invoke-RestMethod -Uri "http://localhost:8080/analyze/players/1" -Method Get

# Analisar jogador com IA
Invoke-RestMethod -Uri "http://localhost:8080/analyze/players/1?ai=true" -Method Get

# Analisar todos os jogadores
Invoke-RestMethod -Uri "http://localhost:8080/analyze/players" -Method Get

# Analisar todos os jogadores com IA
Invoke-RestMethod -Uri "http://localhost:8080/analyze/players?ai=true" -Method Get

# Comparar jogadores
Invoke-RestMethod -Uri "http://localhost:8080/analyze/compare?ids=1&ids=2" -Method Get

# Comparar jogadores com IA
Invoke-RestMethod -Uri "http://localhost:8080/analyze/compare?ids=1&ids=2&ai=true" -Method Get
```

## 🗄️ Modelo de Dados

### Player (Jogador)
```go
type Player struct {
    gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
    Name     string   `json:"name" binding:"required" gorm:"not null"`     // Nome do jogador (obrigatório)
    Age      int      `json:"age" binding:"required,min=1,max=100" gorm:"not null"`      // Idade (1-100)
    Position string   `json:"position" binding:"required" gorm:"not null"` // Posição (obrigatório)
    Team     string   `json:"team" binding:"required" gorm:"not null"`     // Time atual (obrigatório)
    Goals    int      `json:"goals" binding:"min=0" gorm:"default:0"`    // Número de gols (>= 0)
    Tackles  int      `json:"tackles" binding:"min=0" gorm:"default:0"`  // Número de tackles (>= 0)
    Passes   int      `json:"passes" binding:"min=0" gorm:"default:0"`   // Número de passes (>= 0)
}
```

## 🧪 Testes

O projeto inclui testes automatizados para os handlers:

```bash
# Executar todos os testes
go test ./...

# Executar testes com verbose
go test ./handlers -v

# Executar testes com cobertura
go test ./handlers -cover
```

**Nota**: Os testes usam SQLite em memória e podem requerer dependências C no Windows.

## 🐳 Configuração Docker

### Docker Compose
O arquivo `docker-compose.yml` configura três serviços:

1. **go-backend**: Aplicação Go
   - Build do contexto `./go-backend`
   - Porta `8080:8080`
   - Variáveis de ambiente para conexão com banco e Ollama
   - Dependência dos serviços `db` e `ollama`

2. **db**: Banco PostgreSQL
   - Imagem `postgres:15`
   - Porta `5432:5432`
   - Volume persistente para dados
   - Variáveis de ambiente configuradas

3. **ollama**: Serviço Ollama3
   - Imagem `ollama/ollama:latest`
   - Porta `11434:11434`
   - Volume persistente para modelos
   - Download automático do modelo `llama3.2`

### Dockerfile
O projeto utiliza um Dockerfile multi-stage para otimizar o tamanho da imagem final:

1. **Etapa de Build**: Compila a aplicação usando `golang:1.22-alpine`
2. **Etapa de Execução**: Cria uma imagem final baseada em `alpine:latest`

## 📦 Dependências Principais

- `github.com/gin-gonic/gin v1.10.1` - Framework web Gin
- `gorm.io/gorm v1.25.9` - ORM para Go
- `gorm.io/driver/postgres v1.4.6` - Driver PostgreSQL para GORM
- `github.com/stretchr/testify v1.9.0` - Framework de testes
- Dependências de suporte para JSON, validação, e outras funcionalidades

## 🔍 Desenvolvimento

### Estrutura do Código

- **`cmd/main.go`**: Ponto de entrada da aplicação, configuração do servidor e rotas
- **`handlers/playerHandler.go`**: Handlers HTTP para operações CRUD de jogadores
- **`handlers/playerHandler_test.go`**: Testes automatizados dos handlers de jogadores
- **`handlers/analyzeHandler.go`**: Handlers HTTP para análise de jogadores com IA
- **`handlers/analyzeHandler_test.go`**: Testes automatizados dos handlers de análise
- **`handlers/ollamaHandler.go`**: Integração com Ollama3
- **`models/player.go`**: Modelo de dados do jogador usando GORM

### Melhorias Implementadas

1. **Validação de Dados**: Todos os endpoints validam dados de entrada
2. **Tratamento de Erros**: Mensagens de erro em português e códigos HTTP apropriados
3. **Verificação de Existência**: Endpoints verificam se recursos existem antes de operações
4. **Variáveis de Ambiente**: Configuração flexível via variáveis de ambiente
5. **Testes Automatizados**: Cobertura de testes para handlers
6. **Makefile**: Comandos úteis para desenvolvimento e deploy
7. **Análise com IA**: Sistema de análise automatizada de jogadores
8. **Integração com Ollama3**: Análises mais naturais e contextuais
9. **Prompts Especializados**: Prompts específicos para cada posição
10. **Fallback Inteligente**: Sistema de fallback para análise estática
11. **Comparação de Jogadores**: Funcionalidade para comparar performance entre jogadores
12. **Insights Inteligentes**: Geração automática de insights baseados em dados
13. **Rating System**: Sistema de pontuação de 1-10 para jogadores

## 🚀 Deploy

### Produção
Para deploy em produção, recomenda-se:

1. Usar variáveis de ambiente para configurações sensíveis
2. Implementar logging estruturado
3. Adicionar métricas e monitoramento
4. Configurar health checks mais robustos
5. Implementar rate limiting e segurança
6. Configurar backup do banco PostgreSQL
7. Usar secrets management para senhas

### Exemplo de Deploy com Docker Compose em Produção
```yaml
version: '3.8'
services:
  go-backend:
    build:
      context: ./go-backend
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
    restart: unless-stopped
    depends_on:
      - db
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  db:
    image: postgres:15
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  postgres_data:
```

## 🤝 Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença [MIT](LICENSE).

## 👨‍💻 Autor

**Marcus Botelho** - [GitHub](https://github.com/mvcbotelho) | [LinkedIn](https://www.linkedin.com/in/mvcbotelho/)

## 📞 Suporte

Para suporte, abra uma issue no repositório do projeto ou entre em contato através do LinkedIn.

---

**Nota**: Este é um projeto em desenvolvimento. Novas funcionalidades e melhorias estão sendo implementadas continuamente.

## 🤖 Funcionalidades de IA

### Integração com Ollama3

O Scout AI agora integra com o **Ollama3**, um modelo de linguagem local que permite análises mais sofisticadas e naturais dos jogadores. Esta integração oferece:

#### **Vantagens da Integração com Ollama3:**
- **Análises mais naturais**: Linguagem mais fluida e profissional
- **Insights contextuais**: Análises baseadas em conhecimento de futebol
- **Flexibilidade**: Funciona offline e localmente
- **Personalização**: Prompts específicos para scouting
- **Fallback inteligente**: Se Ollama não estiver disponível, usa análise estática

#### **Como Usar a IA:**
Adicione o parâmetro `?ai=true` aos endpoints de análise:

```bash
# Análise com IA
GET /analyze/players/1?ai=true

# Análise comparativa com IA
GET /analyze/players?ai=true

# Comparação de jogadores com IA
GET /analyze/compare?ids=1&ids=2&ai=true
```

#### **Configuração do Ollama3:**
O sistema está configurado para usar o modelo `llama3.2` por padrão, mas pode ser personalizado via variáveis de ambiente:

```yaml
environment:
  - OLLAMA_BASE_URL=http://ollama:11434
  - OLLAMA_MODEL=llama3.2
  - OLLAMA_TEMPERATURE=0.7
  - OLLAMA_TOP_P=0.9
```

#### **Prompts Especializados:**
O sistema gera prompts específicos para cada posição:

- **Atacantes**: Foco em finalização, participação no jogo e eficiência
- **Meio-campistas**: Análise de visão de jogo, distribuição e marcação
- **Zagueiros**: Avaliação defensiva e saída de bola
- **Goleiros**: Comando de área e saída do gol

### Análise Inteligente de Jogadores

O Scout AI utiliza algoritmos inteligentes para analisar dados de jogadores e gerar insights valiosos para scouting:

#### **Sistema de Rating (1-10)**
- Calcula pontuação baseada em eficiência, idade e performance
- Considera diferentes pesos para cada posição
- Ajusta rating baseado na experiência do jogador

#### **Geração de Insights**
- **Baseados na Idade**: Identifica jogadores jovens com potencial ou experientes para liderança
- **Baseados na Posição**: Análises específicas para atacantes, meio-campistas, zagueiros e goleiros
- **Baseados na Performance**: Identifica jogadores com performance excepcional ou que precisam melhorar

#### **Análise Comparativa**
- Compara jogadores em diferentes categorias
- Identifica melhores artilheiros, marcadores e passadores
- Distribuição de posições no elenco
- Estatísticas gerais do time

#### **Recomendações Automáticas**
- Sugere adequação para diferentes níveis de competição
- Identifica jogadores prontos para times de alto nível
- Recomenda tempo de desenvolvimento para jogadores jovens

### Algoritmos Utilizados

#### **Cálculo de Eficiência por Posição**
```go
// Atacante: Foco em gols e participação no jogo
Efficiency = (Goals * 0.6) + (Passes * 0.4)

// Meio-campista: Equilíbrio entre passes, marcação e gols
Efficiency = (Passes * 0.5) + (Tackles * 0.3) + (Goals * 0.2)

// Zagueiro: Foco em marcação e saída de bola
Efficiency = (Tackles * 0.6) + (Passes * 0.4)

// Goleiro: Foco em saída do gol
Efficiency = (Tackles * 0.8) + (Passes * 0.2)
```

#### **Sistema de Rating**
- **Base**: 5 pontos
- **Eficiência > 200**: +3 pontos
- **Eficiência > 150**: +2 pontos
- **Eficiência > 100**: +1 ponto
- **Eficiência < 50**: -1 ponto
- **Idade ideal (25-30)**: +1 ponto
- **Muito jovem (<20)**: -1 ponto

### Adicionando Novas Funcionalidades

1. Crie novos modelos no diretório `models/`
2. Implemente handlers no diretório `handlers/`
3. Adicione testes em `handlers/*_test.go`
4. Adicione novas rotas no arquivo `main.go`
5. Execute `db.AutoMigrate()` para o novo modelo
6. Atualize as dependências se necessário com `go mod tidy` 

## 🔧 Troubleshooting

### Problemas Comuns

#### **Erro no Ollama: "unknown command"**
Se você ver erros como `Error: unknown command "sh" for "ollama"`:

1. **Pare os containers:**
```bash
docker-compose down
```

2. **Reinicie com configuração correta:**
```bash
make deploy-with-ollama
```

#### **Ollama não está respondendo**
Se os endpoints com `?ai=true` retornam análise estática:

1. **Verifique se o Ollama está rodando:**
```bash
make ollama-status
```

2. **Configure o Ollama manualmente:**
```bash
make setup-ollama
```

3. **Verifique os logs:**
```bash
docker-compose logs ollama
```

#### **Modelo não encontrado**
Se o Ollama não consegue baixar o modelo:

1. **Verifique a conexão:**
```bash
curl http://localhost:11434/api/tags
```

2. **Baixe o modelo manualmente:**
```bash
curl -X POST http://localhost:11434/api/pull \
  -H "Content-Type: application/json" \
  -d '{"name": "llama3.2"}'
```

#### **Fallback para Análise Estática**
Se o Ollama não estiver disponível, o sistema automaticamente usa análise estática. Você pode verificar isso no campo `ai_used: false` na resposta.

### Comandos Úteis

```bash
# Verificar status de todos os serviços
make check

# Verificar status do Ollama
make ollama-status

# Configurar Ollama
make setup-ollama

# Logs do Ollama
docker-compose logs ollama

# Reiniciar apenas o Ollama
docker-compose restart ollama
```

#### **Erro de Migração: "insufficient arguments"**
Se você ver erros como `Erro ao fazer migração: insufficient arguments`:

1. **Pare os containers:**
```bash
docker-compose down
```

2. **Limpe os volumes (cuidado - isso apaga os dados):**
```bash
docker-compose down -v
```

3. **Reinicie com configuração correta:**
```bash
make deploy-with-ollama
```

**Causa:** O banco de dados não estava totalmente pronto quando a aplicação tentou fazer a migração.

**Solução:** Adicionamos healthcheck e retry automático para resolver este problema. 