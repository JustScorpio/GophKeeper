package dtos

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// NewCredentials - учётные данные (dto - новая запись)
type NewCredentials struct {
	NewSecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EncryptFields - шифрует только метаданные (логин и пароль не шифруются)
func (d *NewCredentials) EncryptFields(cryptoService *crypto.CryptoService) error {
	if d.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(d.Metadata)
		if err != nil {
			return err
		}
		d.Metadata = encryptedMetadata
	}

	return nil
}
