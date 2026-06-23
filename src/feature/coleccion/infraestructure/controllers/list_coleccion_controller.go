package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/coleccion/application"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type ListColeccionController struct {
	listUC    *application.ListColeccionUseCase
	getByIDUC *application.GetByIDUseCase
	resumenUC *application.ResumenColeccionUseCase
}

func NewListColeccionController(
	l *application.ListColeccionUseCase,
	g *application.GetByIDUseCase,
	r *application.ResumenColeccionUseCase,
) *ListColeccionController {
	return &ListColeccionController{listUC: l, getByIDUC: g, resumenUC: r}
}

// Resumen godoc
// @Summary Resumen de la colección (cantidad y valor total)
// @Tags coleccion
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entities.ResumenColeccion
// @Router /coleccion/resumen [get]
func (ctrl *ListColeccionController) Resumen(c *gin.Context) {
	resumen, err := ctrl.resumenUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resumen)
}

// List godoc
// @Summary Listar vehículos de colección
// @Description Si se pasa ?etiqueta=nombre, filtra los vehículos por esa etiqueta
// @Tags coleccion
// @Security BearerAuth
// @Produce json
// @Param etiqueta query string false "Filtrar por nombre de etiqueta"
// @Success 200 {array} entities.Vehiculo
// @Router /coleccion [get]
func (ctrl *ListColeccionController) List(c *gin.Context) {
	etiquetaFiltro := c.Query("etiqueta")
	vehiculos, err := ctrl.listUC.Execute(c.Request.Context(), etiquetaFiltro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if vehiculos == nil {
		vehiculos = []entities.Vehiculo{}
	}
	c.JSON(http.StatusOK, vehiculos)
}

// GetByID godoc
// @Summary Obtener vehículo por ID
// @Tags coleccion
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID del vehículo"
// @Success 200 {object} entities.Vehiculo
// @Router /coleccion/{id} [get]
func (ctrl *ListColeccionController) GetByID(c *gin.Context) {
	id := c.Param("id")
	v, err := ctrl.getByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}
