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

type PriorityCheckerImpl struct {
}

func (p *PriorityCheckerImpl) PriorityCheck(arrressID int64) error {
	log.Printf("Checking address: %d", arrressID)
	return nil
}

const redisAddr = "127.0.0.1:6379"

func runServer() {
	notifier := &NotifierImpl{}
	priorityChecker := &PriorityCheckerImpl{}

	redis := asynq.RedisClientOpt{Addr: redisAddr}
	server := asynq.NewServer(redis, asynq.Config{
		Concurrency: 1,
	})

	mux := asynq.NewServeMux()
	mux.Handle(
		queue.TaskTypeNotify,
		queue.NewNotifyHandler(notifier),
	)
	mux.Handle(
		queue.TaskTypePriorityCheck,
		queue.NewPriorityCheckHandler(priorityChecker),
	)

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func main() {

	redis := asynq.RedisClientOpt{Addr: redisAddr}
	client := asynq.NewClient(redis)

	go runServer()

	if err := client.Enqueue(queue.NewNotifyTask(123, "hello world")); err != nil {
		log.Fatal(err)
	}
	if err := client.Enqueue(queue.NewPriorityCheckTask(100500)); err != nil {
		log.Fatal(err)
	}

	select {}

}
