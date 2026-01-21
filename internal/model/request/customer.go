package requestModel

type Customer struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phoneNumber,omitempty"`
	WhatsappNumber string `json:"whatsappNumber,omitempty"`
	Address        string `json:"address,omitempty"`
	Age            int    `json:"age,omitempty"`
	Gender         string `json:"gender,omitempty"`
}
