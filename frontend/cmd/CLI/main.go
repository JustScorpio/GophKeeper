package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/frontend/internal/repositories/sqlite"
	"github.com/JustScorpio/GophKeeper/frontend/internal/services"
)

var (
	// build-переменные заполняемые с помощью ldflags -X
	buildVersion = "1.0"
	buildDate    = time.Now().Format("January 2, 2006")
)

// configContent - содержимое конфигурационного файла
//
//go:embed config.json
var configContent []byte

// DBConfiguration - из confog.json
type AppConfiguration struct {
	DbPath     string `json:"db_path"`
	ServerAddr string `json:"server_addr"`
}

type App struct {
	dbManager    *sqlite.DatabaseManager
	authService  *services.AuthService
	apiClient    *services.APIClient
	localStorage *services.StorageService
	syncService  *services.SyncService
	appService   *services.GophkeeperService
	isLoggedIn   bool
	currentUser  string
}

func main() {
	fmt.Printf("%s v%s\n", "GophKeeper", buildVersion)
	fmt.Println("==========================")

	// Инициализация приложения
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.shutdown()

	// // Обработка сигналов для graceful shutdown
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// go func() {
	// 	<-sigChan
	// 	fmt.Println("\nReceived shutdown signal")
	// 	os.Exit(0)
	// }()

	// Основной цикл приложения
	app.run()
}

func initializeApp() (*App, error) {
	var conf AppConfiguration
	if err := json.Unmarshal(configContent, &conf); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	fmt.Printf("Using database: %s\n", conf.DbPath)
	fmt.Printf("Connecting to server: %s\n", conf.ServerAddr)

	// Инициализация базы данных
	dbManager, err := sqlite.NewDatabaseManager(conf.DbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Инициализация сервисов
	authService := services.NewAuthService(conf.ServerAddr)

	apiClient := services.NewAPIClient(conf.ServerAddr)

	localStorage := services.NewStorageService(
		dbManager.BinariesRepo,
		dbManager.CardsRepo,
		dbManager.CredentialsRepo,
		dbManager.TextsRepo,
	)

	syncService := services.NewSyncService(apiClient, localStorage)

	appService := services.NewGophkeeperService(authService, apiClient, localStorage, syncService)

	return &App{
		dbManager:    dbManager,
		authService:  authService,
		apiClient:    apiClient,
		localStorage: localStorage,
		syncService:  syncService,
		appService:   appService,
		isLoggedIn:   false,
		currentUser:  "",
	}, nil
}

func (a *App) shutdown() {
	if a.dbManager != nil {
		a.dbManager.Close()
	}
	fmt.Println("\nGoodbye!")
}

func (a *App) run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		a.showMainMenu()

		fmt.Print("\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.handleLogin(reader)
		case "2":
			a.handleRegister(reader)
		case "3":
			if a.isLoggedIn {
				a.handleDataMenu(reader)
			} else {
				fmt.Println("Please login first!")
			}
		case "4":
			if a.isLoggedIn {
				a.handleSync()
			} else {
				fmt.Println("Please login first!")
			}
		case "5":
			if a.isLoggedIn {
				a.handleLogout()
			} else {
				fmt.Println("You are not logged in!")
			}
		case "6":
			fmt.Println("Exiting...")
			return
		case "help":
			a.showHelp()
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func (a *App) showMainMenu() {
	fmt.Println("\n=== Main Menu ===")
	if a.isLoggedIn {
		fmt.Printf("Logged in as: %s\n", a.currentUser)
		fmt.Println("1. Login (switch user)")
		fmt.Println("2. Register")
		fmt.Println("3. Manage Data")
		fmt.Println("4. Sync Data")
		fmt.Println("5. Logout")
		fmt.Println("6. Exit")
	} else {
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Manage Data (requires login)")
		fmt.Println("4. Sync Data (requires login)")
		fmt.Println("5. Logout")
		fmt.Println("6. Exit")
	}
}

func (a *App) showHelp() {
	fmt.Println("\n=== Available Commands ===")
	fmt.Println("login    - Login to your account")
	fmt.Println("register - Create a new account")
	fmt.Println("data     - Manage your data (binaries, cards, etc.)")
	fmt.Println("sync     - Synchronize data with server")
	fmt.Println("logout   - Logout from current account")
	fmt.Println("exit     - Exit the application")
	fmt.Println("help     - Show this help message")
}

func (a *App) handleLogin(reader *bufio.Reader) {
	fmt.Println("\n=== Login ===")

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("Logging in... ")
	user, err := a.appService.Login(ctx, username, password)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
		return
	}

	a.isLoggedIn = true
	a.currentUser = user.Login
	fmt.Println("SUCCESS")
	fmt.Printf("Welcome, %s!\n", user.Login)
}

func (a *App) handleRegister(reader *bufio.Reader) {
	fmt.Println("\n=== Register ===")

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Confirm Password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)

	if password != confirmPassword {
		fmt.Println("Passwords do not match!")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Print("Registering... ")
	user, err := a.appService.Register(ctx, username, password)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
		return
	}

	fmt.Println("SUCCESS")
	fmt.Printf("Account created for %s. You can now login.\n", user.Login)
}

