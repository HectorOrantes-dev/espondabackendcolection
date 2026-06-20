package routers

import (
	"coleccionbackend/src/core/middleware"
	"coleccionbackend/src/feature/etiquetas/infraestructure/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterEtiquetasRoutes(rg *gin.RouterGroup, ctrl *controllers.EtiquetasController) {
	et := rg.Group("/etiquetas")
	et.Use(middleware.JWTAuth(), middleware.AuditLogger())
	{
		et.POST("", ctrl.Create)
		et.GET("", ctrl.List)
		et.PUT("/:id", ctrl.Update)
		et.DELETE("/:id", ctrl.Delete)
	}
}
