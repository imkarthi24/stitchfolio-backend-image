package entities

import "time"

type OrderStatus string

const (
	DRAFT               OrderStatus = "DRAFT"
	CONFIRMED           OrderStatus = "CONFIRMED"
	DESIGN_CONFIRMED    OrderStatus = "DESIGN_CONFIRMED"
	RAW_MATERIAL_SOURCE OrderStatus = "RAW_MATERIAL_SOURCE"
	CUTTING             OrderStatus = "CUTTING"
	STITCHING           OrderStatus = "STITCHING"
	FINISHING           OrderStatus = "FINISHING"
	READY_FOR_DELIVERY  OrderStatus = "READY_FOR_DELIVERY"
	DELIVERED           OrderStatus = "DELIVERED"
	CANCELLED           OrderStatus = "CANCELLED"
)

type Order struct {
	*Model `mapstructure:",squash"`

	Status OrderStatus `json:"status"`

	Notes string `json:"notes"`

	AdditionalCharges float64 `json:"additionalCharges"`

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

func (Order) TableNameForQuery() string {
	return "\"stich\".\"Orders\" E"
}
