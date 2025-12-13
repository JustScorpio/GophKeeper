package dtos

// NewTextData - текстовые данные (dto - новая запись)
type NewTextData struct {
	NewSecureEntity
	Data string `json:"data"`
}
