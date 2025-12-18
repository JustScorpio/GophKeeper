// entities содержит модели сущностей которые хранятся в БД
package entities

import (
	"crypto/sha256"
	"encoding/hex"

	crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
)

// BinaryData - бинарные данные
type BinaryData struct {
	SecureEntity
	Data []byte `json:"data"`
}

// EncryptFields - шифрует все поля кроме ID
func (entity *BinaryData) EncryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = encryptedMetadata
	}

	if len(entity.Data) > 0 {
		encryptedData, err := cryptoService.EncryptBytes(entity.Data)
		if err != nil {
			return err
		}
		entity.Data = encryptedData
	}

	return nil
}

// DecryptFields - дешифрует все поля кроме ID
func (entity *BinaryData) DecryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = decryptedMetadata
	}

	if len(entity.Data) > 0 {
		decryptedData, err := cryptoService.DecryptBytes(entity.Data)
		if err != nil {
			return err
		}
		entity.Data = decryptedData
	}

	return nil
}

func (entity *BinaryData) GetHash() string {
	// Нулевой байт ([]byte{0}) как разделитель не встретится в данных
	hasher := sha256.New()
	hasher.Write([]byte(entity.ID))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Metadata))
	hasher.Write([]byte{0})
	hasher.Write(entity.Data)

	return hex.EncodeToString(hasher.Sum(nil))
}
