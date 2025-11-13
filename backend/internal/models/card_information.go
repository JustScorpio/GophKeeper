package models

// CardInformation - bank card information to store
type CardInformation struct {
	ID             string `json:"id"`
	Metadata       string `json:"metadata"`
	Number         string `json:"number"`
	CardHolder     string `json:"card_holder"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
