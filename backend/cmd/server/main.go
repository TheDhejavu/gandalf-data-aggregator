package main

import (
	"fmt"
	"gandalf-data-aggregator/config"
	"gandalf-data-aggregator/delivery"
	"gandalf-data-aggregator/models"
	token "gandalf-data-aggregator/pkg/jwt"
	"gandalf-data-aggregator/postgres"
	"gandalf-data-aggregator/repository"
	"gandalf-data-aggregator/service"
	"gandalf-data-aggregator/store"
	"gandalf-data-aggregator/webapi"
	workertask "gandalf-data-aggregator/worker/tasks"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
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

	jwtMaker, err := token.NewJWTMaker(cfg.JWTSecretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize jwt maker")
	}
	workerTask := workertask.NewWorkerTask(cfg)

	service := service.NewService(cfg, repository.NewPostgres(db), webapi.NewGandalfClient(cfg), store.NewSessionStore(), workerTask, jwtMaker)

	router := echo.New()

	srv := delivery.NewServer(router, cfg, service, jwtMaker)
	log.Info().Msgf("Starting server on Port %s", cfg.Port)

	if err := srv.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatal().Err(err).Msg("Server failed to run")
	}
}
