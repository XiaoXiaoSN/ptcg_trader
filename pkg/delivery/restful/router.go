package restful

import (
	"github.com/labstack/echo/v4"
)

// SetRoutes ...
func SetRoutes(e *echo.Echo, h *Handler) {
	gAPIv1 := e.Group("apis/v1")
	{
		gAPIv1.GET("/items/:id", h.GetItemEndpoint)
	}
}
