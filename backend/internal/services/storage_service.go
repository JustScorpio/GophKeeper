// Пакет services содержит структуры и методы, реализующие бизнес-логику приложения
package services

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/customerrors"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/backend/internal/repositories"
)

// StorageService - сервис для взаимодействия с хранилищем
type StorageService struct {
	binariesRepo    repositories.IRepository[entities.BinaryData, dtos.NewBinaryData]
	cardsRepo       repositories.IRepository[entities.CardInformation, dtos.NewCardInformation]
	credentialsRepo repositories.IRepository[entities.Credentials, dtos.NewCredentials]
	textsRepo       repositories.IRepository[entities.TextData, dtos.NewTextData]
	usersRepo       repositories.IRepository[entities.User, dtos.NewUser]

	taskQueue      chan Task // канал-очередь задач
	tasksInProcess sync.WaitGroup

	isShuttingDown atomic.Bool //Использование вместо Bool помогает избежать гонки данных при её обновлении
}

// TaskType - алиас вокруг int, типы задач задач
type TaskType int

const (
	TaskGetAll TaskType = iota
	TaskGet
	TaskCreate
	TaskUpdate
	TaskDelete
)

// EntityType - типы сущностей
type EntityType int

const (
	EntityUser EntityType = iota
	EntityBinary
	EntityCard
	EntityCredentials
	EntityText
)

// Task - задача в очереди задач на обработку сервисом
type Task struct {
	Context    context.Context
	Payload    interface{}
	ResultCh   chan TaskResult
	TaskType   TaskType
	EntityType EntityType
}

// TaskResult - результат обработки задачи Task
type TaskResult struct {
	Result interface{}
	Err    error
}

// NewStorageService - инициализация сервиса
func NewStorageService(usersRepo repositories.IRepository[entities.User, dtos.NewUser],
	binariesRepo repositories.IRepository[entities.BinaryData, dtos.NewBinaryData],
	cardsRepo repositories.IRepository[entities.CardInformation, dtos.NewCardInformation],
	credentialsRepo repositories.IRepository[entities.Credentials, dtos.NewCredentials],
	textsRepo repositories.IRepository[entities.TextData, dtos.NewTextData]) *StorageService {
	service := &StorageService{
		usersRepo:       usersRepo,
		binariesRepo:    binariesRepo,
		cardsRepo:       cardsRepo,
		credentialsRepo: credentialsRepo,
		textsRepo:       textsRepo,
		taskQueue:       make(chan Task, 256),
	}

	go service.taskProcessor()

	return service
}

// taskProcessor - обработчик очереди задач в составе StorageService
func (s *StorageService) taskProcessor() {
	for task := range s.taskQueue {
		var result interface{}
		var err error

		//Если происходит shutdown - прерываем задачи которые уже стоят в очереди
		if s.isShuttingDown.Load() {
			if task.ResultCh != nil {
				task.ResultCh <- TaskResult{
					Err: customerrors.ServiceUnavailableError,
				}
				close(task.ResultCh)
			}

			continue
		}

		switch task.EntityType {
		case EntityBinary:
			result, err = s.processBinaryTask(task)
		case EntityCard:
			result, err = s.processCardTask(task)
		case EntityCredentials:
			result, err = s.processCredentialsTask(task)
		case EntityText:
			result, err = s.processTextTask(task)
		case EntityUser:
			result, err = s.processUserTask(task)
		}

		if task.ResultCh != nil {
			task.ResultCh <- TaskResult{
				Result: result,
				Err:    err,
			}
			close(task.ResultCh)
		}
	}
}

func (s *StorageService) processBinaryTask(task Task) (interface{}, error) {
	switch task.TaskType {
	case TaskCreate:
		dto := task.Payload.(*dtos.NewBinaryData)
		return s.binariesRepo.Create(task.Context, dto)
	case TaskGet:
		id := task.Payload.(string)
		return s.binariesRepo.Get(task.Context, id)
	case TaskGetAll:
		return s.binariesRepo.GetAll(task.Context)
	case TaskUpdate:
		entity := task.Payload.(*entities.BinaryData)
		return s.binariesRepo.Update(task.Context, entity)
	case TaskDelete:
		id := task.Payload.(string)
		return nil, s.binariesRepo.Delete(task.Context, id)
	default:
		return nil, customerrors.UnsupportedOperation
	}
}

