package queue

import (
	"context"

	"github.com/hibiken/asynq"
)

func NewPriorityCheckTask(addressID int64) *asynq.Task {
	payload := map[string]interface{}{"AddressID": addressID}
	return asynq.NewTask(TaskTypePriorityCheck, payload)
}

type IPriorityCheckExecutor interface {
	PriorityCheck(chatID int64) error
}

type PriorityCheckHandler struct {
	executor IPriorityCheckExecutor
}

func NewPriorityCheckHandler(executor IPriorityCheckExecutor) *PriorityCheckHandler {
	return &PriorityCheckHandler{executor}
}

func (h *PriorityCheckHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	addressID, err := task.Payload.GetInt("AddressID")
	if err != nil {
		return err
	}
	if err := h.executor.PriorityCheck(int64(addressID)); err != nil {
		return err
	}
	return nil
}
