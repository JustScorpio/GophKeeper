package entities

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// CardInformation - данные банковской карты
type CardInformation struct {
	SecureEntity
	Number         string `json:"number"`
	CardHolder     string `json:"card_holder"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}

// EncryptFields - шифрует все поля кроме ID
func (c *CardInformation) EncryptFields(cryptoService *crypto.CryptoService) error {
	if c.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(c.Metadata)
		if err != nil {
			return err
		}
		c.Metadata = encryptedMetadata
	}

	if c.Number != "" {
		encryptedNumber, err := cryptoService.Encrypt(c.Number)
		if err != nil {
			return err
		}
		c.Number = encryptedNumber
	}

	if c.CardHolder != "" {
		encryptedCardHolder, err := cryptoService.Encrypt(c.CardHolder)
		if err != nil {
			return err
		}
		c.CardHolder = encryptedCardHolder
	}

	if c.ExpirationDate != "" {
		encryptedExpirationDate, err := cryptoService.Encrypt(c.ExpirationDate)
		if err != nil {
			return err
		}
		c.ExpirationDate = encryptedExpirationDate
	}

	if c.CVV != "" {
		encryptedCVV, err := cryptoService.Encrypt(c.CVV)
		if err != nil {
			return err
		}
		c.CVV = encryptedCVV
	}

	return nil
}

// DecryptFields - дешифрует все поля кроме ID
func (c *CardInformation) DecryptFields(cryptoService *crypto.CryptoService) error {
	if c.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(c.Metadata)
		if err != nil {
			return err
		}
		c.Metadata = decryptedMetadata
	}

	if c.Number != "" {
		decryptedNumber, err := cryptoService.Decrypt(c.Number)
		if err != nil {
			return err
		}
		c.Number = decryptedNumber
	}

	if c.CardHolder != "" {
		decryptedCardHolder, err := cryptoService.Decrypt(c.CardHolder)
		if err != nil {
			return err
		}
		c.CardHolder = decryptedCardHolder
	}

	if c.ExpirationDate != "" {
		decryptedExpirationDate, err := cryptoService.Decrypt(c.ExpirationDate)
		if err != nil {
			return err
		}
		c.ExpirationDate = decryptedExpirationDate
	}

	if c.CVV != "" {
		decryptedCVV, err := cryptoService.Decrypt(c.CVV)
		if err != nil {
			return err
		}
		c.CVV = decryptedCVV
	}

	return nil
}
