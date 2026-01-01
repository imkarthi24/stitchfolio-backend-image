package entities

import "gorm.io/gorm"

type ChannelStatus string

const (
	CHANNEL_ACTIVE   NotificationStatus = "ACTIVE"
	CHANNEL_INACTIVE NotificationStatus = "INACTIVE"
)

type Channel struct {
	*Model `mapstructure:",squash"`
	Name   string        `json:"name,omitempty"`
	Status ChannelStatus `gorm:"default:'ACTIVE';type:string;not null" json:"status,omitempty"`

	//Reference
	OwnerUserID uint  `json:"ownerUserId,omitempty"`
	OwnerUser   *User `gorm:"foreignKey:OwnerUserID;references:ID" json:"-"`
}

func (Channel) TableName() string {
	return "stitch.Channels"
}

func (Channel) TableNameForQuery() string {
	return "\"Channels\" E"
}

func (c *Channel) AfterCreate(tx *gorm.DB) (err error) {

	res := tx.Exec("UPDATE \"Channels\" SET channel_id = ? WHERE id = ?", c.ID, c.ID)
	if res.Error != nil {
		return err
	}

	res = tx.Exec("UPDATE \"Users\" SET channel_id = ? WHERE id = ?", c.ID, c.OwnerUserID)
	if res.Error != nil {
		return err
	}

	return
}
