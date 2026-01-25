package entities

type NotificationStatus string

const (
	NOTIF_PENDING   NotificationStatus = "PENDING"
	NOTIF_PARTIAL   NotificationStatus = "PARTIAL"
	NOTIF_COMPLETED NotificationStatus = "COMPLETED"
	NOTIF_FAULTED   NotificationStatus = "FAULTED"
)

type Notification struct {
	*Model
	Status       NotificationStatus `gorm:"default:'PENDING';type:text;not null" json:"status,omitempty"`
	SourceEntity string             `json:"sourceEntity,omitempty"`
	EntityId     uint               `json:"entityId,omitempty"`

	//Children
	EmailNotifications    []EmailNotification    `json:"-"`
	WhatsappNotifications []WhatsappNotification `json:"-"`
}

func (Notification) TableNameForQuery() string {
	return "\"stich\".\"Notifications\" E"
}

func (n *Notification) AddEmailNotification(email ...EmailNotification) {
	n.EmailNotifications = append(n.EmailNotifications, email...)
}

func (n *Notification) AddWhatsappNotification(whatsapp ...WhatsappNotification) {
	n.WhatsappNotifications = append(n.WhatsappNotifications, whatsapp...)
}
