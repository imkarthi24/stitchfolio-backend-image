package entities

import (
	"time"

	"github.com/imkarthi24/sf-backend/pkg/constants"
	"gorm.io/gorm"
)

type Model struct {
	ID          uint       `gorm:"not null;primarykey" json:"ID,omitempty"`
	CreatedAt   *time.Time `gorm:"<-:create" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	IsActive    bool       `gorm:"default:true;type:bool" json:"isActive"`
	CreatedById *uint      `gorm:"<-:create" json:"createdById,omitempty"`
	UpdatedById *uint      `json:"updatedById,omitempty"`

	//Channel Id must be create only since it will interfere with update operations
	//Use tx.Exec as raw query to update the channelId as done in channel after-create hook
	ChannelId uint `gorm:"<-:create" json:"channelId,omitempty"`
}

func (u *Model) BeforeUpdate(tx *gorm.DB) (err error) {

	//Rare panic scenario
	if u == nil {
		return
	}

	if id, ok := tx.Get(constants.USER_ID); ok {
		u.UpdatedById = id.(*uint)
	}
	return
}

func (u *Model) BeforeCreate(tx *gorm.DB) (err error) {

	//Rare panic scenario
	if u == nil {
		return
	}

	if id, ok := tx.Get(constants.USER_ID); ok {
		u.CreatedById = id.(*uint)
	}

	if id, ok := tx.Get(constants.CHANNEL_ID); ok {
		u.ChannelId = id.(uint)
	}

	return
}
