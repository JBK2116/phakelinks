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

const (
	CharacterSub     PhishingTechnique = "character-substitution"
	HomoGlyphs       PhishingTechnique = "homoglyphs"
	IDNHomograph     PhishingTechnique = "idn-homograph"
	DotManipulation  PhishingTechnique = "dot-manipulation"
	HyphenInsertion  PhishingTechnique = "hyphen-insertion"
	TLDSwap          PhishingTechnique = "top-level-domain-swap"
	SubDomainAbuse   PhishingTechnique = "subdomain-abuse"
	ComboSquatting   PhishingTechnique = "combo-squatting"
	TypoSquatting    PhishingTechnique = "typosquatting"
	Punycode         PhishingTechnique = "punycode"
	PathManipulation PhishingTechnique = "path-manipulation"
	OpenRedirect     PhishingTechnique = "open-redirect"
	AtSymbolAbuse    PhishingTechnique = "at-symbol-abuse"
	PortAbuse        PhishingTechnique = "port-abuse"
	HTTPSDeception   PhishingTechnique = "https-deception"
	LookAlikeDomain  PhishingTechnique = "lookalike-domain"
)

var AllPhishingTechniques = []PhishingTechnique{
	CharacterSub,
	HomoGlyphs,
	IDNHomograph,
	DotManipulation,
	HyphenInsertion,
	TLDSwap,
	SubDomainAbuse,
	ComboSquatting,
	TypoSquatting,
	Punycode,
	PathManipulation,
	OpenRedirect,
	AtSymbolAbuse,
	PortAbuse,
	HTTPSDeception,
	LookAlikeDomain,
}

// CreateLink represents the incoming request payload to generate a phishing URL.
type CreateLinkDTO struct {
	Link    string   `json:"link"`
	Mode    string   `json:"mode"`
	Exclude []string `json:"exclude"`
}

// ReturnLink represents the response payload containing the original and generated phishing URL.
type ReturnLinkDTO struct {
	Link        string `json:"link"`
	FakeLink    string `json:"fake_link"`
	Technique   string `json:"technique,omitempty"`
	Mode        string `json:"mode"`
	Explanation string `json:"explanation,omitempty"`
}

// Explanation represents the AI-generated explanation linked to a specific URL mapping.
type ExplanationDTO struct {
	FakeLink    string `json:"fake_link"`
	Technique   string `json:"technique,omitempty"`
	Explanation string `json:"explanation"`
}

type PrankDTO struct {
	Link string `json:"link"`
	Slug string `json:"slug,omitempty"`
}

// ErrorResponse represents an error that occurs during runtime
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
	Extra   string `json:"extra,omitempty"`
}
