# Переменные
BINARY_NAME=main
APP_CMD_PATH=cmd/main.go
PM2_APP_NAME="Unreal AI Bot"
PM2_PROCESS_ID=6

# Запуск приложения в режиме разработки
dev:
	go run $(APP_CMD_PATH)

# Обновление приложения из репозитория, сборка и перезапуск сервера
update: git-pull build restart-app-server

# Скачивание обновлений из репозитория
git-pull:
	git pull

# Сборка приложения
build:
	CGO_ENABLED=0 go build -o $(BINARY_NAME) -ldflags "-w -s" $(APP_CMD_PATH)

# Перезапуск приложения с помощью pm2 и обновление переменных среды
restart-app-server:
	pm2 restart $(PM2_PROCESS_ID) --name $(PM2_APP_NAME) --update-env

.PHONY: dev update git-pull build restart-app-server