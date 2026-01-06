package entities

type Customer struct {
	*Model `mapstructure:",squash"`

	Name           string `json:"name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	WhatsappNumber string `json:"whatsappNumber"`
	Address        string `json:"address"`

	//transient field
	Source string `json:"source" gorm:"-"`

	Enquiries    []Enquiry     `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE" json:"enquiries"`
	Measurements []Measurement `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE" json:"measurements"`
	Orders       []Order       `gorm:"foreignKey:CustomerId;constraint:OnDelete:CASCADE" json:"orders"`
}

func (Customer) TableName() string {
	return "stitch.Customers"
}

func (Customer) TableNameForQuery() string {
	return "\"stitch\".\"Customers\" E"
}
