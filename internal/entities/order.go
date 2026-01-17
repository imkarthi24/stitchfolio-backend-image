package entities

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	*Model `mapstructure:",squash"`

	Status OrderStatus `json:"status"`

	Notes string `json:"notes"`

	ExpectedDeliveryDate *time.Time `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time `json:"deliveredDate,omitempty"`

	CustomerId *uint     `json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"customer"`

	OrderTakenById *uint `json:"orderTakenById"`
	OrderTakenBy   *User `gorm:"foreignKey:OrderTakenById" json:"orderTakenBy"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderId" json:"orderItems"`

	// Calculated fields (populated via SQL subqueries, not stored in DB)
	OrderQuantity int     `gorm:"-" json:"-"`
	OrderValue    float64 `gorm:"-" json:"-"`
}

func (Order) TableName() string {
	return "Orders"
}

func (Order) TableNameForQuery() string {
	return "\"Orders\" E"
}
