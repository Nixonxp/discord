.PHONY: build
build:
	@printf '\033[36mСборка контейнеров\033[0m\n'
	@docker compose build

.PHONY: up
up: ## (пере)сборка образов и запуск приложения
	@printf '\033[36mЗапуск контейнеров\033[0m\n'
	@docker compose up -d

.PHONY: down
down: ## остановка приложения
	@docker compose down

.PHONY: restart
restart: ## перезапуск приложения
	@docker compose down
	@docker compose up -d

.PHONY: rebuild-app ## пересборка и перезапуск приложения
rebuild-app: build restart