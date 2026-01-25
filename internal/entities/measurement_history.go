package entities

import (
	"time"

	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
)

type MeasurementHistoryAction string

const (
	MeasurementHistoryActionCreated MeasurementHistoryAction = "CREATED"
	MeasurementHistoryActionUpdated MeasurementHistoryAction = "UPDATED"
	MeasurementHistoryActionDeleted MeasurementHistoryAction = "DELETED"
)

type MeasurementHistory struct {
	*Model `mapstructure:",squash"`

	Action MeasurementHistoryAction `gorm:"type:text;not null" json:"action"`
 
	OldValues entitiy_types.JSON `gorm:"type:jsonb" json:"oldValues,omitempty"`

	MeasurementId uint         `json:"measurementId"`
	Measurement   *Measurement `gorm:"foreignKey:MeasurementId" json:"measurement"`

	PerformedAt   time.Time `gorm:"not null" json:"performedAt"`
	PerformedById uint      `json:"performedById"`
	PerformedBy   *User     `gorm:"foreignKey:PerformedById" json:"-"`
}

func (MeasurementHistory) TableNameForQuery() string {
	return "\"stich\".\"MeasurementHistories\" E"
}
