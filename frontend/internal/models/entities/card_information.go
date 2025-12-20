// entities содержит модели сущностей которые хранятся в БД
package entities

import (
	"crypto/sha256"
	"encoding/hex"

	crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
)

// CardInformation - данные банковской карты
type CardInformation struct {
	SecureEntity
	Number         string `json:"number"`
	CardHolder     string `json:"card_holder"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}

// EncryptFields - шифрует все поля кроме ID
func (entity *CardInformation) EncryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = encryptedMetadata
	}

	if entity.Number != "" {
		encryptedNumber, err := cryptoService.Encrypt(entity.Number)
		if err != nil {
			return err
		}
		entity.Number = encryptedNumber
	}

	if entity.CardHolder != "" {
		encryptedCardHolder, err := cryptoService.Encrypt(entity.CardHolder)
		if err != nil {
			return err
		}
		entity.CardHolder = encryptedCardHolder
	}

	if entity.ExpirationDate != "" {
		encryptedExpirationDate, err := cryptoService.Encrypt(entity.ExpirationDate)
		if err != nil {
			return err
		}
		entity.ExpirationDate = encryptedExpirationDate
	}

	if entity.CVV != "" {
		encryptedCVV, err := cryptoService.Encrypt(entity.CVV)
		if err != nil {
			return err
		}
		entity.CVV = encryptedCVV
	}

	return nil
}

// DecryptFields - дешифрует все поля кроме ID
func (entity *CardInformation) DecryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = decryptedMetadata
	}

	if entity.Number != "" {
		decryptedNumber, err := cryptoService.Decrypt(entity.Number)
		if err != nil {
			return err
		}
		entity.Number = decryptedNumber
	}

	if entity.CardHolder != "" {
		decryptedCardHolder, err := cryptoService.Decrypt(entity.CardHolder)
		if err != nil {
			return err
		}
		entity.CardHolder = decryptedCardHolder
	}

	if entity.ExpirationDate != "" {
		decryptedExpirationDate, err := cryptoService.Decrypt(entity.ExpirationDate)
		if err != nil {
			return err
		}
		entity.ExpirationDate = decryptedExpirationDate
	}

	if entity.CVV != "" {
		decryptedCVV, err := cryptoService.Decrypt(entity.CVV)
		if err != nil {
			return err
		}
		entity.CVV = decryptedCVV
	}

	return nil
}

func (entity *CardInformation) GetHash() string {
	// Нулевой байт ([]byte{0}) как разделитель не встретится в данных
	hasher := sha256.New()
	hasher.Write([]byte(entity.ID))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Metadata))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Number))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.CardHolder))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.ExpirationDate))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.CVV))

	return hex.EncodeToString(hasher.Sum(nil))
}
