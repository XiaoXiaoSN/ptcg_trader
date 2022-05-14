package gormrepo

import (
	"context"
	"reflect"

	"ptcg_trader/internal/errors"
	"ptcg_trader/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func buildItemQuery(db *gorm.DB, query model.ItemQuery) *gorm.DB {
	var clauses []clause.Expression

	if query.ForUpdate {
		clauses = append(clauses, selectForUpdate)
	}

	db = db.Model(&model.Item{}).
		Clauses(clauses...).Where(query)
	return db
}

// GetItem get information of the targe item
func (repo *_repository) GetItem(ctx context.Context, query model.ItemQuery) (model.Item, error) {
	var item model.Item

	err := buildItemQuery(repo.DB(ctx), query).
		First(&item).Error
	if err != nil {
		return item, notFoundOrInternalError(err)
	}
	return item, nil
}

// CountItems counting total count of items
func (repo *_repository) CountItems(ctx context.Context, query model.ItemQuery) (int64, error) {
	var total int64

	err := buildItemQuery(repo.DB(ctx), query).
		Count(&total).Error
	if err != nil {
		return total, err
	}
	return total, nil
}

// ListItems list Items
func (repo *_repository) ListItems(ctx context.Context, query model.ItemQuery) ([]model.Item, error) {
	var itemList []model.Item

	db := buildItemQuery(repo.DB(ctx), query)
	if query.PerPage <= 0 {
		query.PerPage = DefaultPerPage
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	db = db.Limit(query.PerPage).
		Offset((query.Page - 1) * query.PerPage)

	err := db.Find(&itemList).Error
	if err != nil {
		return nil, err
	}
	return itemList, nil
}

func buildOrderQuery(db *gorm.DB, query model.OrderQuery) *gorm.DB {
	db = db.Where(query)

	return db
}

// GetOrder get information of the targe order
func (repo *_repository) GetOrder(ctx context.Context, query model.OrderQuery) (model.Order, error) {
	var order model.Order

	err := buildOrderQuery(repo.DB(ctx), query).
		Model(&model.Order{}).
		First(&order).Error
	if err != nil {
		return order, notFoundOrInternalError(err)
	}
	return order, nil
}

// CountOrders counting total count of orders
func (repo *_repository) CountOrders(ctx context.Context, query model.OrderQuery) (int64, error) {
	var total int64

	err := buildOrderQuery(repo.DB(ctx), query).
		Model(&model.Order{}).
		Count(&total).Error
	if err != nil {
		return total, err
	}
	return total, nil
}

// ListOrders list Orders
func (repo *_repository) ListOrders(ctx context.Context, query model.OrderQuery) ([]model.Order, error) {
	var orderList []model.Order

	db := buildOrderQuery(repo.DB(ctx), query).
		Model(&model.Order{})

	if query.PerPage >= 0 {
		if query.PerPage == 0 {
			query.PerPage = DefaultPerPage
		}
		if query.Page == 0 {
			query.Page = 1
		}
		db = db.Limit(query.PerPage).
			Offset((query.Page - 1) * query.PerPage)
	}

	err := db.Find(&orderList).Error
	if err != nil {
		return nil, err
	}
	return orderList, nil
}

// CreateOrder create Order
func (repo *_repository) CreateOrder(ctx context.Context, order *model.Order) error {
	err := repo.DB(ctx).Create(order).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrders update Orders
func (repo *_repository) UpdateOrders(ctx context.Context, query model.OrderQuery, updates model.OrderUpdates) error {
	if reflect.DeepEqual(query, model.OrderQuery{}) {
		return errors.Wrap(errors.ErrInternalError, "UpdateOrders query is empty")
	}

	db := buildOrderQuery(repo.DB(ctx), query)

	err := db.Model(&model.OrderUpdates{}).
		Updates(&updates).Error
	if err != nil {
		return err
	}
	return nil
}

func buildTransactionQuery(db *gorm.DB, query model.TransactionQuery) *gorm.DB {
	db = db.Model(&model.Transaction{}).Where(query)

	return db
}

// GetTransaction get information of the targe transaction
func (repo *_repository) GetTransaction(ctx context.Context, query model.TransactionQuery) (model.Transaction, error) {
	var transaction model.Transaction

	err := buildTransactionQuery(repo.DB(ctx), query).
		First(&transaction).Error
	if err != nil {
		return transaction, notFoundOrInternalError(err)
	}
	return transaction, nil
}

// CountTransactions counting total count of transactions
func (repo *_repository) CountTransactions(ctx context.Context, query model.TransactionQuery) (int64, error) {
	var total int64

	err := buildTransactionQuery(repo.DB(ctx), query).
		Count(&total).Error
	if err != nil {
		return total, err
	}
	return total, nil
}

// ListTransactions list Transactions
func (repo *_repository) ListTransactions(ctx context.Context, query model.TransactionQuery) ([]model.Transaction, error) {
	var transactionList []model.Transaction

	db := buildTransactionQuery(repo.DB(ctx), query)
	if query.PerPage <= 0 {
		query.PerPage = DefaultPerPage
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	db = db.Limit(query.PerPage).
		Offset((query.Page - 1) * query.PerPage)

	err := db.Find(&transactionList).Error
	if err != nil {
		return nil, err
	}
	return transactionList, nil
}

// CreateTransaction create Transaction
func (repo *_repository) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	err := repo.DB(ctx).Model(&model.Transaction{}).
		Create(transaction).Error
	if err != nil {
		return err
	}
	return nil
}
