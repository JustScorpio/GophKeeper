package entities

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// TextData - текстовые данные
type TextData struct {
	SecureEntity
	Data string `json:"data"`
}

// EncryptFields - шифрует все поля кроме ID
func (t *TextData) EncryptFields(cryptoService *crypto.CryptoService) error {
	if t.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(t.Metadata)
		if err != nil {
			return err
		}
		t.Metadata = encryptedMetadata
	}

	if t.Data != "" {
		encryptedData, err := cryptoService.Encrypt(t.Data)
		if err != nil {
			return err
		}
		t.Data = encryptedData
	}

	return nil
}

// DecryptFields - дешифрует все поля кроме ID
func (t *TextData) DecryptFields(cryptoService *crypto.CryptoService) error {
	if t.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(t.Metadata)
		if err != nil {
			return err
		}
		t.Metadata = decryptedMetadata
	}

	if t.Data != "" {
		decryptedData, err := cryptoService.Decrypt(t.Data)
		if err != nil {
			return err
		}
		t.Data = decryptedData
	}

	return nil
}
