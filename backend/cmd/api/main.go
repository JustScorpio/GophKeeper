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

	"google.golang.org/grpc"

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

	fmt.Println("DB connection string: ", databaseConnStr)

	// Инициализация репозиториев
	db, err := postgres.NewDBConnection(databaseConnStr)
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	binariesRepo, err := postgres.NewPgBinariesRepo(db)
	if err != nil {
		return err
	}
	cardsRepo, err := postgres.NewPgCardsRepo(db)
	if err != nil {
		return err
	}
	credentialsRepo, err := postgres.NewPgCredentialsRepo(db)
	if err != nil {
		return err
	}
	textsRepo, err := postgres.NewPgTextsRepo(db)
	if err != nil {
		return err
	}
	usersRepo, err := postgres.NewPgUsersRepo(db)
	if err != nil {
		return err
	}

	// Инициализация сервисов
	storageService := services.NewStorageService(usersRepo, binariesRepo, cardsRepo, credentialsRepo, textsRepo)

	//При наличии переменной окружения или наличии флага - запускаем на HTTPS.
	if _, hasEnv := os.LookupEnv("ENABLE_HTTPS"); hasEnv {
		enableHTTPS = true
	}

	//Сертификат для HTTPS
	var tlsConfig *tls.Config
	if enableHTTPS {
		tlsConfig, err = GetTestTLSConfig()
		if err != nil {
			return err
		}
	}

	// Инициализация обработчиков
	shURLHandler := handlers.NewShURLHandler(storageService, enableHTTPS)

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
	r.Use(auth.AuthMiddleware())
	r.Use(logger.LoggingMiddleware(zapLogger))
	r.Use(gzipencoder.GZIPEncodingMiddleware())
	r.Get("/ping", pingFunc)
	r.Get("/api/user/urls", shURLHandler.GetShURLsByUserID)
	r.Delete("/api/user/urls", shURLHandler.DeleteMany)
	r.Get("/{token}", shURLHandler.GetFullURL)
	r.Post("/api/shorten", shURLHandler.ShortenURL)
	r.Post("/api/shorten/batch", shURLHandler.ShortenURLsBatch)
	r.Post("/", shURLHandler.ShortenURL)
	r.With(cidrWhiteList.CIDRWhitelistMiddleware()).Get("/api/internal/stats", shURLHandler.GetStats)

	server := createHTTPServer(routerAddr, r, tlsConfig)

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

// runServer - запускает сервер в горутине и возвращает канал с ошибкой
func runHTTPServer(server *http.Server) error {
	fmt.Printf("Running server on %s\n", server.Addr)
	return server.ListenAndServe()
}

// gracefulShutdown - graceful shutdown приложения
func gracefulShutdown(service *services.StorageService, servers ...interface{}) error {
	fmt.Println("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем серверы
	for i, server := range servers {
		switch s := server.(type) {
		case *http.Server:
			if s != nil {
				if err := s.Shutdown(ctx); err != nil {
					fmt.Printf("HTTP server %d shutdown error: %v\n", i, err)
				} else {
					fmt.Printf("HTTP server %d stopped\n", i)
				}
			}
		case *grpc.Server:
			if s != nil {
				fmt.Printf("Stopping gRPC server...\n")
				s.GracefulStop()
				fmt.Printf("gRPC server stopped\n")
			}
		}
	}

	// Останавливаем сервис
	service.Shutdown()
	fmt.Println("Service shutdown completed")

	fmt.Println("Graceful shutdown finished")
	return nil
}
