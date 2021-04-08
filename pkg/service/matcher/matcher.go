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
	readyFlag chan struct{}

	repo repository.Repositorier

	buyOrderEngines  map[int64]*OrderEngine
	sellOrderEngines map[int64]*OrderEngine
}

// NewMatch support DI tool to create a new service instance
func NewMatch(params MatchParams) service.Matcher {
	svc := &svc{
		readyFlag: make(chan struct{}),
		repo:      params.Repo,

		buyOrderEngines:  make(map[int64]*OrderEngine), // map[itemID]order-RB-Tree
		sellOrderEngines: make(map[int64]*OrderEngine), // map[itemID]order-RB-Tree
	}

	ctx := log.Logger.WithContext(context.Background())
	go svc.loadUncompletedOrders(ctx)

	return svc
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

// loadUncompletedOrders from database load orders in memory
func (svc *svc) loadUncompletedOrders(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("loadUncompletedOrders start")
	defer log.Ctx(ctx).Info().Msg("loadUncompletedOrders complate.")

	query := model.OrderQuery{
		Status:  model.OrderStatusProgress,
		PerPage: -1,
	}
	orders, err := svc.repo.ListOrders(ctx, query)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("load uncompleted order failed: %+v", err)
		return err
	}

	for i := range orders {
		switch orders[i].OrderType {
		case model.OrderTypeBuy:
			oe, isExist := svc.buyOrderEngines[orders[i].ItemID]
			if !isExist {
				oe = NewOrderEngine()
				svc.buyOrderEngines[orders[i].ItemID] = oe
			}
			oe.Append(&orders[i])

		case model.OrderTypeSell:
			oe, isExist := svc.sellOrderEngines[orders[i].ItemID]
			if !isExist {
				oe = NewOrderEngine()
				svc.sellOrderEngines[orders[i].ItemID] = oe
			}
			oe.Append(&orders[i])
		}
	}

	// statistics orderEngines
	for itemID, oe := range svc.sellOrderEngines {
		log.Ctx(ctx).Debug().
			Int64("item_id", itemID).
			Msgf("sellOrderEnine has %d orders", oe.Size())
	}
	for itemID, oe := range svc.buyOrderEngines {
		log.Ctx(ctx).Debug().
			Int64("item_id", itemID).
			Msgf("buyOrderEnine has %d orders", oe.Size())
	}

	close(svc.readyFlag)
	return nil
}

// AsyncMatchOrders ...
// NOTE: A trader can only buy or sell 1 card in 1 order. QQ
func (svc *svc) AsyncMatchOrders(ctx context.Context, order *model.Order) error {
	<-svc.readyFlag // wait init

	takeSideOrderEngine, makeSideOrderEngine := svc.getOrderEngines(ctx, order)
	if takeSideOrderEngine == nil || makeSideOrderEngine == nil {
		return errors.WithMessage(errors.ErrBadRequest, "unknown order type")
	}

	err := svc.asyncMatchOrders(ctx, order, makeSideOrderEngine, takeSideOrderEngine)
	if err != nil {
		log.Error().Msgf("excute asyncMatchOrders fail: %+v", err)
	}

	return nil
}

func (svc *svc) asyncMatchOrders(ctx context.Context, takeOrder *model.Order, makeSideOrderEngine, takeSideOrderEngine *OrderEngine) (err error) {
	defer func() {
		if _err := recoverPanic(); _err != nil {
			err = _err
		}
	}()

	logger := log.Ctx(ctx).With().
		Int64("take_order_id", takeOrder.ID).
		Str("take_order_price", takeOrder.Price.String()).
		Int8("take_order_side", int8(takeOrder.OrderType)).
		Logger()

	// not matched, insert into rb-tree
	if makeSideOrderEngine.Tree.Size() == 0 {
		logger.Debug().Msg("The makeOrderEngine is empty")
		takeSideOrderEngine.Append(takeOrder)
		return nil
	}

	// takeOrder matched a makeOrder, pop it and update their db status
	var makeOrder *model.Order
	var (
		timeTree                  *rbt.Tree
		orderTreeKey, timeTreeKey interface{}
	)
	switch takeOrder.OrderType {
	case model.OrderTypeBuy:
		// buy order:
		// find that there exists an uncompleted sell order, whose price is the lowest one
		// among all uncompleted sell orders and less than or equal to the price of the buy order.
		orderNode := makeSideOrderEngine.Tree.Left()
		if orderNode.Key.(decimal.Decimal).LessThanOrEqual(takeOrder.Price) {
			orderTreeKey = orderNode.Key

			timeTree = orderNode.Value.(*rbt.Tree)
			timeNode := timeTree.Left()
			timeTreeKey = timeNode.Key
			makeOrder = timeNode.Value.(*model.Order)
		}

	case model.OrderTypeSell:
		// sell order:
		// there exists an uncompleted buy order, whose price is the highest one
		// among all uncompleted buy orders and greater than or equal to the price of the sell order.
		orderNode := makeSideOrderEngine.Tree.Right()
		if orderNode.Key.(decimal.Decimal).GreaterThanOrEqual(takeOrder.Price) {
			orderTreeKey = orderNode.Key

			timeTree = orderNode.Value.(*rbt.Tree)
			timeNode := timeTree.Left()
			timeTreeKey = timeNode.Key
			makeOrder = timeNode.Value.(*model.Order)
		}

	default:
		return errors.WithMessagef(errors.ErrBadRequest, "Unknown order type: %d", takeOrder.OrderType)
	}

	// the input take order not matched any make order, put it into orderEngine
	if makeOrder == nil {
		logger.Debug().Msg("There are not matched any order in the makeOrderEngine")
		takeSideOrderEngine.Append(takeOrder)
		return nil
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
		orderQuery.ID = &takeOrder.ID
		err = txRepo.UpdateOrders(ctx, orderQuery, orderUpdates)
		if err != nil {
			return err
		}

		// create transaction
		tx := model.Transaction{
			ItemID:      takeOrder.ItemID,
			MakeOrderID: makeOrder.ID,
			TakeOrderID: takeOrder.ID,
			FinalPrice:  decimal.Min(takeOrder.Price, makeOrder.Price),
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

	logger.Debug().
		Int64("make_order_id", makeOrder.ID).
		Str("make_order_price", makeOrder.Price.String()).
		Int8("make_order_side", int8(makeOrder.OrderType)).
		Msg("takeOrder matched a make order in the makeOrderEngine")

	// remove makeOrder from makeSideOrderEngine
	if timeTree.Size() == 1 {
		makeSideOrderEngine.Tree.Remove(orderTreeKey)
	} else {
		timeTree.Remove(timeTreeKey)
	}

	return nil
}

func (svc *svc) getOrderEngines(ctx context.Context, order *model.Order) (takeSide *OrderEngine, makeSide *OrderEngine) {
	switch order.OrderType {
	case model.OrderTypeBuy:
		var ok bool
		takeSide, ok = svc.buyOrderEngines[order.ItemID]
		if !ok {
			takeSide = NewOrderEngine()
			svc.buyOrderEngines[order.ItemID] = takeSide
		}
		makeSide, ok = svc.sellOrderEngines[order.ItemID]
		if !ok {
			makeSide = NewOrderEngine()
			svc.sellOrderEngines[order.ItemID] = makeSide
		}

	case model.OrderTypeSell:
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
	}

	return takeSide, makeSide
}

func recoverPanic() error {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg = msg + fmt.Sprintf("%s:%d\n", file, line)
		}
		return fmt.Errorf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
	}
	return nil
}
