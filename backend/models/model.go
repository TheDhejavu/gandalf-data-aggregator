package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:UUID;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" sql:"index"`
}

type User struct {
	Base
	Username   string
	ExternalID string
	Email      string
	AvatarURL  string
	FirstName  string
	LastName   string
}

type DataType string

var (
	DataTypeNetflix DataType = "netflix"
)

type DataKey struct {
	Base
	UserID   uuid.UUID `gorm:"type:UUID"`
	DataType DataType
	Key      string
}

type Activity struct {
	Base
	UserID             uuid.UUID    `gorm:"type:UUID" json:"user_id"`
	ProviderActivityID string       `json:"provider_activity_id" gorm:"uniqueIndex"`
	Subject            []Identifier `json:"subject"`
	Title              string       `json:"title"`
	Date               time.Time    `json:"date,omitempty"`
	Processed          bool         `gorm:"default:false" json:"processed"`
}

type ActivityDataSet struct {
	Limit int         `json:"limit"`
	Page  int         `json:"page"`
	Total int64       `json:"total"`
	Data  []*Activity `json:"data"`
}

type Identifier struct {
	Base
	ActivityID     uuid.UUID `gorm:"type:UUID" json:"activity_id"`
	Value          string    `json:"value"`
	IdentifierType string    `json:"identifier_type"`
}

type ActivityStat struct {
	Base
	UserID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Year   int       `gorm:"primary_key;"`
	Month  int       `gorm:"primary_key;"`
	Total  int
}

type YearData struct {
	Labels []int
	Months []string
}

type YearDataStat struct {
	YearData    map[int]YearData `json:"year_data"`
	CurrentYear string           `json:"current_year"`
}
