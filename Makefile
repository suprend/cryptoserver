.PHONY: build run test clean help

BINARY := cryptoserver

help: ## Показать справку
	@echo "Доступные команды:"
	@echo "  build      - Скомпилировать бинарник"
	@echo "  run        - Запустить сервер (go run)"
	@echo "  test       - Запустить go test ./..."
	@echo "  clean      - Удалить артефакты сборки"
	@echo "  help       - Показать эту справку"

build: ## Скомпилировать бинарник
	@echo "🛠️  Сборка Go-проекта..."
	go build -o $(BINARY) ./cryptoserver.go
	@echo "✅ Сборка завершена"

run: ## Запустить сервер (go run)
	@echo "🚀 Запуск сервера..."
	go run ./cryptoserver.go

test: ## Запустить тесты
	@echo "🧪 Запуск автотестов..."
	go test ./...

clean: ## Очистить артефакты сборки
	@echo "🧹 Очистка..."
	rm -f $(BINARY)
	@echo "✅ Очистка завершена"
