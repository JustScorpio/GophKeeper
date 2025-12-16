// Пакет Main
package main

import (
	"crypto/tls"
	"errors"
	"os"
)

// GetTLSConfigFromFiles - загрузить TLS-конфигурацию из файлов сертификата и ключа
func GetTLSConfigFromFiles(certPath, keyPath string) (tlsConfig *tls.Config, err error) {
	// Проверяем файлы
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return nil, errors.New("файл сертификата не найден: " + certPath)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, errors.New("файл ключа не найден: " + keyPath)
	}

	// Загружаем сертификат и ключ из файлов
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
