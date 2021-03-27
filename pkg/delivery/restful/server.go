package restful

import (
	"net/http"

	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/service"

	"github.com/labstack/echo/v4"
)

// Handler http restful handler
type Handler struct {
	service service.TraderServicer
}

// NewHandler http handler injection
func NewHandler(
	service service.TraderServicer,
) *Handler {
	server := Handler{
		service: service,
	}
	return &server
}

// #########################################################
// GET  /apis/v1/items/:id  HTTP/1.1
// #########################################################

type getItemResp struct {
	Data model.Item
}

// GetItemEndpoint ...
func (h Handler) GetItemEndpoint(c echo.Context) (err error) {
	// TODO: impl it

	return c.JSON(http.StatusOK, getItemResp{})
}
