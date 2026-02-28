package link

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/JBK2116/phakelinks/internal/configs"
	"github.com/JBK2116/phakelinks/types"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// ValidateCreateLinkDTO() ensures that the provided CreateLinkDTO holds valid information in all fields
func ValidateCreateLinkDTO(dto types.CreateLinkDTO) *types.ErrorResponse {
	if dto.Link == "" {
		return &types.ErrorResponse{Error: "MISSING_URL", Message: "A URL is required to create a link."}
	}
	if dto.Mode == "" {
		return &types.ErrorResponse{
			Error:   "MISSING_MODE",
			Message: "A mode must be selected.",
		}
	}
	if dto.Exclude == nil {
		return &types.ErrorResponse{
			Error:   "MISSING_EXCLUDE",
			Message: "An exclude list is required. Pass an empty array if you have no exclusions.",
		}
	}
	if err := ValidateLink(dto.Link); err != nil {
		return &types.ErrorResponse{
			Error:   "INVALID_URL",
			Message: "The URL or domain is not valid. Ensure it includes a scheme (e.g. https://) and a proper domain.",
			Value:   err.Error(),
		}
	}
	if !ValidateMode(dto.Mode) {
		return &types.ErrorResponse{
			Error:   "INVALID_MODE",
			Message: "The provided mode is not valid",
			Value:   dto.Mode,
		}
	}
	if err := ValidateExcludes(dto.Exclude); err != nil {
		return &types.ErrorResponse{
			Error:   "INVALID_EXCLUDE",
			Message: "One or more exclude patterns are invalid. Ensure all patterns are valid URL paths or glob expressions.",
			Extra:   err.Error(),
		}
	}
	return nil
}

// ValidateLink() ensures that the provided CreateLinkDTO holds valid url information
func ValidateLink(link string) error {
	httpClient := &http.Client{}
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return ValidateURL(link, *httpClient)
	} else {
		return ValidateDomain(link, *httpClient)
	}
}

// ValidateURL() checks that the provided URL string is a valid, well-formed HTTP/HTTPS URL.
func ValidateURL(rawURL string, client http.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	url := "https://api.cloudmersive.com/validate/domain/url/full"
	method := "POST"
	body, err := json.Marshal(map[string]string{"URL": rawURL})
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(body))
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ApiKey", configs.Envs.CLOUDMERSIVE_KEY)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	fmt.Println(string(resBody))
	if err != nil {
		return err
	}
	result := make(map[string]any)
	if err := json.Unmarshal(resBody, &result); err != nil {
		return err
	}
	if valid, ok := result["ValidURL"].(bool); !ok || !valid {
		return fmt.Errorf("%s", rawURL)
	}
	return nil
}

// ValidateDomain() checks the the provided domain string is a valid, well-formed web domain
func ValidateDomain(rawDomain string, client http.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	url := "https://api.cloudmersive.com/validate/domain/check"
	method := "POST"
	payload := strings.NewReader(rawDomain)
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ApiKey", configs.Envs.CLOUDMERSIVE_KEY)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	result := make(map[string]bool)
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if !result["ValidDomain"] {
		return fmt.Errorf("%s", rawDomain)
	}
	return nil
}

// ValidateMode() checks that the provided mode is a valid mode defined in `types.go`
func ValidateMode(mode string) bool {
	return mode == string(types.Educational) || mode == string(types.Prank)
}

func ValidateExcludes(exclude []string) error {
	if len(exclude) >= len(types.AllPhishingTechniques) {
		return fmt.Errorf("Length of exclude array must be less than %d", len(types.AllPhishingTechniques))
	}
	seen := make(map[string]struct{})
	for _, v := range exclude {
		if _, exists := seen[v]; exists {
			continue
		}
		if !contains(types.AllPhishingTechniques, v) {
			return fmt.Errorf("%s is an invalid exclude type", v)
		}
		seen[v] = struct{}{}
	}
	return nil
}

// contains() checks if the provided string is in the provided array slice
func contains[T ~string](slice []T, s string) bool {
	for _, v := range slice {
		if string(v) == s {
			return true
		}
	}
	return false
}

// GetRandomPhishingTechnique() returns a random PhishingTechnique enum
func GetRandomPhishingTechnique(excludes []string) string {
	availableTechniques := make([]string, 0)
	for _, v := range types.AllPhishingTechniques {
		if !contains(excludes, string(v)) {
			availableTechniques = append(availableTechniques, string(v))
		}
	}
	randIndex := rand.IntN(len(availableTechniques))
	return availableTechniques[randIndex]
}

// GetEducationalAISummary() queries the OPENAI API for the AI summary, returning the `ExplanationDTO` if successful
func GetEducationalAISummary(phishingTech string, url string) (types.ExplanationDTO, error) {
	duration := time.Minute * 1
	ctx, cancelCtx := context.WithTimeout(context.Background(), duration)
	defer cancelCtx()

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
	dto.Technique = phishingTech
	return dto, nil
}

func GetPrankLink(url string) (types.PrankDTO, error) {
	duration := time.Minute * 1
	ctx, cancelCtx := context.WithTimeout(context.Background(), duration)
	defer cancelCtx()
	client := openai.NewClient(option.WithAPIKey(configs.Envs.OPENAI_KEY))
	question := GetPrankPrompt(url)
	response, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModelGPT4o,
	})
	var dto types.PrankDTO
	if err != nil {
		return dto, err
	}
	if configs.Envs.IsDev {
		dto.Link = fmt.Sprintf("%s:%s/%s", configs.Envs.RedirectHost, configs.Envs.RedirectPort, response.OutputText())
	} else {
		dto.Link = fmt.Sprintf("%s/%s", configs.Envs.RedirectHost, response.OutputText())
	}
	dto.Slug = strings.TrimSpace(response.OutputText())
	return dto, nil
}

// GetAIPrompt() returns a string representing an AI prompt that returns an educational summary
func GetAIPrompt(phishingTech string, url string) string {
	technique := string(phishingTech)

	return fmt.Sprintf(`You are a phishing URL generator and cybersecurity educator.

Given the legitimate URL "%s" and the phishing technique "%s", return a JSON object with exactly two fields:
1. "fake_link": A realistic phishing URL using the specified technique
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
{"fake_link": "...", "explanation": "..."}`, url, technique)
}

// GetPrankPrompt() returns a string representing an AI prompt that returns a sketchy looking link
func GetPrankPrompt(url string) string {
	return fmt.Sprintf(`You are a prank link generator. Given the legitimate URL "%s", generate a single fake-looking suspicious link string that appears related to the domain/content of the URL but looks obviously sketchy (e.g. if given amazon.com, return something like amazon-free-gift-exe.zip or amazon_login_verify-132.exe.zip). Return only the raw link string, no JSON, no explanation, no markdown.`, url)
}
