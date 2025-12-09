package interfaces

import (
	"github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
)

type Encryptable interface {
	EncryptFields(cryptoService *encryption.CryptoService) error
}
type Decryptable interface {
	DecryptFields(cryptoService *encryption.CryptoService) error
}
