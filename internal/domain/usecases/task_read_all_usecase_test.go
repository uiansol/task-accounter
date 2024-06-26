package usecases_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uiansol/task-accounter.git/internal/domain/entities"
	"github.com/uiansol/task-accounter.git/internal/domain/mocks"
	"github.com/uiansol/task-accounter.git/internal/domain/usecases"
)

func TestNewTaskReadAllUseCase(t *testing.T) {
	t.Run("should return a task read use case", func(t *testing.T) {
		taskRepositoryMock := mocks.NewTaskRepositoryInterface(t)
		encrypterMock := mocks.NewEncrypterInterface(t)
		TaskReadAllUseCase := usecases.NewTaskReadAllUseCase(taskRepositoryMock, encrypterMock)

		assert.NotNil(t, TaskReadAllUseCase)
		assert.NotNil(t, TaskReadAllUseCase.TaskRepository)
		assert.NotNil(t, TaskReadAllUseCase.Encrypter)
	})
}

func TestTaskReadAllUseCaseExecute(t *testing.T) {
	t.Run("should return all tasks when user is a manager", func(t *testing.T) {
		taskRepositoryMock := mocks.NewTaskRepositoryInterface(t)
		encrypterMock := mocks.NewEncrypterInterface(t)
		TaskReadAllUseCase := usecases.NewTaskReadAllUseCase(taskRepositoryMock, encrypterMock)

		taskRepositoryMock.On("FindAll").Return([]*entities.Task{
			{
				ID:      "task-id",
				Title:   "Task Title",
				Summary: "Task Description 1",
				OwnerID: "user-id",
				Status:  entities.Open,
			},
			{
				ID:      "task-id-2",
				Title:   "Task Title 2",
				Summary: "Task Description 2",
				OwnerID: "user-id-2",
				Status:  entities.Closed,
			},
		}, nil)
		encrypterMock.On("Decrypt", "Task Description 1").Return("Task Description 1", nil)
		encrypterMock.On("Decrypt", "Task Description 2").Return("Task Description 2", nil)

		output, err := TaskReadAllUseCase.Execute(usecases.TaskReadAllInput{
			User: entities.User{
				ID:   "user-id",
				Role: entities.UserRoleManager,
			},
		})

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Tasks, 2)

		assert.Equal(t, "task-id", output.Tasks[0].ID)
		assert.Equal(t, "Task Title", output.Tasks[0].Title)
		assert.Equal(t, "Task Description 1", output.Tasks[0].Summary)
		assert.Equal(t, "user-id", output.Tasks[0].OwnerID)
		assert.Equal(t, entities.Open, output.Tasks[0].Status)

		assert.Equal(t, "task-id-2", output.Tasks[1].ID)
		assert.Equal(t, "Task Title 2", output.Tasks[1].Title)
		assert.Equal(t, "Task Description 2", output.Tasks[1].Summary)
		assert.Equal(t, "user-id-2", output.Tasks[1].OwnerID)
		assert.Equal(t, entities.Closed, output.Tasks[1].Status)
	})

	t.Run("should return an error when an error occurs while finding all tasks", func(t *testing.T) {
		taskRepositoryMock := mocks.NewTaskRepositoryInterface(t)
		encrypterMock := mocks.NewEncrypterInterface(t)
		TaskReadAllUseCase := usecases.NewTaskReadAllUseCase(taskRepositoryMock, encrypterMock)

		taskRepositoryMock.On("FindAll").Return(nil, assert.AnError)

		output, err := TaskReadAllUseCase.Execute(usecases.TaskReadAllInput{
			User: entities.User{
				ID:   "user-id",
				Role: entities.UserRoleManager,
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, string(usecases.ErrorFindAllTasks)+": assert.AnError general error for testing", err.Error())
		assert.Equal(t, usecases.TaskReadAllOutput{}, output)
	})

	t.Run("should return tasks by user when user is not a manager", func(t *testing.T) {
		taskRepositoryMock := mocks.NewTaskRepositoryInterface(t)
		encrypterMock := mocks.NewEncrypterInterface(t)
		TaskReadAllUseCase := usecases.NewTaskReadAllUseCase(taskRepositoryMock, encrypterMock)

		taskRepositoryMock.On("FindByUserID", "user-id").Return([]*entities.Task{
			{
				ID:      "task-id",
				Title:   "Task Title",
				Summary: "Task Description 1",
				OwnerID: "user-id",
				Status:  entities.Open,
			},
			{
				ID:      "task-id-2",
				Title:   "Task Title 2",
				Summary: "Task Description 2",
				OwnerID: "user-id",
				Status:  entities.Closed,
			},
		}, nil)
		encrypterMock.On("Decrypt", "Task Description 1").Return("Task Description 1", nil)
		encrypterMock.On("Decrypt", "Task Description 2").Return("Task Description 2", nil)

		output, err := TaskReadAllUseCase.Execute(usecases.TaskReadAllInput{
			User: entities.User{
				ID:   "user-id",
				Role: entities.UserRoleTechnician,
			},
		})

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Tasks, 2)

		assert.Equal(t, "task-id", output.Tasks[0].ID)
		assert.Equal(t, "Task Title", output.Tasks[0].Title)
		assert.Equal(t, "Task Description 1", output.Tasks[0].Summary)
		assert.Equal(t, "user-id", output.Tasks[0].OwnerID)
		assert.Equal(t, entities.Open, output.Tasks[0].Status)

		assert.Equal(t, "task-id-2", output.Tasks[1].ID)
		assert.Equal(t, "Task Title 2", output.Tasks[1].Title)
		assert.Equal(t, "Task Description 2", output.Tasks[1].Summary)
		assert.Equal(t, "user-id", output.Tasks[1].OwnerID)
		assert.Equal(t, entities.Closed, output.Tasks[1].Status)
	})

	t.Run("should return an error when an error occurs while finding tasks by user", func(t *testing.T) {
		taskRepositoryMock := mocks.NewTaskRepositoryInterface(t)
		encrypterMock := mocks.NewEncrypterInterface(t)
		TaskReadAllUseCase := usecases.NewTaskReadAllUseCase(taskRepositoryMock, encrypterMock)

		taskRepositoryMock.On("FindByUserID", "user-id").Return(nil, assert.AnError)

		output, err := TaskReadAllUseCase.Execute(usecases.TaskReadAllInput{
			User: entities.User{
				ID:   "user-id",
				Role: entities.UserRoleTechnician,
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, string(usecases.ErrorFindTasksByUser)+": assert.AnError general error for testing", err.Error())
		assert.Equal(t, usecases.TaskReadAllOutput{}, output)
	})
}
