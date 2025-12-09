package dtos

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// NewTextData - текстовые данные (dto - новая запись)
type NewTextData struct {
	NewSecureEntity
	Data string `json:"data"`
}

// EncryptFields - шифрует все поля
func (d *NewTextData) EncryptFields(cryptoService *crypto.CryptoService) error {
	if d.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(d.Metadata)
		if err != nil {
			return err
		}
		d.Metadata = encryptedMetadata
	}

	if d.Data != "" {
		encryptedData, err := cryptoService.Encrypt(d.Data)
		if err != nil {
			return err
		}
		d.Data = encryptedData
	}

	return nil
}
