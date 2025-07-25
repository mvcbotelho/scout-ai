# Makefile para Scout AI

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

# Comandos principais
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila a aplicação
	cd go-backend && go build -o bin/main ./cmd

run: ## Executa a aplicação localmente
	cd go-backend && go run cmd/main.go

test: ## Executa os testes
	cd go-backend && go test ./...

clean: ## Remove arquivos de build
	rm -rf go-backend/bin
	rm -rf go-backend/tmp

# Comandos Docker
docker-build: ## Constrói as imagens Docker
	docker-compose build

docker-up: ## Inicia os containers
	docker-compose up -d

docker-down: ## Para os containers
	docker-compose down

docker-logs: ## Mostra os logs dos containers
	docker-compose logs -f

docker-restart: ## Reinicia os containers
	docker-compose restart

# Comandos de desenvolvimento
dev: ## Inicia o ambiente de desenvolvimento
	docker-compose up --build

dev-logs: ## Mostra logs em modo desenvolvimento
	docker-compose logs -f go-backend

# Comandos de banco de dados
db-reset: ## Reseta o banco de dados (cuidado!)
	docker-compose down -v
	docker-compose up -d db
	sleep 5
	docker-compose up --build

# Comandos de verificação
check: ## Verifica se tudo está funcionando
	@echo "Verificando conectividade..."
	@curl -s http://localhost:8080/ping || echo "Servidor não está rodando"
	@echo "Verificando banco de dados..."
	@docker-compose ps db | grep -q "Up" || echo "Banco de dados não está rodando"
	@echo "Verificando Ollama..."
	@curl -s http://localhost:11434/api/tags || echo "Ollama não está rodando"

# Comandos Ollama
setup-ollama: ## Configura o Ollama (baixa o modelo)
	@echo "Aguardando Ollama estar disponível..."
	@sleep 10
	@echo "Baixando modelo llama3.2..."
	@curl -X POST http://localhost:11434/api/pull -H "Content-Type: application/json" -d '{"name": "llama3.2"}' || echo "Erro ao baixar modelo"
	@echo "Ollama configurado com sucesso!"

ollama-status: ## Verifica status do Ollama
	@curl -s http://localhost:11434/api/tags || echo "Ollama não está disponível"

# Comandos de deploy
deploy: ## Deploy completo
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "Deploy concluído! Acesse http://localhost:8080"
	@echo "Para configurar Ollama, execute: make setup-ollama"

deploy-with-ollama: ## Deploy completo com configuração do Ollama
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "Aguardando serviços iniciarem..."
	@sleep 15
	@echo "Configurando Ollama..."
	@make setup-ollama
	@echo "Deploy completo! Acesse http://localhost:8080"
