package queue

import (
	"context"

	"github.com/hibiken/asynq"
)

func NewNotifyTask(chatID int64, message string) *asynq.Task {
	payload := map[string]interface{}{"ChatID": chatID, "Message": message}
	return asynq.NewTask(TaskTypeNotify, payload)
}

type INotifyExecutor interface {
	Notify(chatID int64, message string) error
}

type NotifyHandler struct {
	executor INotifyExecutor
}

func NewNotifyHandler(executor INotifyExecutor) *NotifyHandler {
	return &NotifyHandler{executor}
}

func (h *NotifyHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	chatID, err := task.Payload.GetInt("ChatID")
	if err != nil {
		return err
	}
	message, err := task.Payload.GetString("Message")
	if err != nil {
		return err
	}
	if err := h.executor.Notify(int64(chatID), message); err != nil {
		return err
	}
	return nil
}
