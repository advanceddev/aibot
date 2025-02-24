package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - Структура конфигурации приложения
type Config struct {
	BotToken    string `yaml:"BOT_TOKEN" env-required:"true" env:"BOT_TOKEN"`
	BotID       string `yaml:"BOT_ID" env:"BOT_ID"`
	ChannelID   int64  `yaml:"CHANNEL_ID" env-required:"true" env:"CHANNEL_ID"`
	APIToken    string `yaml:"API_TOKEN" env-required:"true" env:"API_TOKEN"`
	APIUrl      string `yaml:"API_URL" env-required:"true" env:"API_URL"`
	AdminUserID int64  `yaml:"ADMIN_USER_ID" env-required:"true" env:"ADMIN_USER_ID"`
}

const (
	configPath = "cfg.yaml"
)

// loadConfigFromFile загружает конфигурацию из файла
func loadConfigFromFile(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("Переменные из файла не получены: %w", err)
	}
	return &cfg, nil
}

// loadConfigFromEnv загружает конфигурацию из переменных окружения
func loadConfigFromEnv() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("Переменные из окружения не получены: %w", err)
	}
	return &cfg, nil
}

// MustLoad загружает конфигурацию из файла или из переменных окружения
func MustLoad() *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Файл конфигурации не найден: %s", configPath)
	}

	cfg, err := loadConfigFromEnv()
	if err != nil {
		fmt.Printf("Не удалось получить переменные из окружения: %v. Пробуем получить из файла %s\n", err, configPath)
		cfg, err = loadConfigFromFile(configPath)
		if err != nil {
			log.Fatalf("Ошибка конфигурации: %v", err)
		}
	}

	return cfg
}
