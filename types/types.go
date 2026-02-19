package types

// Mode represents an enum type of the applications state
type Mode string

// const here stores all Mode enums
const (
	Educational Mode = "educational"
	Prank       Mode = "prank"
)

// PhishingTechnique represents an enum type of a PhishingTechnique string
type PhishingTechnique string

// const here stores all PhishingTechnique enums
const (
	CharacterSub    PhishingTechnique = "character-substitution"
	HomoGlyphs      PhishingTechnique = "homoglyphs"
	IDNHomograph    PhishingTechnique = "idn-homograph"
	DotManipulation PhishingTechnique = "dot-manipulation"
	HyphenInsertion PhishingTechnique = "hyphen-insertion"
	TLDSwap         PhishingTechnique = "top-level-domain-swap"
	SubDomainAbuse  PhishingTechnique = "subdomain-abuse"
	ComboSquatting  PhishingTechnique = "combo-squatting"
)

// CreateLink represents the incoming request payload to generate a phishing URL.
type CreateLinkDTO struct {
	URL  string `json:"url"`
	Mode string `json:"mode"`
}

// ReturnLink represents the response payload containing the original and generated phishing URL.
type ReturnLinkDTO struct {
	URL       string `json:"url"`
	FakeURL   string `json:"fake_url"`
	Technique string `json:"technique"`
	Mode      string `json:"mode"`
}

// Explanation represents the AI-generated explanation linked to a specific URL mapping.
type ExplanationDTO struct {
	ID           int64  `json:"id"`
	URLMappingID int64  `json:"url_mapping_id"`
	Explanation  string `json:"explanation"`
}
