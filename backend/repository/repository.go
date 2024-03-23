package repository

import (
	"context"
	"fmt"
	"gandalf-data-aggregator/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Postgres struct {
	Db *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{
		Db: db,
	}
}

func (pg *Postgres) Transaction(
	ctx context.Context,
	callback func(context.Context, *Postgres) error,
) error {
	session := pg.Db.Session(&gorm.Session{})
	err := session.Transaction(func(tx *gorm.DB) error {
		txpg := NewPostgres(tx)
		return callback(ctx, txpg)
	})
	if err != nil {
		return fmt.Errorf("in transaction: %w", err)
	}

	return nil
}

func (s *Postgres) FindOrCreateDataKey(ctx context.Context, dataKey *models.DataKey) (*models.DataKey, error) {
	tx := s.Db.Model(models.DataKey{}).
		Where("user_id = ? AND data_type = ?", dataKey.UserID, dataKey.DataType).
		Assign(&dataKey).
		FirstOrCreate(&dataKey)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return dataKey, nil
}

func (s *Postgres) FindOrCreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	tx := s.Db.Model(models.User{}).Where("external_id = ? ", user.ExternalID).Assign(&user).FirstOrCreate(&user)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (s *Postgres) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user *models.User
	tx := s.Db.Model(models.User{}).Where("id = ?", userID).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (s *Postgres) CreateActivity(ctx context.Context, userID uuid.UUID, activity *models.Activity) (*models.Activity, error) {
	tx := s.Db.Model(models.Activity{}).Where("subject = ? and user_id = ?", activity.Subject, userID).FirstOrCreate(&activity)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return activity, nil
}

func (s *Postgres) CreateActivities(ctx context.Context, activities []*models.Activity) ([]*models.Activity, error) {
	tx := s.Db.Clauses(clause.OnConflict{DoNothing: true}).Model(models.Activity{}).Create(&activities)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return activities, nil
}

func (s *Postgres) GetActivitySetByUser(ctx context.Context, userID uuid.UUID, limit int, page int) (*models.ActivityDataSet, error) {
	currentPage := page - 1
	if currentPage < 0 {
		currentPage = 0
	}

	activitySet := &models.ActivityDataSet{}
	tx := s.Db.Debug().Model(models.Activity{}).
		Distinct().
		Where("user_id = ?", userID).
		Preload("Subject").
		Order("date DESC").
		Count(&activitySet.Total).
		Limit(limit).
		Offset(limit * currentPage).
		Find(&activitySet.Data)

	if tx.Error != nil {
		return nil, tx.Error
	}

	activitySet.Limit = limit
	activitySet.Page = page

	return activitySet, nil
}

func (s *Postgres) BatchUpsertActivityStat(ctx context.Context, stats []*models.ActivityStat) error {
	return s.Db.Debug().
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "year"},
				{Name: "month"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"total": gorm.Expr("EXCLUDED.total + activity_stats.total"),
			}),
		}).
		Model(&models.ActivityStat{}).
		Create(&stats).Error
}

func (s *Postgres) SetActivityStatsToProcessedByUser(ctx context.Context, activityIDs uuid.UUIDs) error {
	tx := s.Db.Model(&models.Activity{}).
		Where("id IN ?", activityIDs).
		Update("processed", true)

	if tx.Error != nil {
		return tx.Error
	}

	fmt.Printf("<< Updated %v records.\n >> ", tx.RowsAffected)
	return nil
}

func (s *Postgres) GetActivityStatsByUser(ctx context.Context, userID uuid.UUID) ([]models.ActivityStat, error) {
	var stats []models.ActivityStat

	tx := s.Db.Model(&models.ActivityStat{}).
		Where("user_id = ?", userID).
		Find(&stats)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return stats, nil
}

func (s *Postgres) FetchUnprocessedUserActivities(ctx context.Context, userID uuid.UUID, limit int, page int) (*models.ActivityDataSet, error) {
	currentPage := page - 1
	if currentPage < 0 {
		currentPage = 0
	}

	activitySet := &models.ActivityDataSet{}
	tx := s.Db.Debug().Model(models.Activity{}).
		Distinct().
		Where("user_id = ? AND processed = ?", userID, false).
		Preload("Subject").
		Order("date DESC").
		Count(&activitySet.Total).
		Limit(limit).
		Offset(limit * currentPage).
		Find(&activitySet.Data)

	if tx.Error != nil {
		return nil, tx.Error
	}

	activitySet.Limit = limit
	activitySet.Page = page

	return activitySet, nil
}
