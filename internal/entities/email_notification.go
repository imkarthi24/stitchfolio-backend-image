package entities

type EmailNotification struct {
	*Model
	Status        string `json:"firstName,omitempty"`
	ToMailAddress string `json:"toMailAddress,omitempty"`
	Subject       string `json:"subject,omitempty"`
	Body          string `json:"body,omitempty"`

	//References
	NotificationId uint          `json:"notificationId,omitempty"`
	Notification   *Notification `json:"-"`
}

func (EmailNotification) TableNameForQuery() string {
	return "\"stich\".\"EmailNotifications\" E"
}
