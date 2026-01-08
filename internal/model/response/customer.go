package responseModel

type Customer struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phoneNumber,omitempty"`
	WhatsappNumber string `json:"whatsappNumber,omitempty"`
	Address        string `json:"address,omitempty"`

	Enquiries    []Enquiry     `json:"enquiries,omitempty"`
	Measurements []Measurement `json:"measurements,omitempty"`
	Orders       []Order       `json:"orders,omitempty"`
}

type CustomerAutoComplete struct {
	CustomerID  uint   `json:"customerId,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}
