package entities

import (
	"time"

	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
)

type OrderHistoryAction string

const (
	OrderHistoryActionCreated OrderHistoryAction = "CREATED"
	OrderHistoryActionUpdated OrderHistoryAction = "UPDATED"
	OrderHistoryActionDeleted OrderHistoryAction = "DELETED"
)

// Order change field constants
const (
	OrderChangeFieldStatus               string = "status"
	OrderChangeFieldExpectedDeliveryDate string = "expectedDeliveryDate"
	OrderChangeFieldDeliveredDate        string = "deliveredDate"
)

type OrderHistory struct {
	*Model `mapstructure:",squash"`

	Action OrderHistoryAction `gorm:"type:text;not null" json:"action"`

	// Comma-separated list of changed fields (e.g., "status,expectedDeliveryDate")
	ChangedFields string `json:"changedFields,omitempty"`

	Status               *OrderStatus `gorm:"type:text" json:"status,omitempty"`
	ExpectedDeliveryDate *time.Time   `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time   `json:"deliveredDate,omitempty"`

	OrderItemId   *uint               `json:"orderItemId,omitempty"`
	OrderItemData *entitiy_types.JSON `gorm:"type:jsonb" json:"orderItemData,omitempty"`

	OrderId uint   `json:"orderId"`
	Order   *Order `gorm:"foreignKey:OrderId" json:"order"`

	PerformedAt   time.Time `gorm:"not null" json:"performedAt"`
	PerformedById uint      `json:"performedById"`
	PerformedBy   *User     `gorm:"foreignKey:PerformedById" json:"-"`
}

func (OrderHistory) TableNameForQuery() string {
	return "\"stich\".\"OrderHistories\" E"
}
