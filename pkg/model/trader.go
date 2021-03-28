package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Item is the model of goods, like PTCG cards
type Item struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	ImageURL  string    `json:"image_url" gorm:"column:image_url"`
	CreatorID int64     `json:"creator_id" gorm:"column:creator_id"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// ItemQuery ...
type ItemQuery struct {
	ID *int64 `json:"id" gorm:"column:id"`

	PerPage int `json:"per_page" gorm:"-"`
	Page    int `json:"page" gorm:"-"`
}

// Order is the model of user order
type Order struct {
	ID        int64           `json:"id"  gorm:"column:id"`
	ItemID    int64           `json:"item_id"  gorm:"column:item_id"`
	CreatedID int64           `json:"created_id"  gorm:"column:created_id"`
	OrderType OrderType       `json:"order_type"  gorm:"column:order_type"`
	Price     decimal.Decimal `json:"price"  gorm:"column:price"`
	Status    OrderStatus     `json:"status"  gorm:"column:status"`
	UpdatedAt time.Time       `json:"updated_at"  gorm:"column:updated_at"`
	CreatedAt time.Time       `json:"created_at"  gorm:"column:created_at"`
}

// Transaction is the result of orders matching
type Transaction struct {
	ID          int64           `json:"id"  gorm:"column:id"`
	BuyOrderID  int64           `json:"buy_order_id"  gorm:"column:buy_order_id"`
	SellOrderID int64           `json:"sell_order_id"  gorm:"column:sell_order_id"`
	FinalPrice  decimal.Decimal `json:"final_price"  gorm:"column:final_price"`
	UpdatedAt   time.Time       `json:"updated_at"  gorm:"column:updated_at"`
	CreatedAt   time.Time       `json:"created_at"  gorm:"column:created_at"`
}
