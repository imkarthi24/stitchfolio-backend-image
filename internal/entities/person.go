package entities

type Person struct {
	*Model `mapstructure:",squash"`

	Name string `json:"name"`

	CustomerId uint      `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`

	Measurements []Measurement `gorm:"foreignKey:PersonId;constraint:OnDelete:CASCADE" json:"measurements"`
}

func (Person) TableName() string {
	return "Persons"
}

func (Person) TableNameForQuery() string {
	return "\"Persons\" E"
}
