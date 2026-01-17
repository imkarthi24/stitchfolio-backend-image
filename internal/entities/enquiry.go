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
	Status  EnquiryStatus `gorm:"type:string;not null" json:"status"`
	Notes   string        `json:"notes"`

	Source              string `json:"source,omitempty"`
	ReferredBy          string `json:"referredBy,omitempty"`
	ReferrerPhoneNumber string `json:"referrerPhoneNumber,omitempty"`

	CustomerId *uint     `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`
}

func (Enquiry) TableName() string {
	return "Enquiries"
}

func (Enquiry) TableNameForQuery() string {
	return "\"Enquiries\" E"
}
