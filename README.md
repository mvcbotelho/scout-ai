# Scout AI

Um projeto backend em Go que utiliza o framework Gin para criar uma API REST simples e eficiente.

## 📋 Descrição

O Scout AI é uma aplicação backend desenvolvida em Go que fornece uma API REST básica. O projeto está estruturado para ser facilmente containerizado e deployado usando Docker.

## 🚀 Tecnologias Utilizadas

- **Go 1.22.2** - Linguagem de programação principal
- **Gin Framework** - Framework web para Go
- **Docker** - Containerização da aplicação
- **Docker Compose** - Orquestração de containers

## 📁 Estrutura do Projeto

```
scout-ai/
├── docker-compose.yml          # Configuração do Docker Compose
├── Makefile                    # Comandos de automação (vazio)
├── README.md                   # Este arquivo
└── go-backend/                 # Código fonte do backend
    ├── cmd/
    │   └── main.go            # Ponto de entrada da aplicação
    ├── Dockerfile             # Configuração do container Docker
    ├── go.mod                 # Dependências do Go
    └── go.sum                 # Checksums das dependências
```

## 🛠️ Pré-requisitos

- [Go 1.22.2](https://golang.org/dl/) ou superior
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## ⚙️ Instalação e Execução

### Opção 1: Execução Local

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd scout-ai
```

2. Navegue para o diretório do backend:
```bash
cd go-backend
```

3. Baixe as dependências:
```bash
go mod download
```

4. Execute a aplicação:
```bash
go run cmd/main.go
```

A aplicação estará disponível em `http://localhost:8080`

### Opção 2: Execução com Docker

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

### Opção 3: Execução Manual com Docker

1. Navegue para o diretório do backend:
```bash
cd go-backend
```

2. Construa a imagem Docker:
```bash
docker build -t scout-ai-backend .
```

3. Execute o container:
```bash
docker run -p 8080:8080 scout-ai-backend
```

## 🔧 Endpoints da API

### Health Check
- **GET** `/ping`
  - **Descrição**: Endpoint de verificação de saúde da aplicação
  - **Resposta**: `{"message": "pong"}`
  - **Status**: 200 OK

**Exemplo de uso:**
```bash
curl http://localhost:8080/ping
```

## 🐳 Configuração Docker

### Dockerfile
O projeto utiliza um Dockerfile multi-stage para otimizar o tamanho da imagem final:

1. **Etapa de Build**: Compila a aplicação usando `golang:1.22-alpine`
2. **Etapa de Execução**: Cria uma imagem final baseada em `alpine:latest`

### Docker Compose
O arquivo `docker-compose.yml` configura:
- Build do contexto `./go-backend`
- Mapeamento da porta `8080:8080`
- Comentário sobre remoção de volume para evitar sobrescrita

## 📦 Dependências Principais

- `github.com/gin-gonic/gin v1.10.1` - Framework web Gin
- Dependências de suporte para JSON, validação, e outras funcionalidades

## 🔍 Desenvolvimento

### Estrutura do Código

O arquivo principal `cmd/main.go` contém:
- Configuração do servidor Gin
- Definição das rotas da API
- Inicialização do servidor na porta 8080

### Adicionando Novas Funcionalidades

1. Crie novos handlers no diretório `cmd/` ou crie um diretório `handlers/`
2. Adicione novas rotas no arquivo `main.go`
3. Atualize as dependências se necessário com `go mod tidy`

## 🚀 Deploy

### Produção
Para deploy em produção, recomenda-se:

1. Usar variáveis de ambiente para configurações
2. Implementar logging estruturado
3. Adicionar métricas e monitoramento
4. Configurar health checks mais robustos
5. Implementar rate limiting e segurança

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
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/ping"]
      interval: 30s
      timeout: 10s
      retries: 3
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

**Marcos Botelho** - [mvcbotelho](https://github.com/mvcbotelho)

## 📞 Suporte

Para suporte, abra uma issue no repositório do projeto ou entre em contato através do email.

---

**Nota**: Este é um projeto em desenvolvimento. Novas funcionalidades e melhorias estão sendo implementadas continuamente. 