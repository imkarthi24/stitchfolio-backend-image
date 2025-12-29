package responseModel

type Order struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Status string `json:"status,omitempty"`

	CustomerId *uint     `json:"customerId,omitempty"`
	Customer   *Customer `json:"customer,omitempty"`

	OrderItems []OrderItem `json:"orderItems,omitempty"`
}

type OrderItem struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Description string  `json:"description,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Total       float64 `json:"total,omitempty"`

	OrderId uint   `json:"orderId,omitempty"`
	Order   *Order `json:"order,omitempty"`

	MeasurementId *uint        `json:"measurementId,omitempty"`
	Measurement   *Measurement `json:"measurement,omitempty"`
}
