package entities

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// Credentials - учётные данные
type Credentials struct {
	SecureEntity
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EncryptFields - шифрует все поля
func (c *Credentials) EncryptFields(cryptoService *crypto.CryptoService) error {
	if c.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(c.Metadata)
		if err != nil {
			return err
		}
		c.Metadata = encryptedMetadata
	}

	if c.Login != "" {
		encryptedData, err := cryptoService.Encrypt(c.Login)
		if err != nil {
			return err
		}
		c.Login = encryptedData
	}

	if c.Password != "" {
		encryptedData, err := cryptoService.Encrypt(c.Password)
		if err != nil {
			return err
		}
		c.Password = encryptedData
	}

	return nil
}

// DecryptFields - дешифрует все поля
func (c *Credentials) DecryptFields(cryptoService *crypto.CryptoService) error {
	if c.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(c.Metadata)
		if err != nil {
			return err
		}
		c.Metadata = decryptedMetadata
	}

	if c.Login != "" {
		decryptedMetadata, err := cryptoService.Decrypt(c.Login)
		if err != nil {
			return err
		}
		c.Login = decryptedMetadata
	}

	if c.Password != "" {
		decryptedMetadata, err := cryptoService.Decrypt(c.Password)
		if err != nil {
			return err
		}
		c.Password = decryptedMetadata
	}

	return nil
}
