package restful

import (
	"net/http"

	"ptcg_trader/internal/errors"
	"ptcg_trader/pkg/model"

	"github.com/labstack/echo/v4"
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
// @Param  itemID    path    uint32  true    "商品ID"
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
// @Param  perPage   query   uint32  false   "一頁顯示幾筆 (default: 50)"
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
