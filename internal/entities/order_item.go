package entities

type OrderItem struct {
	*Model `mapstructure:",squash"`

	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`

	OrderId uint   `json:"orderId"`
	Order   *Order `gorm:"-" json:"-"`

	MeasurementId *uint        `json:"measurementId"`
	Measurement   *Measurement `gorm:"-" json:"-"`
}

func (OrderItem) TableName() string {
	return "stitch.OrderItems"
}

func (OrderItem) TableNameForQuery() string {
	return "\"OrderItems\" E"
}