func (s *StorageService) processCardTask(task Task) (interface{}, error) {
	switch task.TaskType {
	case TaskCreate:
		dto := task.Payload.(*dtos.NewCardInformation)
		return s.cardsRepo.Create(task.Context, dto)
	case TaskGet:
		id := task.Payload.(string)
		return s.cardsRepo.Get(task.Context, id)
	case TaskGetAll:
		return s.cardsRepo.GetAll(task.Context)
	case TaskUpdate:
		entity := task.Payload.(*entities.CardInformation)
		return s.cardsRepo.Update(task.Context, entity)
	case TaskDelete:
		id := task.Payload.(string)
		return nil, s.cardsRepo.Delete(task.Context, id)
	default:
		return nil, customerrors.UnsupportedOperation
	}
}

func (s *StorageService) processCredentialsTask(task Task) (interface{}, error) {
	switch task.TaskType {
	case TaskCreate:
		dto := task.Payload.(*dtos.NewCredentials)
		return s.credentialsRepo.Create(task.Context, dto)
	case TaskGet:
		id := task.Payload.(string)
		return s.credentialsRepo.Get(task.Context, id)
	case TaskGetAll:
		return s.credentialsRepo.GetAll(task.Context)
	case TaskUpdate:
		entity := task.Payload.(*entities.Credentials)
		return s.credentialsRepo.Update(task.Context, entity)
	case TaskDelete:
		id := task.Payload.(string)
		return nil, s.credentialsRepo.Delete(task.Context, id)
	default:
		return nil, customerrors.UnsupportedOperation
	}
}

func (s *StorageService) processTextTask(task Task) (interface{}, error) {
	switch task.TaskType {
	case TaskCreate:
		dto := task.Payload.(*dtos.NewTextData)
		return s.textsRepo.Create(task.Context, dto)
	case TaskGet:
		id := task.Payload.(string)
		return s.textsRepo.Get(task.Context, id)
	case TaskGetAll:
		return s.textsRepo.GetAll(task.Context)
	case TaskUpdate:
		entity := task.Payload.(*entities.TextData)
		return s.textsRepo.Update(task.Context, entity)
	case TaskDelete:
		id := task.Payload.(string)
		return nil, s.textsRepo.Delete(task.Context, id)
	default:
		return nil, customerrors.UnsupportedOperation
	}
}

func (s *StorageService) processUserTask(task Task) (interface{}, error) {
	switch task.TaskType {
	case TaskCreate:
		dto := task.Payload.(*dtos.NewUser)
		return s.createUser(task.Context, dto)
	case TaskGet:
		id := task.Payload.(string)
		return s.usersRepo.Get(task.Context, id)
	case TaskGetAll:
		return s.usersRepo.GetAll(task.Context)
	case TaskUpdate:
		entity := task.Payload.(*entities.User)
		return s.updateUser(task.Context, entity)
	case TaskDelete:
		id := task.Payload.(string)
		return nil, s.usersRepo.Delete(task.Context, id)
	default:
		return nil, customerrors.UnsupportedOperation
	}
}

// enqueueTask - поставить задачу в очередь
func (s *StorageService) enqueueTask(task Task) (interface{}, error) {
	// Проверяем, не начался ли shutdown
	if s.isShuttingDown.Load() {
		return nil, customerrors.ServiceUnavailableError
	}

	if task.ResultCh == nil {
		task.ResultCh = make(chan TaskResult, 1)
	}

	s.tasksInProcess.Add(1) // Увеличиваем счетчик
	s.taskQueue <- task

	select {
	case <-task.Context.Done():
		s.tasksInProcess.Done() // Уменьшаем счётчик
		return nil, task.Context.Err()
	case res := <-task.ResultCh:
		s.tasksInProcess.Done() // Уменьшаем счётчик
		return res.Result, res.Err
	}
}

