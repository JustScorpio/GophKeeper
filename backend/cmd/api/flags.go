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

	// configPath - путь до конфигурационного файла
	configPath string
)

// parseFlags - обрабатывает аргументы командной строки и сохраняет их значения в соответствующих переменных
func parseFlags() {
	flag.StringVar(&routerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&databaseConnStr, "d", "host=127.0.0.1 user=postgres password=1Qwerty dbname=gophkeeperdb port=5432 sslmode=disable", "postgresql connection string (only for postgresql)")
	flag.BoolVar(&enableHTTPS, "s", false, "enable https")
	flag.StringVar(&configPath, "c", "../configs/app_config.json", "path to application config file")
	flag.Parse()
}
