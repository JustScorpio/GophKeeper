package entities

import crypto "github.com/JustScorpio/GophKeeper/frontend/internal/encryption"

// BinaryData - бинарные данные
type BinaryData struct {
	SecureEntity
	Data []byte `json:"data"` //Зашифрованное тоже будет в виде массива байт
}

// EncryptFields - шифрует все поля кроме ID
func (b *BinaryData) EncryptFields(cryptoService *crypto.CryptoService) error {
	if b.Metadata != "" {
		encryptedMetadata, err := cryptoService.Encrypt(b.Metadata)
		if err != nil {
			return err
		}
		b.Metadata = encryptedMetadata
	}

	if len(b.Data) > 0 {
		encryptedData, err := cryptoService.EncryptBytes(b.Data)
		if err != nil {
			return err
		}
		b.Data = encryptedData
	}

	return nil
}

// DecryptFields - дешифрует все поля кроме ID
func (b *BinaryData) DecryptFields(cryptoService *crypto.CryptoService) error {
	if b.Metadata != "" {
		decryptedMetadata, err := cryptoService.Decrypt(b.Metadata)
		if err != nil {
			return err
		}
		b.Metadata = decryptedMetadata
	}

	if len(b.Data) > 0 {
		decryptedData, err := cryptoService.DecryptBytes(b.Data)
		if err != nil {
			return err
		}
		b.Data = decryptedData
	}

	return nil
}
