package dtos

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// NewBinaryData - бинарные данные (dto - новая запись)
type NewBinaryData struct {
	NewSecureEntity
	Data []byte `json:"data"` //Зашифрованное тоже будет в виде массива байт
}

// EncryptFields - шифрует поля DTO
func (d *NewBinaryData) EncryptFields(cryptoService *crypto.CryptoService) error {
	if d.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(d.Metadata)
		if err != nil {
			return err
		}
		d.Metadata = encryptedMetadata
	}

	if len(d.Data) > 0 {
		encryptedData, err := cryptoService.EncryptBytes(d.Data)
		if err != nil {
			return err
		}
		d.Data = encryptedData
	}

	return nil
}