func (a *App) handleLogout() {
	fmt.Printf("\nLogging out %s...\n", a.currentUser)
	a.isLoggedIn = false
	a.currentUser = ""
	fmt.Println("Logged out successfully.")
}

func (a *App) handleSync() {
	fmt.Println("\n=== Sync Data ===")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Print("Synchronizing with server... ")
	err := a.appService.ForceSync(ctx)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Println("All data synchronized successfully.")
	}
}

func (a *App) handleDataMenu(reader *bufio.Reader) {
	for {
		fmt.Println("\n=== Data Management ===")
		fmt.Println("1. Binary Data")
		fmt.Println("2. Card Data")
		fmt.Println("3. Credentials")
		fmt.Println("4. Text Data")
		fmt.Println("5. Back to Main Menu")

		fmt.Print("\nSelect data type: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.handleBinaryData(reader)
		case "2":
			a.handleCardData(reader)
		case "3":
			a.handleCredentials(reader)
		case "4":
			a.handleTextData(reader)
		case "5":
			return
		default:
			fmt.Println("Invalid selection")
		}
	}
}

func (a *App) handleBinaryData(reader *bufio.Reader) {
	for {
		fmt.Println("\n=== Binary Data ===")
		fmt.Println("1. List all binaries")
		fmt.Println("2. Create new binary")
		fmt.Println("3. View binary")
		fmt.Println("4. Update binary")
		fmt.Println("5. Delete binary")
		fmt.Println("6. Back")

		fmt.Print("\nSelect action: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.listBinaries()
		case "2":
			a.createBinary(reader)
		case "3":
			a.viewBinary(reader)
		case "4":
			a.updateBinary(reader)
		case "5":
			a.deleteBinary(reader)
		case "6":
			return
		default:
			fmt.Println("Invalid selection")
		}
	}
}

func (a *App) listBinaries() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaries, err := a.appService.GetAllBinaries(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(binaries) == 0 {
		fmt.Println("No binary data found.")
		return
	}

	fmt.Println("\n=== Binary Data List ===")
	for i, binary := range binaries {
		fmt.Printf("%d. ID: %s, Metadata: %s, Size: %d bytes\n",
			i+1, binary.ID, binary.Metadata, len(binary.Data))
	}
}