// CreateUser - создать пользователя
func (s *StorageService) CreateUser(ctx context.Context, newUser dtos.NewUser) (*entities.User, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskCreate,
		EntityType: EntityUser,
		Context:    ctx,
		Payload:    &newUser,
	})

	return res.(*entities.User), err
}

// GetUser - получить пользователя по логину
func (s *StorageService) GetUser(ctx context.Context, login string) (*entities.User, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGet,
		EntityType: EntityUser,
		Context:    ctx,
		Payload:    login,
	})

	return res.(*entities.User), err
}

// GetAllUsers - получить всех пользователей
func (s *StorageService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGetAll,
		EntityType: EntityUser,
		Context:    ctx,
	})

	return res.([]entities.User), err
}

// UpdateUser - изменить пользователя
func (s *StorageService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskUpdate,
		EntityType: EntityUser,
		Context:    ctx,
		Payload:    user,
	})

	return res.(*entities.User), err
}

// DeleteUser - удалить пользователя
func (s *StorageService) DeleteUser(ctx context.Context, login string) error {
	_, err := s.enqueueTask(Task{
		TaskType:   TaskDelete,
		EntityType: EntityUser,
		Context:    ctx,
		Payload:    login,
	})

	return err
}

// CreateBinary - создать бинарные данные
func (s *StorageService) CreateBinary(ctx context.Context, newBinary *dtos.NewBinaryData) (*entities.BinaryData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskCreate,
		EntityType: EntityBinary,
		Context:    ctx,
		Payload:    newBinary,
	})
	return res.(*entities.BinaryData), err
}

// GetBinary - получить бинарные данные
func (s *StorageService) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGet,
		EntityType: EntityBinary,
		Context:    ctx,
		Payload:    id,
	})
	return res.(*entities.BinaryData), err
}

// GetAllBinaries - получить все бинарные данные
func (s *StorageService) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGetAll,
		EntityType: EntityBinary,
		Context:    ctx,
	})
	return res.([]entities.BinaryData), err
}

// UpdateBinary - изменить бинарные данные
func (s *StorageService) UpdateBinary(ctx context.Context, binary *entities.BinaryData) (*entities.BinaryData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskUpdate,
		EntityType: EntityBinary,
		Context:    ctx,
		Payload:    binary,
	})
	return res.(*entities.BinaryData), err
}

// DeleteBinary - удалить бинарные данные
func (s *StorageService) DeleteBinary(ctx context.Context, id string) error {
	_, err := s.enqueueTask(Task{
		TaskType:   TaskDelete,
		EntityType: EntityBinary,
		Context:    ctx,
		Payload:    id,
	})
	return err
}

// CreateCard - создать данные банковской карты
func (s *StorageService) CreateCard(ctx context.Context, newCard *dtos.NewCardInformation) (*entities.CardInformation, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskCreate,
		EntityType: EntityCard,
		Context:    ctx,
		Payload:    newCard,
	})
	return res.(*entities.CardInformation), err
}

// GetCard - получить данные банковской карты
func (s *StorageService) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGet,
		EntityType: EntityCard,
		Context:    ctx,
		Payload:    id,
	})
	return res.(*entities.CardInformation), err
}

// GetAllCards - получить все данные банковской карты
func (s *StorageService) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGetAll,
		EntityType: EntityCard,
		Context:    ctx,
	})
	return res.([]entities.CardInformation), err
}

// UpdateCard - изменить данные банковской карты
func (s *StorageService) UpdateCard(ctx context.Context, card *entities.CardInformation) (*entities.CardInformation, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskUpdate,
		EntityType: EntityCard,
		Context:    ctx,
		Payload:    card,
	})
	return res.(*entities.CardInformation), err
}

// DeleteCard - удалить данные банковской карты
func (s *StorageService) DeleteCard(ctx context.Context, id string) error {
	_, err := s.enqueueTask(Task{
		TaskType:   TaskDelete,
		EntityType: EntityCard,
		Context:    ctx,
		Payload:    id,
	})
	return err
}

