package entities

type UserChannelDetail struct {
	*Model `mapstructure:",squash"`

	//References
	UserID uint  `json:"UserID,omitempty"`
	User   *User `json:"-,omitempty"`

	UserChannelID uint     `json:"userChannelId,omitempty"` // Refers to channelId
	UserChannel   *Channel `json:"channel,omitempty"`
}

func (UserChannelDetail) TableNameForQuery() string {
	return "\"UserChannelDetails\" E"
}
