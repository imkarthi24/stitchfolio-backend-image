package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
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

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
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

	*j = append((*j)[0:0], data...)
	return nil
}

type Measurement struct {
	*Model `mapstructure:",squash"`

	Values JSON `gorm:"type:jsonb" json:"values"`

	PersonId uint    `json:"personId"`
	Person   *Person `gorm:"foreignKey:PersonId" json:"person"`

	DressTypeId uint       `json:"dressTypeId"`
	DressType   *DressType `gorm:"foreignKey:DressTypeId" json:"dressType"`

	TakenById *uint `json:"takenById"`
	TakenBy   *User `gorm:"foreignKey:TakenById" json:"takenBy"`
}

func (Measurement) TableName() string {
	return "Measurements"
}

func (Measurement) TableNameForQuery() string {
	return "\"Measurements\" E"

}
