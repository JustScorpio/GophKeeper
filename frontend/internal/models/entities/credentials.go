// entities содержит модели сущностей которые хранятся в БД
package entities

import (
	"crypto/sha256"
	"encoding/hex"

	crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
)

// Credentials - учётные данные
type Credentials struct {
	SecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EncryptFields - шифрует все поля
func (entity *Credentials) EncryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = encryptedMetadata
	}

	if entity.Login != "" {
		encryptedData, err := cryptoService.Encrypt(entity.Login)
		if err != nil {
			return err
		}
		entity.Login = encryptedData
	}

	if entity.Password != "" {
		encryptedData, err := cryptoService.Encrypt(entity.Password)
		if err != nil {
			return err
		}
		entity.Password = encryptedData
	}

	return nil
}

// DecryptFields - дешифрует все поля
func (entity *Credentials) DecryptFields(cryptoService *crypto.CryptoService) error {
	if entity.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(entity.Metadata)
		if err != nil {
			return err
		}
		entity.Metadata = decryptedMetadata
	}

	if entity.Login != "" {
		decryptedMetadata, err := cryptoService.Decrypt(entity.Login)
		if err != nil {
			return err
		}
		entity.Login = decryptedMetadata
	}

	if entity.Password != "" {
		decryptedMetadata, err := cryptoService.Decrypt(entity.Password)
		if err != nil {
			return err
		}
		entity.Password = decryptedMetadata
	}

	return nil
}

func (entity *Credentials) GetHash() string {
	// Нулевой байт ([]byte{0}) как разделитель не встретится в данных
	hasher := sha256.New()
	hasher.Write([]byte(entity.ID))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Metadata))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Login))
	hasher.Write([]byte{0})
	hasher.Write([]byte(entity.Password))

	return hex.EncodeToString(hasher.Sum(nil))
}