// CreateCredentials - создать учётные данные
func (s *StorageService) CreateCredentials(ctx context.Context, newCreds *dtos.NewCredentials) (*entities.Credentials, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskCreate,
		EntityType: EntityCredentials,
		Context:    ctx,
		Payload:    newCreds,
	})
	return res.(*entities.Credentials), err
}

// GetCredentials - получить учётные данные
func (s *StorageService) GetCredentials(ctx context.Context, id string) (*entities.Credentials, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGet,
		EntityType: EntityCredentials,
		Context:    ctx,
		Payload:    id,
	})
	return res.(*entities.Credentials), err
}

// GetAllCredentials - получить все учётные данные
func (s *StorageService) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGetAll,
		EntityType: EntityCredentials,
		Context:    ctx,
	})
	return res.([]entities.Credentials), err
}

// UpdateCredentials - изменить учётные данные
func (s *StorageService) UpdateCredentials(ctx context.Context, creds *entities.Credentials) (*entities.Credentials, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskUpdate,
		EntityType: EntityCredentials,
		Context:    ctx,
		Payload:    creds,
	})
	return res.(*entities.Credentials), err
}

// DeleteCredentials - удалить учётные данные
func (s *StorageService) DeleteCredentials(ctx context.Context, id string) error {
	_, err := s.enqueueTask(Task{
		TaskType:   TaskDelete,
		EntityType: EntityCredentials,
		Context:    ctx,
		Payload:    id,
	})
	return err
}

// CreateText - создать текстовые данные
func (s *StorageService) CreateText(ctx context.Context, newText *dtos.NewTextData) (*entities.TextData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskCreate,
		EntityType: EntityText,
		Context:    ctx,
		Payload:    newText,
	})
	return res.(*entities.TextData), err
}

// GetText - получить текстовые данные
func (s *StorageService) GetText(ctx context.Context, id string) (*entities.TextData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGet,
		EntityType: EntityText,
		Context:    ctx,
		Payload:    id,
	})
	return res.(*entities.TextData), err
}

// GetAllTexts - получить все текстовые данные
func (s *StorageService) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskGetAll,
		EntityType: EntityText,
		Context:    ctx,
	})
	return res.([]entities.TextData), err
}

// UpdateText - изменить текстовые данные
func (s *StorageService) UpdateText(ctx context.Context, text *entities.TextData) (*entities.TextData, error) {
	res, err := s.enqueueTask(Task{
		TaskType:   TaskUpdate,
		EntityType: EntityText,
		Context:    ctx,
		Payload:    text,
	})
	return res.(*entities.TextData), err
}

// DeleteText - удалить текстовые данные
func (s *StorageService) DeleteText(ctx context.Context, id string) error {
	_, err := s.enqueueTask(Task{
		TaskType:   TaskDelete,
		EntityType: EntityText,
		Context:    ctx,
		Payload:    id,
	})
	return err
}

// createUser - создать пользователя (инкапсулирует все проверки и бизнес-логику)
func (s *StorageService) createUser(ctx context.Context, newUser *dtos.NewUser) (*entities.User, error) {
	// Проверка наличие пользователя в БД
	existedUser, err := s.usersRepo.Get(ctx, newUser.Login)
	if err != nil {
		return nil, err
	}
	if existedUser != nil {
		return nil, customerrors.AlreadyExistsError
	}

	user, err := s.usersRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// createUser - создать пользователя (инкапсулирует все проверки и бизнес-логику)
func (s *StorageService) updateUser(ctx context.Context, user *entities.User) (*entities.User, error) {

	var curUserLogin = customcontext.GetUserID(ctx)

	if curUserLogin != user.Login {
		// Проверка наличие пользователя в БД
		existedUser, err := s.usersRepo.Get(ctx, user.Login)
		if err != nil {
			return nil, err
		}
		if existedUser != nil {
			return nil, customerrors.AlreadyExistsError
		}
	}

	user, err := s.usersRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Shutdown - инициирует graceful shutdown сервиса
func (s *StorageService) Shutdown() {
	//Помечаем сервис как завершающий работу
	s.isShuttingDown.Store(true)

	//Ждем завершения всех задач
	s.tasksInProcess.Wait()

	//Закрываем очередь
	close(s.taskQueue)
}
