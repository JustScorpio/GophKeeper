package dtos

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// NewCredentials - учётные данные (dto - новая запись)
type NewCredentials struct {
	NewSecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EncryptFields - шифрует все поля
func (d *NewCredentials) EncryptFields(cryptoService *crypto.CryptoService) error {
	if d.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(d.Metadata)
		if err != nil {
			return err
		}
		d.Metadata = encryptedMetadata
	}

	if d.Login != "" {
		encryptedData, err := cryptoService.Encrypt(d.Login)
		if err != nil {
			return err
		}
		d.Login = encryptedData
	}

	if d.Password != "" {
		encryptedData, err := cryptoService.Encrypt(d.Password)
		if err != nil {
			return err
		}
		d.Password = encryptedData
	}

	return nil
}
