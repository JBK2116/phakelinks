package link

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/JBK2116/phakelinks/internal/configs"
	"github.com/JBK2116/phakelinks/types"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
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

// ValidateMode() checks that the provided mode is a valid mode defined in `types.go`
func ValidateMode(mode string) bool {
	return mode == string(types.Educational) || mode == string(types.Prank)
}

// GetRandomPhishingTechnique() returns a random PhishingTechnique enum
func GetRandomPhishingTechnique() types.PhishingTechnique {
	amountOptions := len(types.AllPhishingTechniques)
	randIndex := rand.IntN(amountOptions)
	return types.AllPhishingTechniques[randIndex]
}

func GetEducationalAISummary(phishingTech types.PhishingTechnique, url string) (types.ExplanationDTO, error) {
	duration := time.Minute * 1
	ctx, cancelContext := context.WithTimeout(context.Background(), duration)
	defer cancelContext()

	client := openai.NewClient(
		option.WithAPIKey(configs.Envs.OPENAI_KEY),
	)
	question := GetAIPrompt(phishingTech, url)
	response, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModelGPT4o,
	})
	var dto types.ExplanationDTO
	if err != nil {
		return dto, err
	}
	if err := json.Unmarshal([]byte(response.OutputText()), &dto); err != nil {
		return dto, err
	}
	return dto, nil
}

// GetAIPrompt() returns a string representing an AI prompt
func GetAIPrompt(phishingTech types.PhishingTechnique, url string) string {
	technique := string(phishingTech)

	return fmt.Sprintf(`You are a phishing URL generator and cybersecurity educator.

Given the legitimate URL "%s" and the phishing technique "%s", return a JSON object with exactly two fields:
1. "fake_url": A realistic phishing URL using the specified technique
2. "explanation": A 3-4 sentence explanation covering: what technique is used, why it's effective, and how to spot it

Technique definitions:
- character-substitution: Swap a character for a similar-looking one (e.g. amazon.com → arnazon.com, 0 for o)
- homoglyphs: Replace letters with visually identical Unicode chars from other scripts (e.g. rn → m lookalike)
- idn-homograph: Use internationalized domain name Unicode chars that render identically (e.g. Cyrillic а vs Latin a)
- dot-manipulation: Add, remove, or move dots in the domain (e.g. amazon.com → amaz.on.com)
- hyphen-insertion: Insert hyphens to break up the real domain (e.g. amazon.com → a-mazon.com or amazon-login.com)
- top-level-domain-swap: Change the TLD (e.g. amazon.com → amazon.co or amazon.net or amazon.org)
- subdomain-abuse: Make the real domain a subdomain of a fake one (e.g. amazon.com → amazon.verify-login.com)
- combo-squatting: Append a legitimate-sounding word to the real domain (e.g. amazon.com → amazon-secure.com or amazonlogin.com)
Do not wrap the response in markdown code fences or backticks. Return raw JSON only
Respond with ONLY valid JSON, no markdown, no extra text:
{"fake_url": "...", "explanation": "..."}`, url, technique)
}
