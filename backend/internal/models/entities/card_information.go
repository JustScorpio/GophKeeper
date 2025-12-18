// entities содержит модели сущностей которые хранятся в БД
package entities

// CardInformation - данные банковской карты
type CardInformation struct {
	SecureEntity
	Number         string `json:"number"`
	CardHolder     string `json:"card_holder"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
