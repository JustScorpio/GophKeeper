// Пакет Main
package main

import (
	"flag"
)

var (
	// routerAddr - адрес и порт для запуска сервера
	routerAddr string

	// databaseConnStr - строка подключения к БД (postgres)
	databaseConnStr string

	// enableHTTPS - включение HTTPS
	enableHTTPS bool

	// secretKey - секретный ключ используемый при выдаче токенов авторизации
	secretKey string

	// tlsCertPath - путь до tls-сертификата
	tlsCertPath string

	// tlsKeyPath - путь до ключа tls-сертификата
	tlsKeyPath string
)

// parseFlags - обрабатывает аргументы командной строки и сохраняет их значения в соответствующих переменных
func parseFlags() {
	flag.StringVar(&routerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&databaseConnStr, "d", "host=127.0.0.1 user=postgres password=1Qwerty dbname=gophkeeperdb port=5432 sslmode=disable", "postgresql connection string (only for postgresql)")
	flag.BoolVar(&enableHTTPS, "s", true, "enable https")
	flag.StringVar(&secretKey, "k", "supersecretkey", "secret key for token creation")
	flag.StringVar(&tlsCertPath, "cp", "../tls/localhost+2.pem", "path to tls certificate")
	flag.StringVar(&tlsKeyPath, "kp", "../tls/localhost+2-key.pem", "path to tls certificate key")
	flag.Parse()
}
