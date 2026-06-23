package routers

import (
	"coleccionbackend/src/core/middleware"
	"coleccionbackend/src/feature/coleccion/infraestructure/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterColeccionRoutes(
	rg *gin.RouterGroup,
	createCtrl *controllers.CreateColeccionController,
	listCtrl *controllers.ListColeccionController,
	updateCtrl *controllers.UpdateColeccionController,
	deleteCtrl *controllers.DeleteColeccionController,
	exportCtrl *controllers.ExportColeccionController,
	backupCtrl *controllers.BackupColeccionController,
) {
	col := rg.Group("/coleccion")
	col.Use(middleware.JWTAuth(), middleware.AuditLogger())
	{
		col.POST("", createCtrl.Handle)
		col.GET("", listCtrl.List)
		col.GET("/resumen", listCtrl.Resumen)
		col.GET("/export", exportCtrl.Handle)
		col.GET("/backup", backupCtrl.Handle)
		col.GET("/:id", listCtrl.GetByID)
		col.PUT("/:id", updateCtrl.Handle)
		col.DELETE("/:id", deleteCtrl.Handle)
	}
}
