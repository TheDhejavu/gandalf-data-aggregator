package workertask

import (
	"encoding/json"
	"errors"
	"fmt"
	"gandalf-data-aggregator/config"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TypeActivityDataResolver  = "resolver:data"
	TypeGenerateActivityStats = "generate:stats"
)

type WorkerTask struct {
	client *asynq.Client
	cfg    config.Config
}

func NewWorkerTask(cfg config.Config) WorkerTask {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse redis URL")
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})

	return WorkerTask{
		cfg:    cfg,
		client: client,
	}
}

type QueuePayload struct {
	UserID  uuid.UUID
	DataKey string
}

func (t WorkerTask) EnqueueActivityDataResolver(queuePayload QueuePayload) error {
	payload, err := json.Marshal(queuePayload)
	if err != nil {
		return fmt.Errorf("unable to marsahl page %w", err)
	}

	newTask := asynq.NewTask(
		TypeActivityDataResolver,
		payload,
		asynq.MaxRetry(2),
		asynq.Retention(1*time.Hour),
	)

	_, err = t.client.Enqueue(newTask)
	if err != nil && errors.Is(err, asynq.ErrDuplicateTask) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}
	return nil
}

func (t WorkerTask) EnqueueGenerateActivityStats(queuePayload QueuePayload) error {
	payload, err := json.Marshal(queuePayload)
	if err != nil {
		return fmt.Errorf("unable to marsahl page %w", err)
	}

	newTask := asynq.NewTask(
		TypeGenerateActivityStats,
		payload,
		asynq.Unique(1*time.Hour),
		asynq.MaxRetry(5),
		asynq.Retention(1*time.Hour),
	)

	_, err = t.client.Enqueue(newTask)
	if err != nil && errors.Is(err, asynq.ErrDuplicateTask) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}

	return nil
}
