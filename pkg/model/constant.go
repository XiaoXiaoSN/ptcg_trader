package model

// OrderType define order type
type OrderType int8

const (
	OrderTypeBuy  OrderType = 1
	OrderTypeSell OrderType = 2
)

// OrderStatus define order status
type OrderStatus int8

const (
	OrderStatusProgress  OrderStatus = 1
	OrderStatusCompleted OrderStatus = 2
)
