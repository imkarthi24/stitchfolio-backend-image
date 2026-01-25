package entities

type EnquiryStatus string

const (
	EnquiryStatusNew      EnquiryStatus = "new"
	EnquiryStatusAccepted EnquiryStatus = "accepted"
	EnquiryStatusCallback EnquiryStatus = "callback"
	EnquiryStatusBacklog  EnquiryStatus = "closed"
)

type Enquiry struct {
	*Model `mapstructure:",squash"`

	Subject string        `json:"subject"`
	Status  EnquiryStatus `gorm:"type:text;not null" json:"status"`
	Notes   string        `json:"notes"`

	Source              string `json:"source,omitempty"`
	ReferredBy          string `json:"referredBy,omitempty"`
	ReferrerPhoneNumber string `json:"referrerPhoneNumber,omitempty"`

	CustomerId *uint     `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`
}

func (Enquiry) TableNameForQuery() string {
	return "\"stich\".\"Enquiries\" E"
}
