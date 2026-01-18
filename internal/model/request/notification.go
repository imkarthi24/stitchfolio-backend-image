package requestModel

import "github.com/imkarthi24/sf-backend/pkg/service/email"

type Notification struct {
	SourceEntity string `json:"sourceEntity,omitempty"`
	EntityId     uint   `json:"entityId,omitempty"`
}

type EmaiNotification struct {
	*Notification
	ToMailAddress string `json:"toMailAddress,omitempty"`
	Subject       string `json:"subject,omitempty"`
	Body          string `json:"body,omitempty"`
	EmailContent  *email.EmailContent
}
