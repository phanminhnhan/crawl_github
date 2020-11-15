package middleware


import (
	"github-trending/security"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github-trending/model"
)

func UseJwtMiddleware() echo.MiddlewareFunc{
	config := middleware.JWTConfig{
		 Claims: &model.JwtCustomClaims{},
		 SigningKey: []byte(security.SecretKey),

	}
	return middleware.JWTWithConfig(config)
}