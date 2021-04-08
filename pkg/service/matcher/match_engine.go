package matcher

import (
	"ptcg_trader/pkg/model"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"github.com/shopspring/decimal"
)

// NewOrderEngine create a order tree instance
// each nodes contain another rb-tree that be sorted by createdTime (order ID)
func NewOrderEngine() *OrderEngine {
	rbtTree := rbt.NewWith(orderComparator)

	return &OrderEngine{
		Tree: rbtTree,
		lock: &sync.Mutex{},
	}
}

// OrderEngine ...
type OrderEngine struct {
	Tree *rbt.Tree
	lock *sync.Mutex
}

// Append a new order into orderEngine
func (oe *OrderEngine) Append(order *model.Order) {
	subTree, found := oe.Tree.Get(order.Price)
	if found {
		timeTree := subTree.(*rbt.Tree)
		timeTree.Put(order.ID, order)
	}

	timeTree := rbt.NewWith(utils.Int64Comparator)
	timeTree.Put(order.ID, order)
	oe.Tree.Put(order.Price, timeTree)
}

// String make OrderEngine readable
func (oe *OrderEngine) String() string {
	if oe == nil || oe.Tree == nil {
		return ""
	}
	return oe.Tree.String()
}

// Size return the actual size
func (oe *OrderEngine) Size() int {
	if oe == nil || oe.Tree == nil {
		return 0
	}

	var size int

	it := oe.Tree.Iterator()
	for it.Next() {
		orderTree := it.Value().(*rbt.Tree)
		size += orderTree.Size()
	}
	return size
}

// orderComparator compare order's parice and creating time
func orderComparator(a, b interface{}) int {
	aWeight := a.(decimal.Decimal)
	bWeight := b.(decimal.Decimal)

	return aWeight.Cmp(bWeight)
}
