package utils

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/purell"
)

// SanitizeURL - проверяет и чистит URL
// CWE-89 - SQL Injection, CWE-918 - Server Side Request Forgery (SSRF)
func SanitizeURL(rawURL string) (string, error) {
	// Проверяем, что входящий URL не пустой
	if rawURL == "" {
		return "", fmt.Errorf("empty URL")
	}

	// Проверяем и исправляем протокол
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	} else {
		schema := strings.ToLower(strings.Split(rawURL, "://")[0])
		if schema != "http" && schema != "https" {
			rawURL = strings.Replace(rawURL, schema+"://", "https://", 1)
		}
	}

	// Проверяем и исправляем URL с использованием библиотеки purell
	sanitizedURL, err := purell.NormalizeURLString(
		rawURL, purell.FlagRemoveDuplicateSlashes|
			purell.FlagRemoveFragment|
			purell.FlagRemoveTrailingSlash|
			purell.FlagsSafe|
			purell.FlagLowercaseScheme,
	)
	if err != nil {
		// Если невалидный URL, возвращаем ошибку
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	return sanitizedURL, nil
}
