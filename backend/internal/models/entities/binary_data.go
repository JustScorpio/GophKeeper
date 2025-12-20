// entities содержит модели сущностей которые хранятся в БД
package entities

// BinaryData - бинарные данные
type BinaryData struct {
	SecureEntity
	Data []byte `json:"data"`
}
