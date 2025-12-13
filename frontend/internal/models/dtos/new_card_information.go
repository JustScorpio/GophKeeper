package dtos

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// NewCardInformation - данные банковской карты (dto - новая запись)
type NewCardInformation struct {
	NewSecureEntity
	Number         string `json:"number"`
	CardHolder     string `json:"card_holder"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}

// EncryptFields - шифрует все поля
func (d *NewCardInformation) EncryptFields(cryptoService *crypto.CryptoService) error {
	if d.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(d.Metadata)
		if err != nil {
			return err
		}
		d.Metadata = encryptedMetadata
	}

	if d.Number != "" {
		encryptedNumber, err := cryptoService.Encrypt(d.Number)
		if err != nil {
			return err
		}
		d.Number = encryptedNumber
	}

	if d.CardHolder != "" {
		encryptedCardHolder, err := cryptoService.Encrypt(d.CardHolder)
		if err != nil {
			return err
		}
		d.CardHolder = encryptedCardHolder
	}

	if d.ExpirationDate != "" {
		encryptedExpirationDate, err := cryptoService.Encrypt(d.ExpirationDate)
		if err != nil {
			return err
		}
		d.ExpirationDate = encryptedExpirationDate
	}

	if d.CVV != "" {
		encryptedCVV, err := cryptoService.Encrypt(d.CVV)
		if err != nil {
			return err
		}
		d.CVV = encryptedCVV
	}

	return nil
}
