package entity

import (
	"testing"
	"time"
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

type Person struct {
	*Model `mapstructure:",squash"`

	Name   string `json:"name"`
	Gender string `gorm:"type:text" json:"gender"`
	Age    int    `json:"age"`
}

func TestTemplate(t *testing.T) {
	entities := []interface{}{
		Person{},
	}

	ExportTemplate(entities)
}
