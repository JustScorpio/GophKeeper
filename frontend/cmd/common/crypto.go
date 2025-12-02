package main

// // На клиенте
// package client

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/base64"
// 	"fmt"
// 	"io"
// )

// // DeriveKey - создать ключ из пароля
// func DeriveKey(password string) []byte {
// 	hash := sha256.Sum256([]byte(password))
// 	return hash[:]
// }

// // EncryptData - зашифровать данные
// func EncryptData(data []byte, password string) (string, error) {
// 	key := DeriveKey(password)

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return "", err
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return "", err
// 	}

// 	ciphertext := gcm.Seal(nonce, nonce, data, nil)
// 	return base64.StdEncoding.EncodeToString(ciphertext), nil
// }

// // DecryptData - расшифровывать данные
// func DecryptData(encryptedData string, password string) ([]byte, error) {
// 	key := DeriveKey(password)

// 	data, err := base64.StdEncoding.DecodeString(encryptedData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nonceSize := gcm.NonceSize()
// 	if len(data) < nonceSize {
// 		return nil, fmt.Errorf("ciphertext too short")
// 	}

// 	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
// 	return gcm.Open(nil, nonce, ciphertext, nil)
// }
