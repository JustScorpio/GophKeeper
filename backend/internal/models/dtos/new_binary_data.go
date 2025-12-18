// dtos содержит объекты для транспортировки данных
package dtos

// NewBinaryData - бинарные данные (dto - новая запись)
type NewBinaryData struct {
	NewSecureEntity
	Data []byte `json:"data"`
}
