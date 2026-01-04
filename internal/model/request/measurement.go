package requestModel

import (
	"encoding/json"
)

type Measurement struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	MeasurementDate string          `json:"measurementDate,omitempty"`
	MeasurementBy   string          `json:"measurementBy,omitempty"`
	DressType       string          `json:"dressType,omitempty"`
	Measurements    json.RawMessage `json:"measurements,omitempty"`

	CustomerId *uint `json:"customerId,omitempty"`
}
