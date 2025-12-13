// Пакет Main
package main

import (
	"encoding/json"
	"os"
)

// parseConfig - обрабатывает параметры запуска приложения из конфигурационного файла
func parseConfig(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var appConfig struct {
		ServerAddress string `json:"server_address"`
		DatabaseDSN   string `json:"database_dsn"`
		EnableHTTPS   bool   `json:"enable_https"`
	}

	err = json.Unmarshal(content, &appConfig)
	if err != nil {
		return err
	}

	routerAddr = appConfig.ServerAddress
	databaseConnStr = appConfig.DatabaseDSN
	enableHTTPS = appConfig.EnableHTTPS

	return nil
}
