package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSConfig define echo middleware cors config
var CORSConfig = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	},
	AllowHeaders: []string{
		"*",
		echo.HeaderAuthorization,
		echo.HeaderContentType,
		echo.HeaderOrigin,
		echo.HeaderContentLength,
	},
})