func (a *App) createBinary(reader *bufio.Reader) {
	fmt.Println("\n=== Create Binary Data ===")

	fmt.Print("Enter metadata (description): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	fmt.Print("Enter file path to load binary data (or press Enter to skip): ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	var data []byte
	if filePath != "" {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		data = fileData
		fmt.Printf("Loaded %d bytes from file\n", len(data))
	} else {
		fmt.Print("Enter base64 encoded data: ")
		base64Data, _ := reader.ReadString('\n')
		base64Data = strings.TrimSpace(base64Data)

		// В реальном приложении нужно декодировать base64
		data = []byte(base64Data) // временно, для демонстрации
		fmt.Printf("Using %d bytes of data\n", len(data))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dto := &dtos.NewBinaryData{
		Data:            data,
		NewSecureEntity: dtos.NewSecureEntity{Metadata: metadata},
	}

	fmt.Print("Creating binary data... ")
	binary, err := a.appService.CreateBinary(ctx, dto)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Created binary with ID: %s\n", binary.ID)
	}
}

func (a *App) viewBinary(reader *bufio.Reader) {
	fmt.Print("\nEnter binary ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binary, err := a.appService.GetBinary(ctx, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if binary == nil {
		fmt.Println("Binary not found")
		return
	}

	fmt.Println("\n=== Binary Details ===")
	fmt.Printf("ID: %s\n", binary.ID)
	fmt.Printf("Metadata: %s\n", binary.Metadata)
	fmt.Printf("Size: %d bytes\n", len(binary.Data))

	fmt.Print("\nSave to file? (y/n): ")
	saveChoice, _ := reader.ReadString('\n')
	saveChoice = strings.TrimSpace(strings.ToLower(saveChoice))

	if saveChoice == "y" || saveChoice == "yes" {
		fmt.Print("Enter filename: ")
		filename, _ := reader.ReadString('\n')
		filename = strings.TrimSpace(filename)

		if err := os.WriteFile(filename, binary.Data, 0644); err != nil {
			fmt.Printf("Error saving file: %v\n", err)
		} else {
			fmt.Printf("Data saved to %s\n", filename)
		}
	}
}

func (a *App) updateBinary(reader *bufio.Reader) {
	fmt.Print("\nEnter binary ID to update: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Сначала получаем текущие данные
	existing, err := a.appService.GetBinary(ctx, id)
	if err != nil {
		fmt.Printf("Error getting binary: %v\n", err)
		return
	}

	if existing == nil {
		fmt.Println("Binary not found")
		return
	}

	fmt.Println("\n=== Current Binary Data ===")
	fmt.Printf("ID: %s\n", existing.ID)
	fmt.Printf("Metadata: %s\n", existing.Metadata)
	fmt.Printf("Size: %d bytes\n", len(existing.Data))

	fmt.Println("\n=== Update Options ===")
	fmt.Println("1. Update metadata only")
	fmt.Println("2. Update data from file")
	fmt.Println("3. Update both metadata and data")
	fmt.Print("\nSelect option: ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	var newMetadata string
	var newData []byte

	switch option {
	case "1": // Только метаданные
		fmt.Print("New metadata: ")
		newMetadata, _ = reader.ReadString('\n')
		newMetadata = strings.TrimSpace(newMetadata)
		newData = existing.Data

	case "2": // Только данные
		fmt.Print("Enter file path with new binary data: ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)

		if filePath == "" {
			fmt.Println("File path cannot be empty")
			return
		}

		fileData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		newMetadata = existing.Metadata
		newData = fileData
		fmt.Printf("Loaded %d bytes from file\n", len(fileData))

	case "3": // И метаданные, и данные
		fmt.Print("New metadata: ")
		newMetadata, _ = reader.ReadString('\n')
		newMetadata = strings.TrimSpace(newMetadata)

		fmt.Print("Enter file path with new binary data: ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)

		if filePath == "" {
			fmt.Println("File path cannot be empty")
			return
		}

		fileData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		newData = fileData
		fmt.Printf("Loaded %d bytes from file\n", len(fileData))

	default:
		fmt.Println("Invalid option")
		return
	}

	// Проверка размера данных
	if len(newData) > 10*1024*1024 { // 10MB limit
		fmt.Println("Error: File too large (max 10MB)")
		return
	}

	updatedBinary := &entities.BinaryData{
		Data:         newData,
		SecureEntity: entities.SecureEntity{Metadata: newMetadata},
	}

	fmt.Print("\nUpdating binary... ")
	updated, err := a.appService.UpdateBinary(ctx, updatedBinary)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Updated binary with ID: %s\n", updated.ID)
		fmt.Printf("New size: %d bytes\n", len(updated.Data))
	}
}

func (a *App) deleteBinary(reader *bufio.Reader) {
	fmt.Print("\nEnter binary ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Are you sure? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("Deleting binary... ")
	err := a.appService.DeleteBinary(ctx, id)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Println("Binary deleted")
	}
}

func (a *App) handleCardData(reader *bufio.Reader) {
	for {
		fmt.Println("\n=== Card Data Management ===")
		fmt.Println("1. List all cards")
		fmt.Println("2. Create new card")
		fmt.Println("3. View card")
		fmt.Println("4. Update card")
		fmt.Println("5. Delete card")
		fmt.Println("6. Back")

		fmt.Print("\nSelect action: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.listCards()
		case "2":
			a.createCard(reader)
		case "3":
			a.viewCard(reader)
		case "4":
			a.updateCard(reader)
		case "5":
			a.deleteCard(reader)
		case "6":
			return
		default:
			fmt.Println("Invalid selection")
		}
	}
}

func (a *App) listCards() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cards, err := a.appService.GetAllCards(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(cards) == 0 {
		fmt.Println("No cards found.")
		return
	}

	fmt.Println("\n=== Card List ===")
	for i, card := range cards {
		// Маскируем номер карты для безопасности
		maskedNumber := maskCardNumber(card.Number)
		fmt.Printf("%d. ID: %s, Card: %s (%s), Expires: %s\n",
			i+1, card.ID, maskedNumber, card.CardHolder, card.ExpirationDate)
	}
}

func maskCardNumber(number string) string {
	if len(number) < 4 {
		return number
	}
	lastFour := number[len(number)-4:]
	return "**** **** **** " + lastFour
}

func (a *App) createCard(reader *bufio.Reader) {
	fmt.Println("\n=== Create New Card ===")

	fmt.Print("Card number: ")
	number, _ := reader.ReadString('\n')
	number = strings.TrimSpace(number)

	fmt.Print("Card holder name: ")
	cardHolder, _ := reader.ReadString('\n')
	cardHolder = strings.TrimSpace(cardHolder)

	fmt.Print("Expiration date (MM/YY): ")
	expirationDate, _ := reader.ReadString('\n')
	expirationDate = strings.TrimSpace(expirationDate)

	fmt.Print("CVV: ")
	cvv, _ := reader.ReadString('\n')
	cvv = strings.TrimSpace(cvv)

	fmt.Print("Metadata (description): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dto := &dtos.NewCardInformation{
		Number:          number,
		CardHolder:      cardHolder,
		ExpirationDate:  expirationDate,
		CVV:             cvv,
		NewSecureEntity: dtos.NewSecureEntity{Metadata: metadata},
	}

	fmt.Print("Creating card... ")
	card, err := a.appService.CreateCard(ctx, dto)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Created card with ID: %s\n", card.ID)
	}
}

func (a *App) viewCard(reader *bufio.Reader) {
	fmt.Print("\nEnter card ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	card, err := a.appService.GetCard(ctx, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if card == nil {
		fmt.Println("Card not found")
		return
	}

	fmt.Println("\n=== Card Details ===")
	fmt.Printf("ID: %s\n", card.ID)
	fmt.Printf("Card number: %s\n", maskCardNumber(card.Number))
	fmt.Printf("Card holder: %s\n", card.CardHolder)
	fmt.Printf("Expiration date: %s\n", card.ExpirationDate)
	fmt.Printf("CVV: %s\n", card.CVV)
	fmt.Printf("Metadata: %s\n", card.Metadata)
}

func (a *App) updateCard(reader *bufio.Reader) {
	fmt.Print("\nEnter card ID to update: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existing, err := a.appService.GetCard(ctx, id)
	if err != nil {
		fmt.Printf("Error getting card: %v\n", err)
		return
	}

	if existing == nil {
		fmt.Println("Card not found")
		return
	}

	fmt.Println("\n=== Current Card Data ===")
	fmt.Printf("ID: %s\n", existing.ID)
	fmt.Printf("Card holder: %s\n", existing.CardHolder)
	fmt.Printf("Expiration date: %s\n", existing.ExpirationDate)
	fmt.Printf("Metadata: %s\n", existing.Metadata)

	fmt.Println("\n=== Update Card ===")

	fmt.Print("Card number (press Enter to keep current): ")
	number, _ := reader.ReadString('\n')
	number = strings.TrimSpace(number)
	if number == "" {
		number = existing.Number
	}

	fmt.Print("Card holder name (press Enter to keep current): ")
	cardHolder, _ := reader.ReadString('\n')
	cardHolder = strings.TrimSpace(cardHolder)
	if cardHolder == "" {
		cardHolder = existing.CardHolder
	}

	fmt.Print("Expiration date MM/YY (press Enter to keep current): ")
	expirationDate, _ := reader.ReadString('\n')
	expirationDate = strings.TrimSpace(expirationDate)
	if expirationDate == "" {
		expirationDate = existing.ExpirationDate
	}

	fmt.Print("CVV (press Enter to keep current): ")
	cvv, _ := reader.ReadString('\n')
	cvv = strings.TrimSpace(cvv)
	if cvv == "" {
		cvv = existing.CVV
	}

	fmt.Print("Metadata (press Enter to keep current): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)
	if metadata == "" {
		metadata = existing.Metadata
	}

	updatedCard := &entities.CardInformation{
		Number:         number,
		CardHolder:     cardHolder,
		ExpirationDate: expirationDate,
		CVV:            cvv,
		SecureEntity:   entities.SecureEntity{Metadata: metadata},
	}

	fmt.Print("\nUpdating card... ")
	updated, err := a.appService.UpdateCard(ctx, updatedCard)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Updated card with ID: %s\n", updated.ID)
	}
}

func (a *App) deleteCard(reader *bufio.Reader) {
	fmt.Print("\nEnter card ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Are you sure? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("Deleting card... ")
	err := a.appService.DeleteCard(ctx, id)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Println("Card deleted")
	}
}

func (a *App) handleCredentials(reader *bufio.Reader) {
	for {
		fmt.Println("\n=== Credentials Management ===")
		fmt.Println("1. List all credentials")
		fmt.Println("2. Create new credentials")
		fmt.Println("3. View credentials")
		fmt.Println("4. Update credentials")
		fmt.Println("5. Delete credentials")
		fmt.Println("6. Back")

		fmt.Print("\nSelect action: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.listCredentials()
		case "2":
			a.createCredentials(reader)
		case "3":
			a.viewCredentials(reader)
		case "4":
			a.updateCredentials(reader)
		case "5":
			a.deleteCredentials(reader)
		case "6":
			return
		default:
			fmt.Println("Invalid selection")
		}
	}
}

func (a *App) listCredentials() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credentials, err := a.appService.GetAllCredentials(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(credentials) == 0 {
		fmt.Println("No credentials found.")
		return
	}

	fmt.Println("\n=== Credentials List ===")
	for i, cred := range credentials {
		fmt.Printf("%d. ID: %s, Login: %s, Metadata: %s\n",
			i+1, cred.ID, cred.Login, cred.Metadata)
	}
}

func (a *App) createCredentials(reader *bufio.Reader) {
	fmt.Println("\n=== Create New Credentials ===")

	fmt.Print("Login/Username: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Metadata (description, e.g., 'Gmail account'): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dto := &dtos.NewCredentials{
		Login:           login,
		Password:        password,
		NewSecureEntity: dtos.NewSecureEntity{Metadata: metadata},
	}

	fmt.Print("Creating credentials... ")
	creds, err := a.appService.CreateCredentials(ctx, dto)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Created credentials with ID: %s\n", creds.ID)
	}
}

func (a *App) viewCredentials(reader *bufio.Reader) {
	fmt.Print("\nEnter credentials ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	creds, err := a.appService.GetCredentials(ctx, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if creds == nil {
		fmt.Println("Credentials not found")
		return
	}

	fmt.Println("\n=== Credentials Details ===")
	fmt.Printf("ID: %s\n", creds.ID)
	fmt.Printf("Login: %s\n", creds.Login)
	fmt.Printf("Password: %s\n", creds.Password)
	fmt.Printf("Metadata: %s\n", creds.Metadata)
}

func (a *App) updateCredentials(reader *bufio.Reader) {
	fmt.Print("\nEnter credentials ID to update: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existing, err := a.appService.GetCredentials(ctx, id)
	if err != nil {
		fmt.Printf("Error getting credentials: %v\n", err)
		return
	}

	if existing == nil {
		fmt.Println("Credentials not found")
		return
	}

	fmt.Println("\n=== Current Credentials ===")
	fmt.Printf("ID: %s\n", existing.ID)
	fmt.Printf("Login: %s\n", existing.Login)
	fmt.Printf("Metadata: %s\n", existing.Metadata)

	fmt.Println("\n=== Update Credentials ===")

	fmt.Print("Login (press Enter to keep current): ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)
	if login == "" {
		login = existing.Login
	}

	fmt.Print("Password (press Enter to keep current): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	if password == "" {
		password = existing.Password
	}

	fmt.Print("Metadata (press Enter to keep current): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)
	if metadata == "" {
		metadata = existing.Metadata
	}

	updatedCreds := &entities.Credentials{
		Login:        login,
		Password:     password,
		SecureEntity: entities.SecureEntity{Metadata: metadata},
	}

	fmt.Print("\nUpdating credentials... ")
	updated, err := a.appService.UpdateCredentials(ctx, updatedCreds)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Updated credentials with ID: %s\n", updated.ID)
	}
}

func (a *App) deleteCredentials(reader *bufio.Reader) {
	fmt.Print("\nEnter credentials ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Are you sure? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("Deleting credentials... ")
	err := a.appService.DeleteCredentials(ctx, id)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Println("Credentials deleted")
	}
}

func (a *App) handleTextData(reader *bufio.Reader) {
	for {
		fmt.Println("\n=== Text Data Management ===")
		fmt.Println("1. List all texts")
		fmt.Println("2. Create new text")
		fmt.Println("3. View text")
		fmt.Println("4. Update text")
		fmt.Println("5. Delete text")
		fmt.Println("6. Back")

		fmt.Print("\nSelect action: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.listTexts()
		case "2":
			a.createText(reader)
		case "3":
			a.viewText(reader)
		case "4":
			a.updateText(reader)
		case "5":
			a.deleteText(reader)
		case "6":
			return
		default:
			fmt.Println("Invalid selection")
		}
	}
}

func (a *App) listTexts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	texts, err := a.appService.GetAllTexts(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(texts) == 0 {
		fmt.Println("No texts found.")
		return
	}

	fmt.Println("\n=== Text List ===")
	for i, text := range texts {
		// Обрезаем текст для предпросмотра
		preview := text.Data
		if len(preview) > 50 {
			preview = preview[:47] + "..."
		}
		fmt.Printf("%d. ID: %s, Metadata: %s\n", i+1, text.ID, text.Metadata)
		fmt.Printf("   Preview: %s\n", preview)
	}
}

func (a *App) createText(reader *bufio.Reader) {
	fmt.Println("\n=== Create New Text ===")

	fmt.Print("Metadata (description): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	fmt.Println("Enter text content (end with empty line or Ctrl+D):")
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(line) == "" {
			break
		}
		lines = append(lines, line)
	}

	content := strings.Join(lines, "")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dto := &dtos.NewTextData{
		Data:            content,
		NewSecureEntity: dtos.NewSecureEntity{Metadata: metadata},
	}

	fmt.Print("Creating text... ")
	text, err := a.appService.CreateText(ctx, dto)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Created text with ID: %s\n", text.ID)
		fmt.Printf("Length: %d characters\n", len(text.Data))
	}
}

func (a *App) viewText(reader *bufio.Reader) {
	fmt.Print("\nEnter text ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	text, err := a.appService.GetText(ctx, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if text == nil {
		fmt.Println("Text not found")
		return
	}

	fmt.Println("\n=== Text Details ===")
	fmt.Printf("ID: %s\n", text.ID)
	fmt.Printf("Metadata: %s\n", text.Metadata)
	fmt.Printf("Length: %d characters\n", len(text.Data))
	fmt.Printf("\n=== Content ===\n%s\n", text.Data)
}

func (a *App) updateText(reader *bufio.Reader) {
	fmt.Print("\nEnter text ID to update: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existing, err := a.appService.GetText(ctx, id)
	if err != nil {
		fmt.Printf("Error getting text: %v\n", err)
		return
	}

	if existing == nil {
		fmt.Println("Text not found")
		return
	}

	fmt.Println("\n=== Current Text ===")
	fmt.Printf("ID: %s\n", existing.ID)
	fmt.Printf("Metadata: %s\n", existing.Metadata)
	fmt.Printf("Length: %d characters\n", len(existing.Data))

	fmt.Println("\n=== Update Options ===")
	fmt.Println("1. Update metadata only")
	fmt.Println("2. Update content only")
	fmt.Println("3. Update both metadata and content")
	fmt.Print("\nSelect option: ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	var newMetadata string
	var newContent string

	switch option {
	case "1": // Только метаданные
		fmt.Print("New metadata: ")
		newMetadata, _ = reader.ReadString('\n')
		newMetadata = strings.TrimSpace(newMetadata)
		newContent = existing.Data

	case "2": // Только содержимое
		fmt.Println("Enter new text content (end with empty line):")
		var lines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil || strings.TrimSpace(line) == "" {
				break
			}
			lines = append(lines, line)
		}

		newMetadata = existing.Metadata
		newContent = strings.Join(lines, "")
		fmt.Printf("New length: %d characters\n", len(newContent))

	case "3": // И метаданные, и содержимое
		fmt.Print("New metadata: ")
		newMetadata, _ = reader.ReadString('\n')
		newMetadata = strings.TrimSpace(newMetadata)

		fmt.Println("Enter new text content (end with empty line):")
		var lines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil || strings.TrimSpace(line) == "" {
				break
			}
			lines = append(lines, line)
		}

		newContent = strings.Join(lines, "")
		fmt.Printf("New length: %d characters\n", len(newContent))

	default:
		fmt.Println("Invalid option")
		return
	}

	// Проверка размера (если есть лимит)
	if len(newContent) > 1*1024*1024 { // 1MB limit
		fmt.Println("Error: Text too large (max 1MB)")
		return
	}

	updatedText := &entities.TextData{
		Data:         newContent,
		SecureEntity: entities.SecureEntity{Metadata: newMetadata},
	}

	fmt.Print("\nUpdating text... ")
	updated, err := a.appService.UpdateText(ctx, updatedText)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Printf("Updated text with ID: %s\n", updated.ID)
		fmt.Printf("New length: %d characters\n", len(updated.Data))
	}
}

func (a *App) deleteText(reader *bufio.Reader) {
	fmt.Print("\nEnter text ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Are you sure? (yes/no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("Deleting text... ")
	err := a.appService.DeleteText(ctx, id)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
		fmt.Println("Text deleted")
	}
}
