package requestModel

type Enquiry struct {
	ID         uint   `json:"id,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
	Subject    string `json:"subject,omitempty"`
	Notes      string `json:"notes,omitempty"`
	Status     string `json:"status,omitempty"`
	CustomerId *uint  `json:"customerId,omitempty"`

	// Customer fields
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phoneNumber,omitempty"`
	WhatsappNumber string `json:"whatsappNumber,omitempty"`
	Address        string `json:"address,omitempty"`

	Source              string `json:"source,omitempty"`
	ReferredBy          string `json:"referredBy,omitempty"`
	ReferrerPhoneNumber string `json:"referrerPhoneNumber,omitempty"`
}
