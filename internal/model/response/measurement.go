package responseModel

import (
	"encoding/json"
	"time"
)

type Measurement struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	MeasurementDate string          `json:"measurementDate,omitempty"`
	MeasurementBy   string          `json:"measurementBy,omitempty"`
	DressType       string          `json:"dressType,omitempty"`
	Measurements    json.RawMessage `json:"measurements,omitempty"`

	CustomerId   *uint     `json:"customerId,omitempty"`
	Customer     *Customer `json:"customer,omitempty"`
	CustomerName string    `json:"customerName,omitempty"`

	MeasurementTakenById *uint `json:"measurementTakenById,omitempty"`
	// MeasurementTakenBy   *User  `json:"measurementTakenBy,omitempty"`
	MeasurementTakenBy string `json:"measurementTakenBy,omitempty"`

	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	UpdatedById *uint      `json:"updatedById,omitempty"`
}
