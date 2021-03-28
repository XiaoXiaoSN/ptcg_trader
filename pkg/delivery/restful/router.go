// @Version 0.0.1
// @Title PTCG Trader API v1
// @Description PTCG Trader API
// @ContactName Xiao.Xiao
// @ContactEmail freedom85812@gmail.com
// @Server http://www.fake.com FakeServerHost
package restful

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
	}
}

// RestfulHandlerParams define params for create handler
type RestfulHandlerParams struct {
	fx.In

	Svc service.TraderServicer
}

// Handler http restful handler
type Handler struct {
	svc service.TraderServicer
}

// NewHandler http handler injection
func NewHandler(params RestfulHandlerParams) *Handler {
	server := Handler{
		svc: params.Svc,
	}
	return &server
}
