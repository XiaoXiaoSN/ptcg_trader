package gormrepo

import (
	"context"
	"ptcg_trader/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type _matchOrder struct {
	ItemID    int64             `gorm:"column:item_id"`
	OrderType model.OrderType   `gorm:"column:order_type"`
	Status    model.OrderStatus `gorm:"column:status"`

	Price     clause.Expression `gorm:"-"`
	ForUpdate bool              `gorm:"-"`
	ForShare  bool              `gorm:"-"`
}

func (w *_matchOrder) Where(db *gorm.DB) *gorm.DB {
	var clauses []clause.Expression

	if w.Price != nil {
		clauses = append(clauses, w.Price)
	}
	if w.ForUpdate {
		clauses = append(clauses, selectForUpdate)
	} else if w.ForShare {
		clauses = append(clauses, selectForShare)
	}

	db = db.Clauses(clauses...).Where(w)
	return db
}

// MatchOrders check that are there two orders can be matched
func (repo *_repository) MatchOrders(ctx context.Context, order *model.Order) (model.Order, error) {
	var matchedOrder model.Order

	_matchOrder := _matchOrder{
		ItemID: order.ItemID,
		Status: model.OrderStatusProgress,
	}
	// buy order should match the other sell order, and vice versa
	if order.OrderType == model.OrderTypeBuy {
		_matchOrder.OrderType = model.OrderTypeSell
		_matchOrder.Price = &clause.Lte{
			Column: "price",
			Value:  order.Price,
		}
	} else {
		_matchOrder.OrderType = model.OrderTypeBuy
		_matchOrder.Price = &clause.Gte{
			Column: "price",
			Value:  order.Price,
		}
	}

	err := repo.DB(ctx).
		Scopes(_matchOrder.Where).
		Order("id DESC").
		First(&matchedOrder).Error
	if err != nil {
		return matchedOrder, notFoundOrInternalError(err)
	}

	return matchedOrder, nil
}
