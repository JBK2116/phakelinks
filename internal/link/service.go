package link

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/JBK2116/phakelinks/types"
)

// ValidateCreateLinkDTO() ensures that the provided CreateLinkDTO holds valid information in all fields
func ValidateCreateLinkDTO(dto types.CreateLinkDTO) error {
	if dto.URL == "" {
		return fmt.Errorf("url is required")
	}
	if dto.Mode == "" {
		return fmt.Errorf("mode is required")
	}
	if !ValidateOriginalURL(dto.URL) {
		return fmt.Errorf("Invalid URL or domain: %s", dto.URL)
	}
	if !ValidateMode(dto.Mode) {
		return fmt.Errorf("Invalid mode: %s", dto.Mode)
	}
	return nil
}

// ValidateOriginalURL() ensures that the provided CreateLinkDTO holds valid url information
func ValidateOriginalURL(originalURL string) bool {
	if strings.HasPrefix(originalURL, "http://") || strings.HasPrefix(originalURL, "https://") {
		return ValidateURL(originalURL)
	} else {
		return ValidateDomain(originalURL)
	}
}

// ValidateURL() checks that the provided URL string is a valid, well-formed HTTP/HTTPS URL.
func ValidateURL(rawURL string) bool {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return false
	}
	return true
}

// ValidateDomain() checks the the provided domain string is a valid, well-formed web domain
func ValidateDomain(rawDomain string) bool {
	regex := `^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	ok, _ := regexp.MatchString(regex, rawDomain)
	return ok
}

func ValidateMode(mode string) bool {
	return mode == string(types.Educational) || mode == string(types.Prank)
}
