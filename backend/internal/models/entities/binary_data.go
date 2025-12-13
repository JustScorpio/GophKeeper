package entities

// BinaryData - бинарные данные
type BinaryData struct {
	SecureEntity
	Data []byte `json:"data"`
}
