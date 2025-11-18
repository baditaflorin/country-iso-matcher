package domain

// Country represents a country with its ISO codes and multilingual information
type Country struct {
	ISO2    string            `json:"iso2"`
	ISO3    string            `json:"iso3"`
	Names   map[string]string `json:"names"`   // Language code -> Name
	Aliases []string          `json:"aliases"` // All aliases for this country
}

// CountryResponse is the API response structure
type CountryResponse struct {
	Query        string `json:"query"`
	OfficialName string `json:"officialName"`
	ISO2Code     string `json:"iso2Code"`
	ISO3Code     string `json:"iso3Code"`
}

// Legacy support for backward compatibility
func (c *Country) Code() string {
	return c.ISO2
}

// GetOfficialName returns the English name or the first available name
func (c *Country) GetOfficialName() string {
	if name, ok := c.Names["en"]; ok {
		return name
	}
	// Fallback to first available name
	for _, name := range c.Names {
		return name
	}
	return ""
}

func NewCountryResponse(query string, country *Country) *CountryResponse {
	return &CountryResponse{
		Query:        query,
		OfficialName: country.GetOfficialName(),
		ISO2Code:     country.ISO2,
		ISO3Code:     country.ISO3,
	}
}
