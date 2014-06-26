package balanced

type Address struct {
	City        string `json:"city"`
	Line2       string `json:"line2"`
	Line1       string `json:"line1"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}
