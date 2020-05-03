package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/my1562/queue"
)

type NotifierImpl struct {
}

func (n *NotifierImpl) Notify(chatID int64, message string) error {
	log.Printf("Sending to %d: %s", chatID, message)
	return nil
}

const redisAddr = "127.0.0.1:6379"

func runServer() {
	notifier := &NotifierImpl{}
	notifyHandler := queue.NewNotifyHandler(notifier)
	redis := asynq.RedisClientOpt{Addr: redisAddr}
	server := asynq.NewServer(redis, asynq.Config{
		Concurrency: 1,
	})
	mux := asynq.NewServeMux()
	mux.Handle(queue.TaskTypeNotify, notifyHandler)
	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func main() {

	redis := asynq.RedisClientOpt{Addr: redisAddr}
	client := asynq.NewClient(redis)

	go runServer()

	task := queue.NewNotifyTask(123, "hello world")
	if err := client.Enqueue(task); err != nil {
		log.Fatal(err)
	}

	select {}

}
