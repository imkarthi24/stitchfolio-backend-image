package entities

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	*Model `mapstructure:",squash"`

	Status OrderStatus `json:"status"`

	CustomerId *uint     `json:"customerId"`
	Customer   *Customer `gorm:"-" json:"-"`

	OrderItems []OrderItem `json:"-"`
}

func (Order) TableName() string {
	return "Orders"
}

func (Order) TableNameForQuery() string {
	return "\"Orders\" E"
}
