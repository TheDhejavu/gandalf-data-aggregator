package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var maxIdleTime = 10 * time.Minute
var maxConnLifetime = time.Hour
var maxIdleConn = 10
var maxOpenConn = 100

func dbConn(dsn string) (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("open sql db failed: %w", err)
	}

	sqlDB.SetConnMaxIdleTime(maxIdleTime)
	sqlDB.SetConnMaxLifetime(maxConnLifetime)
	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetMaxOpenConns(maxOpenConn)

	return sqlDB, nil
}
func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	retries := 10

	dbConn, err := dbConn(dsn)
	if err != nil {
		panic(err)
	}

	cfg := postgres.Config{DSN: dsn, Conn: dbConn}
	for i := 1; i <= retries; i++ {
		time.Sleep(time.Second)
		db, err = gorm.Open(postgres.New(cfg), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			log.Warn().Err(err).Msgf("could not connect to postgres , retry %v / %v", i, retries)
			continue
		}
		return db, nil
	}
	return nil, fmt.Errorf("open postgres: %w", err)
}
