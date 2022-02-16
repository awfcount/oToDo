package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoListFolder(userID string, todoListFolderName string) (entity.TodoListFolder, error) {
	folder := entity.TodoListFolder{
		Name:   todoListFolderName,
		UserID: userID,
	}
	if err := dal.InsertTodoListFolder(&folder); err != nil {
		return entity.TodoListFolder{}, err
	}

	return folder, nil
}

func GetTodoListFolder(userID, todoListFolderID string) (entity.TodoListFolder, error) {
	return OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID string) ([]entity.TodoListFolder, error) {
	vec, err := dal.SelectTodoListFolders(userID)
	return vec, fmt.Errorf("fails to get user: %w", err)
}

func DeleteTodoListFolder(userID, todoListFolderID string) (entity.TodoListFolder, error) {
	write := func(err error) (entity.TodoListFolder, error) {
		return entity.TodoListFolder{}, err
	}

	folder, err := OwnTodoListFolder(userID, todoListFolderID)
	if err != nil {
		return write(err)
	}

	// TODO[feat] Whether to cascade delete todo lists
	// Cascade delete todo lists
	if _, err = dal.DeleteTodoListsByFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to cascade delete todo lists: %w", err))
	}

	if err = dal.DeleteTodoListFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to delete todo list folder: %w", err))
	}

	return folder, nil
}

// Verify permission
func OwnTodoListFolder(userID, todoListFolderID string) (entity.TodoListFolder, error) {
	todoListFolder, err := dal.SelectTodoListFolder(todoListFolderID)
	if err != nil {
		return entity.TodoListFolder{}, fmt.Errorf("fails to get todo list folder: %v", todoListFolderID)
	}

	if todoListFolder.UserID != userID {
		return entity.TodoListFolder{}, utils.NewErrorWithForbidden("unable to handle non-owned todo list folder: %v", todoListFolderID)
	}

	return todoListFolder, nil
}
