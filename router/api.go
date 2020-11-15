package router

import (
	"github-trending/handler"
	"github-trending/middleware"
	"github.com/labstack/echo"
)

type API struct {
	Echo *echo.Echo
	UserHandler handler.UserHandler
	RepoHandler handler.RepoHandler
}

func (api *API) SetupRouter (){

	api.Echo.POST("/user/sign-up",api.UserHandler.HandleSignUp )
	api.Echo.POST("/user/sign-in",api.UserHandler.HandleSignIn )
	api.Echo.GET("user/test", api.UserHandler.TestJWT, middleware.UseJwtMiddleware())
	api.Echo.GET("user/user-profile", api.UserHandler.UserProfile, middleware.UseJwtMiddleware())
	api.Echo.PUT("user/user-profile/update", api.UserHandler.UpdateUser, middleware.UseJwtMiddleware())

	github := api.Echo.Group("/github", middleware.UseJwtMiddleware())
	github.GET("/trending", api.RepoHandler.RepoTrending)

	// bookmark
	bookmark := api.Echo.Group("/bookmark", middleware.UseJwtMiddleware())
	bookmark.GET("/list", api.RepoHandler.SelectBookmarks)
	bookmark.POST("/add", api.RepoHandler.Bookmark)
	bookmark.DELETE("/delete", api.RepoHandler.DelBookmark)
}