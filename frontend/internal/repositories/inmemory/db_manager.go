// inmemory - репозиторий хранящий данные воперативной памяти
package inmemory

// InMemoryRepositories - структура содержащая все репозитории
type DatabaseManager struct {
	BinariesRepo    *InMemoryBinariesRepo
	CardsRepo       *InMemoryCardsRepo
	CredentialsRepo *InMemoryCredentialsRepo
	TextsRepo       *InMemoryTextsRepo
}

// NewDatabaseManager - создание менеджера репозиториев
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		BinariesRepo:    NewInMemoryBinariesRepo(),
		CardsRepo:       NewInMemoryCardsRepo(),
		CredentialsRepo: NewInMemoryCredentialsRepo(),
		TextsRepo:       NewInMemoryTextsRepo(),
	}
}
