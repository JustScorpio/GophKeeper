package entities

// TextData - текстовые данные
type TextData struct {
	SecureEntity
	Data string `json:"data"`
}
