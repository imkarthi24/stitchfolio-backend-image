package requestModel

type UpdateEnquiryAndCustomer struct {
	// Enquiry fields
	ID         uint   `json:"id,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
	Subject    string `json:"subject,omitempty"`
	Notes      string `json:"notes,omitempty"`
	Status     string `json:"status,omitempty"`
	CustomerId *uint  `json:"customerId,omitempty"`

	Source              string `json:"source,omitempty"`
	ReferredBy          string `json:"referredBy,omitempty"`
	ReferrerPhoneNumber string `json:"referrerPhoneNumber,omitempty"`

	// Customer fields
	CustomerIsActive bool   `json:"customerIsActive,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Email            string `json:"email,omitempty"`
	PhoneNumber      string `json:"phoneNumber,omitempty"`
	WhatsappNumber   string `json:"whatsappNumber,omitempty"`
	Address          string `json:"address,omitempty"`
}
