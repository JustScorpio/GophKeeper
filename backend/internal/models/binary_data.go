package models

// BinaryData - files to store
type BinaryData struct {
	ID       string `json:"id"`
	Metadata string `json:"metadata"`
	Data     []byte `json:"data"`
}
