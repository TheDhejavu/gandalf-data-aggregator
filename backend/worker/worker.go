package main

import (
	"gandalf-data-aggregator/config"
	"gandalf-data-aggregator/models"
	"gandalf-data-aggregator/postgres"
	"gandalf-data-aggregator/repository"
	"gandalf-data-aggregator/service"
	"gandalf-data-aggregator/store"
	"gandalf-data-aggregator/webapi"
	"gandalf-data-aggregator/worker/handler"
	workertask "gandalf-data-aggregator/worker/tasks"
	"gandalf-data-aggregator/workerqueue"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msgf("Starting Scheduler worker....")
	var cfg config.Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("clean env failed to read env variables")
	}

	db, err := postgres.NewPostgresConnection(cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start postgres connection")
	}

	err = db.AutoMigrate(&models.User{}, &models.Activity{}, &models.Identifier{}, &models.DataKey{}, &models.ActivityStat{})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to auto migrate database")
	}

	workerTask := workertask.NewWorkerTask(cfg)
	service := service.NewService(cfg, repository.NewPostgres(db), webapi.NewGandalfClient(cfg), store.NewSessionStore(), workerTask, nil)

	srv := workerqueue.NewAsyncqServer(cfg)

	jobHandler := handler.NewJobHandler(service, cfg, workerTask)
	srv.AsynqSrvMux.HandleFunc(workertask.TypeActivityDataResolver, jobHandler.ResolveUserActivityData)
	srv.AsynqSrvMux.HandleFunc(workertask.TypeGenerateActivityStats, jobHandler.GenerateActivityStats)

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msgf("could not run server: %v", err)
	}
}
