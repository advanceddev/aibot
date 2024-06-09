package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumStrings(t *testing.T) {
	t.Run("Конкатенация нескольких строк", func(t *testing.T) {
		result := SumStrings("Hello, ", "World!")
		expected := "Hello, World!"
		assert.Equal(t, expected, result, "Несоответствие ожидаемым результатам.")
	})

	t.Run("Конкатенация пустых строк", func(t *testing.T) {
		result := SumStrings("", "")
		expected := ""
		assert.Equal(t, expected, result, "Несоответствие ожидаемым результатам: должна быть пустая строка.")
	})

	t.Run("Конкатенация одной строки", func(t *testing.T) {
		result := SumStrings("Just one string")
		expected := "Just one string"
		assert.Equal(t, expected, result, "Несоответствие ожидаемым результатам: должа быть та же строка без изменений.")
	})

	t.Run("Конкатенация с использованием специальных символов", func(t *testing.T) {
		result := SumStrings("Special ", "chars: ", "∆ß©️∑¨ˆ¥Øƒ∆")
		expected := "Special chars: ∆ß©️∑¨ˆ¥Øƒ∆"
		assert.Equal(t, expected, result, "Несоответствие ожидаемым результатам.")
	})
}

func TestSanitizeURL(t *testing.T) {
	t.Run("Sanitizing valid URL", func(t *testing.T) {
		rawURL := "https://example.com/path?q=123"
		sanitizedURL, err := SanitizeURL(rawURL)

		assert.NoError(t, err, "При обработке валидного URL не должно быть ошибки.")
		assert.Equal(t, rawURL, sanitizedURL, "Обработанный URL должен соответствовать исходному.")
	})

	t.Run("Sanitizing malformed URL", func(t *testing.T) {
		rawURL := "hts://example.com"
		targetURL := "https://example.com"
		sanitizedURL, err := SanitizeURL(rawURL)

		assert.NoError(t, err, "При обработке URL с ошибкой протокола ошибка должна устраняться.")
		assert.Equal(t, targetURL, sanitizedURL, "Ошибка в протоколе должна быть исправлена.")
	})

	t.Run("Sanitizing URL without protocol", func(t *testing.T) {
		rawURL := "example.com"
		targetURL := "https://example.com"
		sanitizedURL, err := SanitizeURL(rawURL)

		assert.NoError(t, err, "Обработка URL без протокола не должна вызывать ошибку.")
		assert.Equal(t, targetURL, sanitizedURL, "Обработанный URL должен иметь протокол https для URL без протокола.")
	})
}
