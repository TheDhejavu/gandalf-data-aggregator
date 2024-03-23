package workerqueue

import (
	"fmt"
	"gandalf-data-aggregator/config"
	"net/http"

	"strings"

	"github.com/rs/zerolog/log"

	redis "github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

type AsynqServer struct {
	AsynqSrv       *asynq.Server
	AsynqSrvMux    *asynq.ServeMux
	AsynqScheduler *asynq.Scheduler
	AsynqInspector *asynq.Inspector
}

func getRedisConnOpt(opt *redis.Options) asynq.RedisConnOpt {
	return asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	}
}

func NewAsyncqServer(cfg config.Config) *AsynqServer {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse redis URL")
	}

	redisOpts := getRedisConnOpt(opt)
	asyncSrv := asynq.NewServer(
		redisOpts,
		asynq.Config{
			Concurrency: 30,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	asynqSrvMux := asynq.NewServeMux()
	scheduler := asynq.NewScheduler(getRedisConnOpt(opt), nil)

	return &AsynqServer{
		AsynqSrv:       asyncSrv,
		AsynqSrvMux:    asynqSrvMux,
		AsynqScheduler: scheduler,
		AsynqInspector: asynq.NewInspector(redisOpts),
	}
}

func (srv AsynqServer) Run() error {
	if err := srv.AsynqScheduler.Start(); err != nil {
		return fmt.Errorf("unable to start asynq scheduler %w", err)
	}
	return srv.AsynqSrv.Run(srv.AsynqSrvMux)
}

func (srv AsynqServer) DeleteArchivedTasks(queueName string) error {
	log.Info().Msgf("Delete archived tasks for %s", queueName)

	n, err := srv.AsynqInspector.DeleteAllArchivedTasks(queueName)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}
		return err
	}

	log.Info().Msgf("Deleted == %s(%d)", queueName, n)
	return nil
}

func RegisterMonitoringHandler(mux *http.ServeMux, cfg config.Config) {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse redis URL")
	}

	// frontend monitoring
	asynqmonHandler := asynqmon.New(asynqmon.Options{
		RootPath: "/monitoring",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     opt.Addr,
			Password: opt.Password,
			DB:       opt.DB,
		},
	})
	mux.Handle(fmt.Sprintf("%s/", asynqmonHandler.RootPath()), asynqmonHandler)
}
