package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoStep(userID, todoID uuid.UUID, name string) (entity.TodoStep, error) {
	_, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.TodoStep{}, fmt.Errorf("todo not found: %v", todoID)
	}

	step := entity.TodoStep{
		ID:     uuid.New(),
		Name:   name,
		TodoID: todoID,
	}
	if _, err = dal.InsertTodoStep(step); err != nil {
		return entity.TodoStep{}, fmt.Errorf("fails to create todo step")
	}

	return step, nil
}

func UpdateTodoStep(userID uuid.UUID, step entity.TodoStep) (entity.TodoStep, error) {
	oldStep, err := OwnTodoStep(userID, step.ID)
	if err != nil {
		return entity.TodoStep{}, err
	}

	if step.TodoID != oldStep.TodoID {
		return entity.TodoStep{}, fmt.Errorf("unable to update todo id")
	}

	if step.Done && !oldStep.Done {
		step.DoneAt = time.Now()
	}

	if _, err = dal.UpdateTodoStep(step); err != nil {
		return entity.TodoStep{}, fmt.Errorf("fails to update todo step")
	}

	return step, nil
}

func DeleteTodoStep(userID, todoStepID uuid.UUID) (entity.TodoStep, error) {
	if _, err := OwnTodoStep(userID, todoStepID); err != nil {
		return entity.TodoStep{}, err
	}

	return dal.DeleteTodoStep(todoStepID)
}

func OwnTodoStep(userID, todoStepID uuid.UUID) (entity.TodoStep, error) {
	step, err := dal.GetTodoStep(todoStepID)
	if err != nil {
		return entity.TodoStep{}, fmt.Errorf("fails to get todo step: %v", todoStepID)
	}

	_, err = OwnTodo(userID, step.TodoID)
	if err != nil {
		return entity.TodoStep{}, utils.NewErrorWithForbidden("unable to handle non-owned todo: %v", step.ID)
	}

	return step, nil
}
