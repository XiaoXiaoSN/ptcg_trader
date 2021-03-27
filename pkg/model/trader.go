package model

// Item is the model of goods, like PTCG cards
type Item struct {
	ID int64 `json:"id" gorm:"column:id"`
}
