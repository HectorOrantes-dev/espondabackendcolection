package dependencies_coleccion

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/coleccion/application"
	"coleccionbackend/src/feature/coleccion/domain"
	"coleccionbackend/src/feature/coleccion/infraestructure/adapters"
	"coleccionbackend/src/feature/coleccion/infraestructure/controllers"
	"coleccionbackend/src/feature/coleccion/infraestructure/routers"
)

type ColeccionDependencies struct {
	createCtrl *controllers.CreateColeccionController
	listCtrl   *controllers.ListColeccionController
	updateCtrl *controllers.UpdateColeccionController
	deleteCtrl *controllers.DeleteColeccionController
	exportCtrl *controllers.ExportColeccionController
	backupCtrl *controllers.BackupColeccionController
}

func NewColeccionDependencies(db *sql.DB, imgService domain.ImageService) *ColeccionDependencies {
	repo := adapters.NewSupabaseColeccionRepository(db)

	createUC := application.NewCreateColeccionUseCase(repo, imgService)
	listUC := application.NewListColeccionUseCase(repo)
	getByIDUC := application.NewGetByIDUseCase(repo)
	updateUC := application.NewUpdateColeccionUseCase(repo, imgService)
	deleteUC := application.NewDeleteColeccionUseCase(repo, imgService)

	return &ColeccionDependencies{
		createCtrl: controllers.NewCreateColeccionController(createUC),
		listCtrl:   controllers.NewListColeccionController(listUC, getByIDUC),
		updateCtrl: controllers.NewUpdateColeccionController(updateUC),
		deleteCtrl: controllers.NewDeleteColeccionController(deleteUC),
		exportCtrl: controllers.NewExportColeccionController(listUC),
		backupCtrl: controllers.NewBackupColeccionController(listUC, imgService),
	}
}

func (d *ColeccionDependencies) RegisterRoutes(rg *gin.RouterGroup) {
	routers.RegisterColeccionRoutes(rg, d.createCtrl, d.listCtrl, d.updateCtrl, d.deleteCtrl, d.exportCtrl, d.backupCtrl)
}
