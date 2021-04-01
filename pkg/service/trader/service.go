package trader

import (
	"context"

	"ptcg_trader/internal/errors"
	"ptcg_trader/internal/redis"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"

	"github.com/shopspring/decimal"
	"go.uber.org/fx"
)

// ServiceParams define params for create service
type ServiceParams struct {
	fx.In

	Repo  repository.Repositorier
	Match service.Matcher
	Redis redis.Redis
}

type svc struct {
	repo  repository.Repositorier
	match service.Matcher
	redis redis.Redis
}

// NewService support DI tool to create a new service instance
func NewService(params ServiceParams) service.TraderServicer {
	return &svc{
		repo:  params.Repo,
		match: params.Match,
		redis: params.Redis,
	}
}

// Get target item by item ID
func (svc *svc) GetItem(ctx context.Context, query model.ItemQuery) (model.Item, error) {
	item, err := svc.repo.GetItem(ctx, query)
	if err != nil {
		return item, err
	}

	return item, nil
}

// List items by query condition
func (svc *svc) ListItems(ctx context.Context, query model.ItemQuery) ([]model.Item, int64, error) {
	items, err := svc.repo.ListItems(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := svc.repo.CountItems(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Get target order by order ID
func (svc *svc) GetOrder(ctx context.Context, query model.OrderQuery) (model.Order, error) {
	order, err := svc.repo.GetOrder(ctx, query)
	if err != nil {
		return order, err
	}

	return order, nil
}

// List orders by query condition
func (svc *svc) ListOrders(ctx context.Context, query model.OrderQuery) ([]model.Order, int64, error) {
	orders, err := svc.repo.ListOrders(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := svc.repo.CountOrders(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Create a Order and check are there any orders can be matched
func (svc *svc) CreateOrder(ctx context.Context, order *model.Order) error {
	// // try to get redis lock
	// lockKey := fmt.Sprintf("trader.item.%d.lock", order.ItemID)
	// ok, err := svc.redis.RedisLock(ctx, lockKey, "", time.Second)
	// if err != nil {
	// 	return err
	// }
	// if !ok {
	// 	return errors.Wrap(errors.ErrDataConflict, "Order failed, please try again")
	// }
	// defer svc.redis.RedisUnlock(ctx, lockKey, "")

	// order matched! change order status and create a transaction record
	err := svc.repo.Transaction(ctx, func(ctx context.Context, txRepo repository.Repositorier) error {
		// get database row lock
		itemQuery := model.ItemQuery{
			ID:        &order.ItemID,
			ForUpdate: true,
		}
		_, err := txRepo.GetItem(ctx, itemQuery)
		if err != nil {
			return err
		}

		// check matched
		matchedOrder, err := svc.match.WithRepo(txRepo).MatchOrders(ctx, order)
		if err != nil && !errors.Is(err, errors.ErrResourceNotFound) {
			return err
		}
		// no any matched, simply create order
		if errors.Is(err, errors.ErrResourceNotFound) {
			err := txRepo.CreateOrder(ctx, order)
			if err != nil {
				return err
			}
			return nil
		}

		// update status of matched order
		orderQuery := model.OrderQuery{
			ID: &matchedOrder.ID,
		}
		orderUpdates := model.OrderUpdates{
			Status: model.OrderStatusCompleted,
		}
		err = txRepo.UpdateOrders(ctx, orderQuery, orderUpdates)
		if err != nil {
			return err
		}

		// create the new order with completed status
		order.Status = model.OrderStatusCompleted
		err = txRepo.CreateOrder(ctx, order)
		if err != nil {
			return err
		}

		// create transaction
		tx := model.Transaction{
			MakeOrderID: matchedOrder.ID,
			TakeOrderID: order.ID,
			FinalPrice:  decimal.Min(order.Price, matchedOrder.Price),
		}
		err = txRepo.CreateTransaction(ctx, &tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Update orders by query condition
func (svc *svc) UpdateOrders(ctx context.Context, query model.OrderQuery, updates model.OrderUpdates) error {
	err := svc.repo.UpdateOrders(ctx, query, updates)
	if err != nil {
		return err
	}
	return nil
}

// Get target transaction by transaction ID
func (svc *svc) GetTransaction(ctx context.Context, query model.TransactionQuery) (model.Transaction, error) {
	transaction, err := svc.repo.GetTransaction(ctx, query)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// List transactions by query condition
func (svc *svc) ListTransactions(ctx context.Context, query model.TransactionQuery) ([]model.Transaction, int64, error) {
	transactions, err := svc.repo.ListTransactions(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := svc.repo.CountTransactions(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// Create transaction
func (svc *svc) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	err := svc.repo.CreateTransaction(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}
