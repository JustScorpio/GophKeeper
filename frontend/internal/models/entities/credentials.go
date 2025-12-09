package entities

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// Credentials - учётные данные
type Credentials struct {
	SecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EncryptFields - шифрует только метаданные (логин и пароль не шифруются)
func (c *Credentials) EncryptFields(cryptoService *crypto.CryptoService) error {
	if c.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(c.Metadata)
		if err != nil {
			return err
		}
		c.Metadata = encryptedMetadata
	}

	return nil
}

// DecryptFields - дешифрует только метаданные
func (c *Credentials) DecryptFields(cryptoService *crypto.CryptoService) error {
	if c.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(c.Metadata)
		if err != nil {
			return err
		}
		c.Metadata = decryptedMetadata
	}

	return nil
}
