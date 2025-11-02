package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskSayHello = "task:say_hello"
)

type Queue struct {
	client *asynq.Client
	server *asynq.Server
}

func NewQueue(redisAddr string) *Queue {
	r := asynq.RedisClientOpt{Addr: redisAddr}

	server := asynq.NewServer(
		r,
		asynq.Config{
			Concurrency: 1,
		},
	)

	client := asynq.NewClient(r)

	q := &Queue{
		client: client,
		server: server,
	}

	return q
}

func (q *Queue) Start() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSayHello, func(ctx context.Context, t *asynq.Task) error {
		fmt.Println("Hi Redis is running 🚀")
		return nil
	})

	go func() {
		if err := q.server.Run(mux); err != nil {
			log.Fatalf("Error while running: %v", err)
		}
	}()
}

func (q *Queue) EnqueueHello() error {
	task := asynq.NewTask(TaskSayHello, nil)
	_, err := q.client.Enqueue(task, asynq.MaxRetry(1))
	return err
}

func (q *Queue) StartScheduler() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			if err := q.EnqueueHello(); err != nil {
				log.Println("Error to add enqueue:", err)
			}
		}
	}()
}
