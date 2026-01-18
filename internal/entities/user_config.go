package entities

type UserConfig struct {
	*Model `mapstructure:",squash"`

	Config string `json:"config,omitempty"`
	//References
	UserID uint  `json:"UserID,omitempty"`
	User   *User `json:"-,omitempty"`
}

func (UserConfig) TableName() string {
	return "stich.UserConfigs"
}

func (UserConfig) TableNameForQuery() string {
	return "\"stich\".\"UserConfigs\" E"
}
