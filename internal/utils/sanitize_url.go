package utils

import (
	"fmt"

	"github.com/PuerkitoBio/purell"
)

// SanitizeURL - проверяет и чистит URL
// CWE-89 - SQL Injection, CWE-918 - Server Side Request Forgery (SSRF)
func SanitizeURL(rawURL string) (string, error) {
	sanitizedURL, err := purell.NormalizeURLString(rawURL, purell.FlagsSafe)
	if err != nil {
		return "", fmt.Errorf("error sanitizing URL: %w", err)
	}
	return sanitizedURL, nil
}
