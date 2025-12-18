// entities содержит модели сущностей которые хранятся в БД
package entities

import (
	"crypto/sha256"
	"encoding/hex"

	crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
)

// TextData - текстовые данные
type TextData struct {
	SecureEntity
	Data string `json:"data"`
}

// EncryptFields - шифрует все поля кроме ID
func (entity *TextData) EncryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = encryptedMetadata
	}

	if entity.Data != "" {
		encryptedData, err := cryptoService.Encrypt(entity.Data)
		if err != nil {
			return err
		}
		entity.Data = encryptedData
	}

	return nil
}

// DecryptFields - дешифрует все поля кроме ID
func (entity *TextData) DecryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = decryptedMetadata
	}

	if entity.Data != "" {
		decryptedData, err := cryptoService.Decrypt(entity.Data)
		if err != nil {
			return err
		}
		entity.Data = decryptedData
	}

	return nil
}

func (entity *TextData) GetHash() string {
	// Нулевой байт ([]byte{0}) как разделитель не встретится в данных
	hasher := sha256.New()
	hasher.Write([]byte(entity.ID))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Metadata))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Data))

	return hex.EncodeToString(hasher.Sum(nil))
}
