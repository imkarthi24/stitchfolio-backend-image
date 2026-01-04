package responseModel

type Customer struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Name           string `json:"name,omitempty"`
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
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}
