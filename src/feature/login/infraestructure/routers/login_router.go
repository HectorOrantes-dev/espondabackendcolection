package routers

import (
	"coleccionbackend/src/feature/login/infraestructure/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(rg *gin.RouterGroup, ctrl *controllers.LoginController) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", ctrl.Login)
		auth.POST("/logout", ctrl.Logout)
		auth.POST("/refresh", ctrl.Refresh)
	}
}
