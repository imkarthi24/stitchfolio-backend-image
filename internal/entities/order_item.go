package entities

import "time"

type OrderItem struct {
	*Model `mapstructure:",squash"`

	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
	AdditionalCharges float64 `json:"additionalCharges"`

	ExpectedDeliveryDate *time.Time `json:"expectedDeliveryDate,omitempty"`
	DeliveredDate        *time.Time `json:"deliveredDate,omitempty"`

	PersonId *uint   `json:"personId,omitempty"`
	Person   *Person `gorm:"foreignKey:PersonId" json:"person,omitempty"`

	MeasurementId *uint        `json:"measurementId,omitempty"`
	Measurement   *Measurement `gorm:"foreignKey:MeasurementId" json:"measurement,omitempty"`

	OrderId uint   `json:"orderId"`
	Order   *Order `gorm:"foreignKey:OrderId" json:"order"`
}

func (OrderItem) TableNameForQuery() string {
	return "\"stich\".\"OrderItems\" E"
}
