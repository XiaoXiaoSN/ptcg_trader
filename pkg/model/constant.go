package model

// OrderType define order type
type OrderType int8

// OrderType enum
const (
	OrderTypeBuy  OrderType = 1
	OrderTypeSell OrderType = 2
)

// OrderStatus define order status
type OrderStatus int8

// OrderStatus enum
const (
	OrderStatusProgress  OrderStatus = 1
	OrderStatusCompleted OrderStatus = 2
	OrderStatusCancel    OrderStatus = 3
)
