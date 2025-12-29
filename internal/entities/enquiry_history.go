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

	EmployeeComment string         `json:"employeeComment,omitempty"`
	CustomerComment string         `json:"customerComment,omitempty"`
	VisitingDate    *time.Time     `json:"visitingDate,omitempty"`
	CallBackDate    *time.Time     `json:"callBackDate,omitempty"`
	EnquiryDate     time.Time      `gorm:"not null" json:"enquiryDate,omitempty"`
	ResponseStatus  ResponseStatus `gorm:"type:string;not null" json:"responseStatus,omitempty"`

	EnquiryId  uint  `json:"enquiryId,omitempty"`
	EmployeeId uint  `json:"employeeId,omitempty"`
	Employee   *User `json:"employee,omitempty"`
}

func (EnquiryHistory) TableName() string {
	return "stitch.EnquiryHistories"
}

func (EnquiryHistory) TableNameForQuery() string {
	return "\"EnquiryHistories\" E"
}
