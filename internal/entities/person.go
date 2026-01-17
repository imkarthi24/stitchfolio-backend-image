package entities

type Person struct {
	*Model `mapstructure:",squash"`

	Name   string `json:"name"`
	Gender string `gorm:"type:text" json:"gender"`
	Age    int    `json:"age"`

	CustomerId uint      `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`

	Measurements []Measurement `gorm:"foreignKey:PersonId;constraint:OnDelete:CASCADE" json:"measurements"`
}

func (Person) TableName() string {
	return "stitch.Persons"
}

func (Person) TableNameForQuery() string {
	return "\"stitch\".\"Persons\" E"
}
