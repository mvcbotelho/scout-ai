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
    │   └── playerHandler_test.go # Testes dos handlers
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

A aplicação estará disponível em `http://localhost:8080`
O banco PostgreSQL estará disponível na porta `5432`

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
O arquivo `docker-compose.yml` configura dois serviços:

1. **go-backend**: Aplicação Go
   - Build do contexto `./go-backend`
   - Porta `8080:8080`
   - Variáveis de ambiente para conexão com banco
   - Dependência do serviço `db`

2. **db**: Banco PostgreSQL
   - Imagem `postgres:15`
   - Porta `5432:5432`
   - Volume persistente para dados
   - Variáveis de ambiente configuradas

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
- **`handlers/playerHandler_test.go`**: Testes automatizados dos handlers
- **`models/player.go`**: Modelo de dados do jogador usando GORM

### Melhorias Implementadas

1. **Validação de Dados**: Todos os endpoints validam dados de entrada
2. **Tratamento de Erros**: Mensagens de erro em português e códigos HTTP apropriados
3. **Verificação de Existência**: Endpoints verificam se recursos existem antes de operações
4. **Variáveis de Ambiente**: Configuração flexível via variáveis de ambiente
5. **Testes Automatizados**: Cobertura de testes para handlers
6. **Makefile**: Comandos úteis para desenvolvimento e deploy

### Adicionando Novas Funcionalidades

1. Crie novos modelos no diretório `models/`
2. Implemente handlers no diretório `handlers/`
3. Adicione testes em `handlers/*_test.go`
4. Adicione novas rotas no arquivo `main.go`
5. Execute `db.AutoMigrate()` para o novo modelo
6. Atualize as dependências se necessário com `go mod tidy`

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