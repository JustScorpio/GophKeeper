// dtos содержит объекты для транспортировки данных
package dtos

// NewCardInformation - данные банковской карты (dto - новая запись)
type NewCardInformation struct {
	NewSecureEntity
	Number         string `json:"number"`
	CardHolder     string `json:"card_holder"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
