package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gandalf-data-aggregator/config"
	"gandalf-data-aggregator/service"
	workertask "gandalf-data-aggregator/worker/tasks"

	"github.com/go-redis/redis"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type JobHandler struct {
	service        *service.Service
	client         *asynq.Client
	cfg            config.Config
	wt             workertask.WorkerTask
	asynqInspector *asynq.Inspector
}

func NewJobHandler(
	service *service.Service,
	cfg config.Config,
	workerTask workertask.WorkerTask,
) JobHandler {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse redis URL")
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})

	return JobHandler{
		cfg:     cfg,
		client:  client,
		service: service,
		wt:      workerTask,
		asynqInspector: asynq.NewInspector(
			asynq.RedisClientOpt{
				Addr:     opt.Addr,
				Password: opt.Password,
				DB:       opt.DB,
			},
		),
	}
}

func (t JobHandler) GenerateActivityStats(ctx context.Context, task *asynq.Task) error {
	var payload workertask.QueuePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	err := t.service.GenerateActivityStats(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("unable to generate activitity stats: %w", err)
	}

	return nil
}

func (t JobHandler) ResolveUserActivityData(ctx context.Context, task *asynq.Task) error {
	var payload workertask.QueuePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", err)
	}

	err := t.service.FetchAndDumpUserActivities(ctx, payload.UserID, payload.DataKey)
	if err != nil {
		return fmt.Errorf("unable to fetch and dump user activities: %w", err)
	}

	_ = t.wt.EnqueueGenerateActivityStats(payload)

	return nil
}
