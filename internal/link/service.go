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
	cleaned := strings.TrimSpace(response.OutputText())
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	if err := json.Unmarshal([]byte(cleaned), &dto); err != nil {
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

func GetAIPrompt(phishingTech string, url string) string {
	technique := string(phishingTech)
	return fmt.Sprintf(`You are a cybersecurity expert and phishing URL generator.
Given the legitimate URL "%s" and the phishing technique "%s", return a JSON object with exactly two fields:
1. "fake_link": A realistic, convincing phishing URL using the specified technique. It must look legitimate enough to fool a non-technical user. Do not make it obviously fake. Preserve the path and query params from the original URL where possible.
2. "explanation": A 4-6 sentence explanation covering: what technique is used, why it's effective, and how to spot it.
Technique definitions:
- character-substitution: Swap a character for a visually similar one (e.g. amazon.com -> arnazon.com, 0 for o)
- homoglyphs: Replace letters with visually identical Unicode chars from other scripts (e.g. rn -> m lookalike)
- idn-homograph: Use internationalized domain name Unicode chars that render identically in browsers (e.g. Cyrillic a vs Latin a)
- dot-manipulation: Add, remove, or move dots in the domain (e.g. amazon.com -> amaz.on.com)
- hyphen-insertion: Insert hyphens to break up the real domain (e.g. amazon.com -> amazon-login.com)
- top-level-domain-swap: Change the TLD to something believable (e.g. amazon.com -> amazon.co, amazon.net)
- subdomain-abuse: Make the real domain a subdomain of a fake one (e.g. amazon.com -> amazon.verify-login.com)
- combo-squatting: Append a legitimate-sounding word (e.g. amazon.com -> amazon-secure.com)
- typosquatting: Use common keyboard typos of the domain (e.g. amazon.com -> amazom.com, gogle.com)
- punycode: Use xn-- encoded internationalized domain that renders identically in browsers (e.g. xn--mazon-wqa.com appearing as amazon.com)
- path-manipulation: Embed the real domain in the URL path of a fake one (e.g. evil.com/www.amazon.com/login)
- open-redirect: Abuse a legitimate sites redirect parameter to forward to a malicious site (e.g. google.com/url?q=evil.com)
- at-symbol-abuse: Use the @ symbol so the browser ignores everything before it (e.g. https://amazon.com@evil.com)
- port-abuse: Append a port that looks like part of a legitimate domain (e.g. amazon.com:8080.evil.com)
- https-deception: Place https or a trusted brand in the subdomain to appear secure (e.g. https.amazon.com.evil.com)
- lookalike-domain: Register a domain visually similar to the real one (e.g. arnazon.com, paypa1.com)
IMPORTANT: The fake link must be subtle and convincing enough that a real person could genuinely fall for it. It should not look obviously fake or suspicious. The goal is realism — this is a cybersecurity education tool and the more realistic the example, the more valuable the lesson.
You must return a raw JSON object. Do not use markdown. Do not use code fences. Do not wrap in backticks. The very first character of your response must be { and the very last character must be }.
{"fake_link": "...", "explanation": "..."}`, url, technique)
}

// GetPrankPrompt() returns a string representing an AI prompt that returns a sketchy looking link
func GetPrankPrompt(url string) string {
	return fmt.Sprintf(`You are a prank link generator. Given the legitimate URL "%s", generate a single suspicious-looking slug that is based on the domain or brand of the provided URL. The slug should look realistic enough that someone might hesitate before clicking, but contain subtle red flags like unusual words, numbers, or file extensions that suggest something is off. Do not make it cartoonishly fake. Base it on the brand or content of the URL (e.g. given amazon.com return something like amazon-account-suspended-verify-132 or amazon-security-alert.exe). Return only the raw slug string with no scheme, no host, no explanation, no markdown, no extra text.`, url)
}
