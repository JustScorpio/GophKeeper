// Пакет Main
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JustScorpio/GophKeeper/backend/internal/handlers"
	"github.com/JustScorpio/GophKeeper/backend/internal/middleware/auth"
	"github.com/JustScorpio/GophKeeper/backend/internal/middleware/gzipencoder"
	"github.com/JustScorpio/GophKeeper/backend/internal/middleware/logger"
	"github.com/JustScorpio/GophKeeper/backend/internal/repositories/postgres"
	"github.com/JustScorpio/GophKeeper/backend/internal/services"

	_ "net/http/pprof"

	"github.com/go-chi/chi"
)

var (
	// build-переменные заполняемые с помощью ldflags -X
	buildVersion = "1.0"
	buildDate    = time.Now().Format("January 2, 2006")
)

// main - вызывается автоматически при запуске приложения
func main() {
	// вывести аргументы
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)

	// обрабатываем аргументы командной строки
	parseFlags()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run - функция полезна при инициализации зависимостей сервера перед запуском
func run() error {
	//Проверяем указан ли конфигурационный файл.
	if envConfigPath, hasEnv := os.LookupEnv("CONFIG_PATH"); hasEnv {
		configPath = envConfigPath
	}

	//Заполняем параметры из конфига
	if configPath != "" {
		err := parseConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Строка подключения к базе данных postgres
	if envDBAddr, hasEnv := os.LookupEnv("DATABASE_URI"); hasEnv {
		databaseConnStr = envDBAddr
	}

	// Получение ключа для генерации токенов
	if envSecretKey, hasEnv := os.LookupEnv("AUTH_SECRET_KEY"); hasEnv {
		secretKey = envSecretKey
		auth.Init(secretKey)
	}

	fmt.Println("DB connection string: ", databaseConnStr)

	//Инициализация репозиториев
	dbManager, err := postgres.NewDatabaseManager(databaseConnStr)
	if err != nil {
		return err
	}
	defer dbManager.DB.Close(context.Background())

	// Инициализация сервисов
	storageService := services.NewStorageService(dbManager.UsersRepo, dbManager.BinariesRepo, dbManager.CardsRepo, dbManager.CredentialsRepo, dbManager.TextsRepo)

	//При наличии переменной окружения или флага - запускаем на HTTPS
	_, hasEnv := os.LookupEnv("ENABLE_HTTPS")
	enableHTTPS = hasEnv || enableHTTPS

	//Сертификат для HTTPS
	var tlsConfig *tls.Config
	if enableHTTPS {
		// Чтение пути до сертификата и его ключа
		if envTlsCertPath, hasEnv := os.LookupEnv("TLS_CERT_PATH"); hasEnv {
			tlsCertPath = envTlsCertPath
		}
		if envKeyCertPath, hasEnv := os.LookupEnv("TLS_KEY_PATH"); hasEnv {
			tlsKeyPath = envKeyCertPath
		}

		tlsConfig, err = GetTLSConfigFromFiles(tlsCertPath, tlsKeyPath)
		if err != nil {
			return err
		}
	}

	// Инициализация обработчиков
	handler := handlers.NewGophkeeperHandler(storageService)

	// Инициализация логгера
	zapLogger, err := logger.NewLogger("Info", true)
	if err != nil {
		return err
	}
	defer zapLogger.Sync()

	// Берём адрес сервера из переменной окружения. Иначе - из аргумента
	if envServerAddr, hasEnv := os.LookupEnv("SERVER_ADDRESS"); hasEnv {
		routerAddr = envServerAddr
	}

	// Канал для получения сигналов ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Запуск сервера
	r := chi.NewRouter()

	//Базовые middleware
	r.Use(logger.LoggingMiddleware(zapLogger))
	r.Use(gzipencoder.GZIPEncodingMiddleware())

	//Публичные маршруты
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handler.Register)
		r.Post("/api/user/login", handler.Login)
	})

	//Защищённые маршруты с auth middleware
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware())
		r.Post("/api/user/binaries", handler.CreateBinary)
		r.Get("/api/user/binaries/{id}", handler.GetBinary)
		r.Get("/api/user/binaries", handler.GetAllBinaries)
		r.Put("/api/user/binaries", handler.UpdateBinary)
		r.Delete("/api/user/binaries/{id}", handler.DeleteBinary)

		r.Post("/api/user/card", handler.CreateCard)
		r.Get("/api/user/cards/{id}", handler.GetCard)
		r.Get("/api/user/cards", handler.GetAllCards)
		r.Put("/api/user/cards", handler.UpdateCard)
		r.Delete("/api/user/cards/{id}", handler.DeleteCard)

		r.Post("/api/user/credentials", handler.CreateCredentials)
		r.Get("/api/user/credentials/{id}", handler.GetCredentials)
		r.Get("/api/user/credentials", handler.GetAllCredentials)
		r.Put("/api/user/credentials", handler.UpdateCredentials)
		r.Delete("/api/user/credentials/{id}", handler.DeleteCredentials)

		r.Post("/api/user/texts", handler.CreateText)
		r.Get("/api/user/texts/{id}", handler.GetText)
		r.Get("/api/user/texts", handler.GetAllTexts)
		r.Put("/api/user/texts", handler.UpdateText)
		r.Delete("/api/user/texts/{id}", handler.DeleteText)
	})

	server := createHTTPServer(routerAddr, r, tlsConfig)
	fmt.Println("Running server on", routerAddr)

	// Запуск сервера в горутине
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- runHTTPServer(server)
	}()

	// Ожидание сигнала остановки или ошибки сервера
	select {
	case <-stop:
		fmt.Println("Received shutdown signal")
	case err := <-serverErr:
		fmt.Printf("Server error: %v\n", err)
	}

	return gracefulShutdown(storageService, server)
}

// createServer - создает и настраивает HTTP сервер
func createHTTPServer(addr string, handler http.Handler, tlsConfig *tls.Config) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	if tlsConfig != nil {
		server.TLSConfig = tlsConfig
	}

	return server
}

// runHTTPServer - запускает сервер в горутине и возвращает канал с ошибкой
func runHTTPServer(server *http.Server) error {
	if server.TLSConfig != nil {
		return server.ListenAndServeTLS("", "")
	}
	return server.ListenAndServe()
}

// gracefulShutdown - graceful shutdown приложения
func gracefulShutdown(service *services.StorageService, server *http.Server) error {
	fmt.Println("Starting graceful shutdown...")

	// Останавливаем прием новых соединений
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем сервис (прекращаем обработку задач)
	if service != nil {
		service.Shutdown()
		fmt.Println("Service shutdown completed")
	}

	// Останавливаем HTTP сервер
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
		return err
	}

	fmt.Println("HTTP server shutdown completed")
	return nil
}
