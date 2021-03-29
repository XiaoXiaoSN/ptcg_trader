package restful

// @Version 0.0.1
// @Title PTCG Trader API v1
// @Description PTCG Trader API
// @ContactName Xiao.Xiao
// @ContactEmail freedom85812@gmail.com
// @Server http://www.fake.com FakeServerHost

import (
	"ptcg_trader/pkg/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// SetRoutes ...
func SetRoutes(e *echo.Echo, h *Handler) {
	gAPIv1 := e.Group("apis/v1")
	{
		gAPIv1.GET("/items", h.listItemEndpoint)
		gAPIv1.GET("/items/:itemID", h.getItemEndpoint)

		gAPIv1.GET("/orders", h.listOrderEndpoint)
		gAPIv1.GET("/orders/:orderID", h.getOrderEndpoint)
		gAPIv1.POST("/orders", h.createOrderEndpoint)

		gAPIv1.GET("/transactions", h.listTransactionEndpoint)
	}
}

// HandlerParams define params for create handler
type HandlerParams struct {
	fx.In

	Svc service.TraderServicer
}

// Handler http restful handler
type Handler struct {
	svc service.TraderServicer
}

// NewHandler http handler injection
func NewHandler(params HandlerParams) *Handler {
	server := Handler{
		svc: params.Svc,
	}
	return &server
}
