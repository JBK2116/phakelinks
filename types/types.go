package types

// CreateLink represents the incoming request payload to generate a phishing URL.
type CreateLinkDTO struct {
	OriginalURL string `json:"original_url"`
	Mode        string `json:"mode"`
}

// ReturnLink represents the response payload containing the original and generated phishing URL.
type ReturnLinkDTO struct {
	OriginalURL string `json:"original_url"`
	FakeURL     string `json:"fake_url"`
	Technique   string `json:"technique"`
	Mode        string `json:"mode"`
}

// Explanation represents the AI-generated explanation linked to a specific URL mapping.
type ExplanationDTO struct {
	ID           int64  `json:"id"`
	URLMappingID int64  `json:"url_mapping_id"`
	Explanation  string `json:"explanation"`
}
