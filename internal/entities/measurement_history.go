package entities

import "time"

type MeasurementHistoryAction string

const (
	MeasurementHistoryActionCreated MeasurementHistoryAction = "CREATED"
	MeasurementHistoryActionUpdated MeasurementHistoryAction = "UPDATED"
	MeasurementHistoryActionDeleted MeasurementHistoryAction = "DELETED"
)

type MeasurementHistory struct {
	*Model `mapstructure:",squash"`

	Action MeasurementHistoryAction `gorm:"type:string;not null" json:"action"`

	// Comma-separated list of changed fields (e.g., "measurements,measurementDate")
	ChangedValues string `json:"changedValues,omitempty"`

	OldValues JSON `gorm:"type:jsonb" json:"oldValues,omitempty"`

	MeasurementId uint         `json:"measurementId"`
	Measurement   *Measurement `gorm:"foreignKey:MeasurementId" json:"measurement"`

	PerformedAt   time.Time `gorm:"not null" json:"performedAt"`
	PerformedById uint      `json:"performedById"`
	PerformedBy   *User     `gorm:"foreignKey:PerformedById" json:"-"`
}

func (MeasurementHistory) TableName() string {
	return "MeasurementHistories"
}

func (MeasurementHistory) TableNameForQuery() string {
	return "\"MeasurementHistories\" E"
}
