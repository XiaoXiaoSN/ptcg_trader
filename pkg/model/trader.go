package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Item is the model of goods, like PTCG cards
type Item struct {
	ID        int64     `json:"id" gorm:"column:id" example:"1"`
	Name      string    `json:"name" gorm:"column:name" example:"Pikachu"`
	ImageURL  string    `json:"image_url" gorm:"column:image_url" example:"https://imgur.com/NTSEJxX"`
	CreatorID int64     `json:"creator_id" gorm:"column:creator_id" example:"1"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// ItemQuery ...
type ItemQuery struct {
	ID *int64 `json:"id" gorm:"column:id"`

	PerPage int `json:"per_page" gorm:"-"`
	Page    int `json:"page" gorm:"-"`

	ForUpdate bool `gorm:"-"`
}

// Order is the model of user order
type Order struct {
	ID        int64           `json:"id" gorm:"column:id" example:"1"`
	ItemID    int64           `json:"item_id" gorm:"column:item_id" example:"1"`
	CreatorID int64           `json:"creator_id" gorm:"column:creator_id" example:"1"`
	OrderType OrderType       `json:"order_type" gorm:"column:order_type" example:"2"`
	Price     decimal.Decimal `json:"price" gorm:"column:price"`
	Status    OrderStatus     `json:"status" gorm:"column:status"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
}

// OrderQuery ...
type OrderQuery struct {
	ID     *int64      `json:"id" gorm:"column:id"`
	Status OrderStatus `json:"status" gorm:"column:status"`

	PerPage int `json:"per_page" gorm:"-"`
	Page    int `json:"page" gorm:"-"`
}

// OrderUpdates ...
type OrderUpdates struct {
	Status OrderStatus `json:"status"  gorm:"column:status"`
}

// TableName impl gorm Tabler
func (*OrderUpdates) TableName() string {
	return "orders"
}

// Transaction is the result of orders matching
type Transaction struct {
	ID          int64           `json:"id" gorm:"column:id" example:"1"`
	ItemID      int64           `json:"item_id" gorm:"column:item_id" example:"1"`
	MakeOrderID int64           `json:"make_order_id" gorm:"column:make_order_id" example:"1"`
	TakeOrderID int64           `json:"take_order_id" gorm:"column:take_order_id" example:"2"`
	FinalPrice  decimal.Decimal `json:"final_price" gorm:"column:final_price"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at"`
}

// TransactionQuery ...
type TransactionQuery struct {
	ID     *int64 `json:"id" gorm:"column:id"`
	ItemID *int64 `json:"item_id" gorm:"column:item_id"`

	PerPage int `json:"per_page" gorm:"-"`
	Page    int `json:"page" gorm:"-"`
}
