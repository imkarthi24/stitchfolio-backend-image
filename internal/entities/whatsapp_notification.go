package entities

type WhatsappNotification struct {
	*Model
	Status              string `json:"firstName,omitempty"`
	ReceipientNumber    string `json:"receipientNumber,omitempty"`
	ReceipientExtension string `json:"receipientExtension,omitempty"`
	Subject             string `json:"subject,omitempty"`
	Body                string `json:"body,omitempty"`

	//References
	NotificationID uint          `json:"notificationID,omitempty"`
	Notification   *Notification `json:"-"`
}

func (WhatsappNotification) TableNameForQuery() string {
	return "\"stich\".\"WhatsappNotifications\" E"
}
