package entities

type UserConfig struct {
	*Model `mapstructure:",squash"`

	Config string `json:"config,omitempty"`
	//References
	UserID uint  `json:"UserID,omitempty"`
	User   *User `json:"-,omitempty"`
}

func (UserConfig) TableNameForQuery() string {
	return "\"stich\".\"UserConfigs\" E"
}
