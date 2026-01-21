package entities

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
	OTHER  Gender = "OTHER"
)

type Person struct {
	*Model `mapstructure:",squash"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    Gender `gorm:"type:text" json:"gender"`
	Age       *int   `json:"age"`

	CustomerId uint      `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`

	Measurements []Measurement `gorm:"foreignKey:PersonId;constraint:OnDelete:CASCADE" json:"measurements"`
}

func (Person) TableName() string {
	return "stich.Persons"
}

func (Person) TableNameForQuery() string {
	return "\"stich\".\"Persons\" E"
}
