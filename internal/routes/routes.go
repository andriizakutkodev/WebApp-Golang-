package routes

import (
	"net/http"

	g "github.com/gin-gonic/gin"
)

func RegisterRoutes(ge *g.Engine) {
	authGroup := ge.Group("/auth")
	{
		authGroup.POST("/register", handleRegister)
		authGroup.POST("/login", handleLogin)
	}
}

func handleRegister(ctx *g.Context) {
	ctx.String(http.StatusOK, "register")
}

func handleLogin(ctx *g.Context) {
	ctx.String(http.StatusOK, "login")
}
