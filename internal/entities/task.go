package entities

import "time"

type Task struct {
	*Model `mapstructure:",squash"`

	Title        string     `gorm:"size:255;not null" json:"title"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	IsCompleted  bool       `gorm:"default:false" json:"isCompleted"`
	Priority     *int       `json:"priority,omitempty"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
	ReminderDate *time.Time `json:"reminderDate,omitempty"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`

	AssignedToId *uint `json:"assignedToId,omitempty"`
	AssignedTo   *User `gorm:"foreignKey:AssignedToId" json:"assignedTo,omitempty"`
}

func (Task) TableNameForQuery() string {
	return "\"stich\".\"Tasks\" E"
}
