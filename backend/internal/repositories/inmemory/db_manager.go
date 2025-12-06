package inmemory

// InMemoryRepositories - структура содержащая все репозитории
type DatabaseManager struct {
	Users       *InMemoryUsersRepo
	Binaries    *InMemoryBinariesRepo
	Cards       *InMemoryCardsRepo
	Credentials *InMemoryCredentialsRepo
	Texts       *InMemoryTextsRepo
}

// NewDatabaseManager - создание менеджера репозиториев
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		Users:       NewInMemoryUsersRepo(),
		Binaries:    NewInMemoryBinariesRepo(),
		Cards:       NewInMemoryCardsRepo(),
		Credentials: NewInMemoryCredentialsRepo(),
		Texts:       NewInMemoryTextsRepo(),
	}
}
