package entities

type OrderItem struct {
	*Model `mapstructure:",squash"`

	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`

	OrderId uint   `json:"orderId"`
	Order   *Order `gorm:"foreignKey:OrderId" json:"order"`

	MeasurementId *uint        `json:"measurementId"`
	Measurement   *Measurement `gorm:"foreignKey:MeasurementId" json:"measurement"`
}

func (OrderItem) TableName() string {
	return "OrderItems"
}

func (OrderItem) TableNameForQuery() string {
	return "\"OrderItems\" E"
}
