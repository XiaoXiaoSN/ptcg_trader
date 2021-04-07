package matcher

import (
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
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

type OrderEngine struct {
	Tree *rbt.Tree
	lock *sync.Mutex
}

// orderComparator compare order's parice and creating time
func orderComparator(a, b interface{}) int {
	aWeight := a.(decimal.Decimal)
	bWeight := b.(decimal.Decimal)

	return aWeight.Cmp(bWeight)
}

// TODO: load from database
