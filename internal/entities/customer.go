package entities

type Customer struct {
	*Model `mapstructure:",squash"`

	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	WhatsappNumber string `json:"whatsappNumber"`
	Address        string `json:"address"`

	//transient field
	Source string `json:"source" gorm:"-"`

	Persons   []Person  `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE" json:"persons"`
	Enquiries []Enquiry `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE" json:"enquiries"`
	Orders    []Order   `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE" json:"orders"`
}

func (Customer) TableName() string {
	return "Customers"
}

func (Customer) TableNameForQuery() string {
	return "\"Customers\" E"
}
