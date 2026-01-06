package requestModel

type User struct {
	ID                  uint                `json:"id,omitempty"`
	IsActive            bool                `json:"isActive,omitempty"`
	FirstName           string              `json:"firstName,omitempty"`
	LastName            string              `json:"lastName,omitempty"`
	Extension           string              `json:"extension,omitempty"`
	PhoneNumber         string              `json:"phoneNumber,omitempty"`
	Email               string              `json:"email,omitempty"`
	Password            string              `json:"password,omitempty"`
	Role                string              `json:"role,omitempty"`
	IsLoginDisabled     bool                `json:"isLoginDisabled,omitempty"`
	IsLoggedIn          bool                `json:"isLoggedIn,omitempty"`
	LastLoginTime       string              `json:"lastLoginTime,omitempty"`
	LoginFailureCounter int                 `json:"loginFailureCounter,omitempty"`
	ResetPasswordString *string             `json:"resetPasswordString,omitempty"`
	Experience          string              `json:"experience,omitempty"`
	Department          string              `json:"department,omitempty"`
	UserChannelDetails  []UserChannelDetail `json:"userChannelDetails,omitempty"`
}

type UserConfig struct {
	ID       uint   `json:"id,omitempty"`
	Config   string `json:"config,omitempty"`
	UserID   uint   `json:"userID,omitempty"`
	IsActive bool   `json:"isActive"`
}

type UserChannelDetail struct {
	ID        uint `json:"id,omitempty"`
	UserID    uint `json:"userID,omitempty"`
	IsActive  bool `json:"isActive"`
	ChannelId uint `json:"channelID,omitempty"`
}
