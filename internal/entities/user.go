package entities

import (
	"time"
)

type RoleType string

const (
	SYSTEM_ADMIN RoleType = "SYSTEM ADMIN" // to be used for managing system only,not exposed
	DEV          RoleType = "DEV"          // to be used by developer, only read access in PROD

	SUPERADMIN RoleType = "SUPER ADMIN" // to be used for the owner of the application (customer)
	ADMIN      RoleType = "ADMIN"       // to be used by the application administrator

	STAFF      RoleType = "STAFF"
	OUTSOURCED RoleType = "OUTSOURCED"
	VIEWER     RoleType = "VIEWER"
)

type User struct {
	*Model              `mapstructure:",squash"`
	FirstName           string     `json:"firstName,omitempty"`
	LastName            string     `json:"lastName,omitempty"`
	Extension           string     `gorm:"not null" json:"extension,omitempty"`
	PhoneNumber         string     `gorm:"unique;not null" json:"phoneNumber,omitempty"`
	Email               string     `gorm:"unique;not null" json:"email,omitempty"`
	Password            string     `gorm:"not null" json:"-"`
	Role                RoleType   `gorm:"type:text;not null" json:"role,omitempty"`
	IsLoginDisabled     bool       `json:"isLoginDisabled"`
	IsLoggedIn          bool       `json:"isLoggedIn"`
	LastLoginTime       *time.Time `json:"lastLoginTime,omitempty"`
	LoginFailureCounter int16      `json:"loginFailureCounter,omitempty"`
	ResetPasswordString *string    `json:"resetPasswordString"`

	Experience string `json:"experience"`
	Department string `json:"department"`

	//References
	UserChannelDetails []UserChannelDetail `json:"userChannelDetails,omitempty"`
}

func (User) TableNameForQuery() string {
	return `"Users" E`
}
