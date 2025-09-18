package domain

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CountryResponse struct {
	Query        string `json:"query"`
	OfficialName string `json:"officialName"`
	ISOCode      string `json:"isoCode"`
}

func NewCountryResponse(query, officialName, isoCode string) *CountryResponse {
	return &CountryResponse{
		Query:        query,
		OfficialName: officialName,
		ISOCode:      isoCode,
	}
}
