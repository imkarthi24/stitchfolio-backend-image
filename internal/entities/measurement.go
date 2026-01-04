package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

// validateMeasurementsStructure validates that the JSON bytes represent a valid measurements array
// Expected structure:
// [
//
//	{
//	  "type": "Pant",
//	  "values": {"Height": "28", "Hip": "30", ...}
//	}
//
// ]
func validateMeasurementsStructure(bytes []byte) error {
	var temp interface{}
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return errors.New("invalid JSON format")
	}

	// Allow null, empty array, or array of objects
	if temp != nil {
		arr, ok := temp.([]interface{})
		if !ok {
			return errors.New("measurements must be a JSON array")
		}
		for i, item := range arr {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("measurements array item at index %d must be an object", i)
			}
			// Check for 'type' field
			if _, exists := itemMap["type"]; !exists {
				return fmt.Errorf("measurements array item at index %d missing 'type' field", i)
			}
			// Check for 'values' field
			if _, exists := itemMap["values"]; !exists {
				return fmt.Errorf("measurements array item at index %d missing 'values' field", i)
			}
			// Validate 'values' is an object
			if _, ok := itemMap["values"].(map[string]interface{}); !ok {
				return fmt.Errorf("measurements array item at index %d 'values' must be an object", i)
			}
		}
	}

	return nil
}

// NewJSON creates a new JSON type with validation
func NewJSON(data []byte) (JSON, error) {
	if len(data) == 0 {
		return nil, nil
	}

	err := validateMeasurementsStructure(data)
	if err != nil {
		return nil, err
	}

	return JSON(data), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to unmarshal JSON value")
	}

	// Validate structure
	err := validateMeasurementsStructure(bytes)
	if err != nil {
		return err
	}

	result := json.RawMessage{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	*j = JSON(result)
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	// Validate structure before accepting the data
	err := validateMeasurementsStructure(data)
	if err != nil {
		return err
	}

	*j = append((*j)[0:0], data...)
	return nil
}

type Measurement struct {
	*Model `mapstructure:",squash"`

	MeasurementDate *time.Time `json:"measurementDate"`
	MeasurementBy   string     `json:"measurementBy"`
	DressType       string     `json:"dressType"`
	Measurements    JSON       `gorm:"type:jsonb" json:"measurements"`

	CustomerId *uint     `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`
}

func (Measurement) TableName() string {
	return "stitch.Measurements"
}

func (Measurement) TableNameForQuery() string {
	return "\"Measurements\" E"
}
