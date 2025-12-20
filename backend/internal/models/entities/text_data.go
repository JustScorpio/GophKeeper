// entities содержит модели сущностей которые хранятся в БД
package entities

// TextData - текстовые данные
type TextData struct {
	SecureEntity
	Data string `json:"data"`
}
