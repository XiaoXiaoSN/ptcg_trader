package matcher

import (
	"context"
	"fmt"
	"runtime"

	"ptcg_trader/internal/errors"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"go.uber.org/fx"
)

// MatchParams define params for create service
type MatchParams struct {
	fx.In

	Repo repository.Repositorier
}

type svc struct {
	repo repository.Repositorier

	buyOrderEngines  map[int64]*OrderEngine
	sellOrderEngines map[int64]*OrderEngine
}

// NewMatch support DI tool to create a new service instance
func NewMatch(params MatchParams) service.Matcher {
	return &svc{
		repo: params.Repo,

		buyOrderEngines:  make(map[int64]*OrderEngine), // map[itemID]order-RB-Tree
		sellOrderEngines: make(map[int64]*OrderEngine), // map[itemID]order-RB-Tree
	}
}

func (svc *svc) WithRepo(repo repository.Repositorier) service.Matcher {
	if repo == nil {
		repo = svc.repo
	}

	return NewMatch(MatchParams{
		Repo: repo,
	})
}

func (svc *svc) MatchOrders(ctx context.Context, order *model.Order) (*model.Order, error) {
	matchedOrder, err := svc.repo.MatchOrders(ctx, order)
	if err != nil {
		return nil, err
	}

	return &matchedOrder, nil
}

// AsyncMatchOrders ...
// NOTE: A trader can only buy or sell 1 card in 1 order. QQ
func (svc *svc) AsyncMatchOrders(ctx context.Context, order *model.Order) error {
	// save order to redis?

	takeSideOrderEngine, makeSideOrderEngine := svc.getOrderEngines(ctx, order)
	if takeSideOrderEngine == nil || makeSideOrderEngine == nil {
		return errors.WithMessage(errors.ErrBadRequest, "unknown order type")
	}
	makeSideOrderEngine.lock.Lock()
	defer makeSideOrderEngine.lock.Unlock()

	// not matched, insert into rb-tree
	if takeSideOrderEngine.Tree.Size() == 0 {
		subTree := rbt.NewWith(utils.Int64Comparator)
		subTree.Put(order.ID, order)
		makeSideOrderEngine.Tree.Put(order.Price, subTree)

		return nil
	}

	takeSideOrderEngine.lock.Lock()
	go svc.asyncMatchOrders(ctx, order, takeSideOrderEngine)
	return nil
}

func (svc *svc) asyncMatchOrders(ctx context.Context, order *model.Order, takeSideOrderEngine *OrderEngine) (err error) {
	defer func() {
		recoverError()
		if err != nil {
			log.Error().Msgf("asyncMatchOrders fail: %+v", err)
		}
	}()

	// takeOrder matched a makeOrder, pop it and update their db status
	var makeOrder *model.Order
	var (
		timeTree                  *rbt.Tree
		orderTreeKey, timeTreeKey interface{}
	)
	switch order.OrderType {
	case model.OrderTypeBuy:
		// buy order, find the order with lowest price
		orderNode := takeSideOrderEngine.Tree.Left()
		orderTreeKey = orderNode.Key

		timeTree := orderNode.Value.(*rbt.Tree)
		timeNode := timeTree.Left()
		timeTreeKey = timeNode.Key
		makeOrder = timeNode.Value.(*model.Order)

	case model.OrderTypeSell:
		// sell order, find the order with highest price
		orderNode := takeSideOrderEngine.Tree.Right()
		orderTreeKey = orderNode.Key

		timeTree := orderNode.Value.(*rbt.Tree)
		timeNode := timeTree.Left()
		timeTreeKey = timeNode.Key
		makeOrder = timeNode.Value.(*model.Order)

	default:
		return errors.WithMessagef(errors.ErrBadRequest, "Unknown order type: %d", order.OrderType)
	}

	// order matched! change order status and create a transaction record
	err = svc.repo.Transaction(ctx, func(ctx context.Context, txRepo repository.Repositorier) error {
		// update status of matched order
		orderQuery := model.OrderQuery{
			ID: &makeOrder.ID,
		}
		orderUpdates := model.OrderUpdates{
			Status: model.OrderStatusCompleted,
		}
		err := txRepo.UpdateOrders(ctx, orderQuery, orderUpdates)
		if err != nil {
			return err
		}

		// create the new order with completed status
		orderQuery.ID = &order.ID
		err = txRepo.UpdateOrders(ctx, orderQuery, orderUpdates)
		if err != nil {
			return err
		}

		// create transaction
		tx := model.Transaction{
			ItemID:      order.ItemID,
			MakeOrderID: makeOrder.ID,
			TakeOrderID: order.ID,
			FinalPrice:  decimal.Min(order.Price, makeOrder.Price),
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

	// remove makeOrder from takeSideOrderEngine
	if timeTree.Size() == 1 {
		takeSideOrderEngine.Tree.Remove(orderTreeKey)
	} else {
		timeTree.Remove(timeTreeKey)
	}

	return nil
}

func (svc *svc) getOrderEngines(ctx context.Context, order *model.Order) (takeSide *OrderEngine, makeSide *OrderEngine) {
	switch order.OrderType {
	case model.OrderTypeBuy:
		var ok bool
		takeSide, ok = svc.sellOrderEngines[order.ItemID]
		if !ok {
			takeSide = NewOrderEngine()
			svc.sellOrderEngines[order.ItemID] = takeSide
		}
		makeSide, ok = svc.buyOrderEngines[order.ItemID]
		if !ok {
			makeSide = NewOrderEngine()
			svc.buyOrderEngines[order.ItemID] = makeSide
		}

	case model.OrderTypeSell:
		var ok bool
		takeSide, ok = svc.buyOrderEngines[order.ItemID]
		if !ok {
			takeSide = NewOrderEngine()
			svc.buyOrderEngines[order.ItemID] = takeSide
		}
		makeSide, ok = svc.buyOrderEngines[order.ItemID]
		if !ok {
			makeSide = NewOrderEngine()
			svc.buyOrderEngines[order.ItemID] = makeSide
		}
	}

	return takeSide, makeSide
}

func recoverError() {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg = msg + fmt.Sprintf("%s:%d\n", file, line)
		}
		log.Error().Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
	}
}
