package entities

import "time"

type ResponseStatus string

const (
	RESPONDED    ResponseStatus = "RESPONDED"
	NOTRESPONDED ResponseStatus = "NOT RESPONDED"
	REFUSED      ResponseStatus = "REFUSED"
	ACCEPTED     ResponseStatus = "ACCEPTED"
)

type EnquiryHistory struct {
	*Model `mapstructure:",squash"`

	Status *EnquiryStatus `gorm:"type:text" json:"status,omitempty"`

	EmployeeComment string         `json:"employeeComment,omitempty"`
	CustomerComment string         `json:"customerComment,omitempty"`
	VisitingDate    *time.Time     `json:"visitingDate,omitempty"`
	CallBackDate    *time.Time     `json:"callBackDate,omitempty"`
	EnquiryDate     *time.Time     `json:"enquiryDate,omitempty"`
	ResponseStatus  ResponseStatus `gorm:"type:text" json:"responseStatus,omitempty"`

	EnquiryId  uint  `json:"enquiryId,omitempty"`
	EmployeeId uint  `json:"employeeId,omitempty"`
	Employee   *User `json:"employee,omitempty"`

	// History tracking fields
	PerformedAt   time.Time `gorm:"not null" json:"performedAt"`
	PerformedById uint      `json:"performedById"`
	PerformedBy   *User     `gorm:"foreignKey:PerformedById" json:"-"`
}

func (EnquiryHistory) TableNameForQuery() string {
	return "\"stich\".\"EnquiryHistories\" E"
}
