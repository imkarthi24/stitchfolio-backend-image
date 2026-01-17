package entities

import entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"

type Measurement struct {
	*Model `mapstructure:",squash"`

	Value entitiy_types.JSON `gorm:"type:jsonb" json:"values"`

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
