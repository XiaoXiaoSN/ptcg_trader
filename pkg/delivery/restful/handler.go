package restful

import (
	"net/http"

	"ptcg_trader/internal/ctxutil"
	"ptcg_trader/internal/errors"
	"ptcg_trader/pkg/model"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type getItemReq struct {
	ID int64 `param:"itemID"`
}

type getItemResp struct {
	Data model.Item `json:"data"`
}

// @Title getItemEndpoint
// @Description 取得指定 商品ID 的商品
// @
// @Param  item_id   path    uint32  true    "商品ID"
// @
// @Success  200  object  getItemResp           "商品資訊"
// @Failure  400  object  StautsBadRequestResp  "不正確的查詢參數"
// @Failure  404  object  StautsNotFoundResp    "查詢的資源不存在"
// @
// @Resource 關於商品們
// @Router /apis/v1/items/{itemID} [get]
func (h Handler) getItemEndpoint(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// receive request paraments
	var req getItemReq
	err = c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	if req.ID <= 0 {
		return errors.Wrap(errors.ErrBadRequest, "param itemID invalid")
	}

	query := model.ItemQuery{
		ID: &req.ID,
	}
	item, err := h.svc.GetItem(ctx, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, getItemResp{
		Data: item,
	})
}

type listItemReq struct {
	PaginationQuery
}

type listItemResp struct {
	Meta listItemRespMeta `json:"meta"`
	Data []model.Item     `json:"data"`
}

type listItemRespMeta struct {
	Pagination Pagination `json:"pagination"`
}

// @Title List Items
// @Description 列表可以買賣的商品清單
// @
// @Param  page      query   uint32  false   "頁數 (default: 1)"
// @Param  per_page  query   uint32  false   "一頁顯示幾筆 (default: 50)"
// @
// @Success  200  object  listItemResp  "商品列表"
// @Failure  400  object  StautsBadRequestResp  "不正確的查詢參數"
// @
// @Resource 關於商品們
// @Router /apis/v1/items [get]
func (h Handler) listItemEndpoint(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// receive request paraments
	var req listItemReq
	err = c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	err = req.PaginationQuery.ValidateAndSet()
	if err != nil {
		return err
	}

	query := model.ItemQuery{
		PerPage: req.PerPage,
		Page:    req.Page,
	}
	itemList, total, err := h.svc.ListItems(ctx, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, listItemResp{
		Data: itemList,
		Meta: listItemRespMeta{
			Pagination: NewPagination(req.Page, req.PerPage, int(total)),
		},
	})
}

type getOrderReq struct {
	ID int64 `param:"orderID"`
}

type getOrderResp struct {
	Data model.Order `json:"data"`
}

// @Title getOrderEndpoint
// @Description 取得指定 訂單ID 的訂單
// @
// @Param  order_id   path    uint32  true    "訂單ID"
// @
// @Success  200  object  getOrderResp           "訂單資訊"
// @Failure  400  object  StautsBadRequestResp  "不正確的查詢參數"
// @Failure  404  object  StautsNotFoundResp    "查詢的資源不存在"
// @
// @Resource 關於訂單們
// @Router /apis/v1/orders/{orderID} [get]
func (h Handler) getOrderEndpoint(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// receive request paraments
	var req getOrderReq
	err = c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	if req.ID <= 0 {
		return errors.Wrap(errors.ErrBadRequest, "param orderID invalid")
	}

	query := model.OrderQuery{
		ID: &req.ID,
	}
	order, err := h.svc.GetOrder(ctx, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, getOrderResp{
		Data: order,
	})
}

type createOrderReq struct {
	ItemID    int64           `json:"item_id" gorm:"column:item_id"`
	OrderType model.OrderType `json:"order_type" gorm:"column:order_type"`
	Price     string          `json:"price" gorm:"column:price"`
}

type createOrderResp struct {
	Data model.Order `json:"data"`
}

// @Title getOrderEndpoint
// @Description 建立訂單
// @
// @Param  reqBody  body    createOrderReq  true   "建立訂單資訊"
// @
// @Success  200  object  createOrderResp       "訂單資訊"
// @Failure  400  object  StautsBadRequestResp  "不正確的查詢參數"
// @Failure  404  object  StautsNotFoundResp    "查詢的資源不存在"
// @
// @Resource 關於訂單們
// @Router /apis/v1/orders [post]
func (h Handler) createOrderEndpoint(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// receive request paraments
	var req createOrderReq
	err = c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	if req.OrderType == 0 {
		return errors.Wrap(errors.ErrBadRequest, "order type cannot be empty")
	}

	var orderPrice decimal.Decimal
	orderPrice, err = decimal.NewFromString(req.Price)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	// check params range: 0 < priceOrder <= 10
	if orderPrice.LessThanOrEqual(decimal.Zero) {
		return errors.Wrap(errors.ErrBadRequest, "order Price is less then or equal zero")
	}
	if orderPrice.GreaterThan(decimal.NewFromInt(10)) {
		return errors.Wrap(errors.ErrBadRequest, "order Price is greater then 10")
	}

	// var creatorID int64
	creatorID, err := ctxutil.IdentityIDFromCtx(ctx)
	if err != nil {
		return err
	}

	order := &model.Order{
		ItemID:    req.ItemID,
		CreatorID: creatorID,
		OrderType: req.OrderType,
		Price:     orderPrice,
		Status:    model.OrderStatusProgress,
	}
	err = h.svc.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, getOrderResp{
		Data: *order,
	})
}

