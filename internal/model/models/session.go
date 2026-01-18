package models

import "github.com/imkarthi24/sf-backend/internal/entities"

type Session struct {
	Email                 string            `json:"email,omitempty"`
	Role                  entities.RoleType `json:"role,omitempty"`
	FirstName             string            `json:"firstName,omitempty"`
	LastName              string            `json:"lastName,omitempty"`
	UserId                *uint             `json:"userId,omitempty"`
	ChannelId             uint              `json:"channelId,omitempty"`
	ChannelName           string            `json:"channelName,omitempty"`
	AccessibleLocationIds []uint            `json:"accessibleLocationIds,omitempty"`
	IsSystemSession       bool              `json:"-,omitempty"`
}
