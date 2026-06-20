package dependencies_etiquetas

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/etiquetas/application"
	"coleccionbackend/src/feature/etiquetas/infraestructure/adapters"
	"coleccionbackend/src/feature/etiquetas/infraestructure/controllers"
	"coleccionbackend/src/feature/etiquetas/infraestructure/routers"
)

type EtiquetasDependencies struct {
	controller *controllers.EtiquetasController
}

func NewEtiquetasDependencies(db *sql.DB) *EtiquetasDependencies {
	repo := adapters.NewSupabaseEtiquetasRepository(db)

	createUC := application.NewCreateEtiquetaUseCase(repo)
	listUC := application.NewListEtiquetasUseCase(repo)
	updateUC := application.NewUpdateEtiquetaUseCase(repo)
	deleteUC := application.NewDeleteEtiquetaUseCase(repo)

	ctrl := controllers.NewEtiquetasController(createUC, listUC, updateUC, deleteUC)
	return &EtiquetasDependencies{controller: ctrl}
}

func (d *EtiquetasDependencies) RegisterRoutes(rg *gin.RouterGroup) {
	routers.RegisterEtiquetasRoutes(rg, d.controller)
}