type listOrderReq struct {
	PaginationQuery
}

type listOrderResp struct {
	Meta listOrderRespMeta `json:"meta"`
	Data []model.Order     `json:"data"`
}

type listOrderRespMeta struct {
	Pagination Pagination `json:"pagination"`
}

// @Title List Orders
// @Description 列表下定買單、賣單的訂單們
// @
// @Param  page      query   uint32  false   "頁數 (default: 1)"
// @Param  per_page  query   uint32  false   "一頁顯示幾筆 (default: 50)"
// @
// @Success  200  object  listOrderResp  "訂單列表"
// @Failure  400  object  StautsBadRequestResp  "不正確的查詢參數"
// @
// @Resource 關於訂單們
// @Router /apis/v1/orders [get]
func (h Handler) listOrderEndpoint(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// receive request paraments
	var req listOrderReq
	err = c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	err = req.PaginationQuery.ValidateAndSet()
	if err != nil {
		return err
	}

	query := model.OrderQuery{
		PerPage: req.PerPage,
		Page:    req.Page,
	}
	orderList, total, err := h.svc.ListOrders(ctx, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, listOrderResp{
		Data: orderList,
		Meta: listOrderRespMeta{
			Pagination: NewPagination(req.Page, req.PerPage, int(total)),
		},
	})
}

type listTransactionReq struct {
	ItemID int64 `query:"item_id"`

	PaginationQuery
}

type listTransactionResp struct {
	Meta listTransactionRespMeta `json:"meta"`
	Data []model.Transaction     `json:"data"`
}

type listTransactionRespMeta struct {
	Pagination Pagination `json:"pagination"`
}

// @Title List Transactions
// @Description 列表交易紀錄們
// @
// @Param  page      query   uint32  false   "頁數 (default: 1)"
// @Param  per_page  query   uint32  false   "一頁顯示幾筆 (default: 50)"
// @Param  item_id   query   uint32  false   "商品ID"
// @
// @Success  200  object  listTransactionResp  "交易紀錄列表"
// @Failure  400  object  StautsBadRequestResp  "不正確的查詢參數"
// @
// @Resource 關於交易紀錄們
// @Router /apis/v1/transactions [get]
func (h Handler) listTransactionEndpoint(c echo.Context) (err error) {
	ctx := c.Request().Context()

	// receive request paraments
	var req listTransactionReq
	err = c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrBadRequest, err.Error())
	}
	err = req.PaginationQuery.ValidateAndSet()
	if err != nil {
		return err
	}

	query := model.TransactionQuery{
		ItemID: &req.ItemID,

		PerPage: req.PerPage,
		Page:    req.Page,
	}
	transactionList, total, err := h.svc.ListTransactions(ctx, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, listTransactionResp{
		Data: transactionList,
		Meta: listTransactionRespMeta{
			Pagination: NewPagination(req.Page, req.PerPage, int(total)),
		},
	})
}
