package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/etiquetas/application"
	"coleccionbackend/src/feature/etiquetas/domain/entities"
)

type EtiquetasController struct {
	createUC *application.CreateEtiquetaUseCase
	listUC   *application.ListEtiquetasUseCase
	updateUC *application.UpdateEtiquetaUseCase
	deleteUC *application.DeleteEtiquetaUseCase
}

func NewEtiquetasController(
	c *application.CreateEtiquetaUseCase,
	l *application.ListEtiquetasUseCase,
	u *application.UpdateEtiquetaUseCase,
	d *application.DeleteEtiquetaUseCase,
) *EtiquetasController {
	return &EtiquetasController{createUC: c, listUC: l, updateUC: u, deleteUC: d}
}

type etiquetaRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

// Create godoc
// @Summary Crear etiqueta
// @Tags etiquetas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body etiquetaRequest true "Nombre de la etiqueta"
// @Success 201 {object} entities.Etiqueta
// @Router /etiquetas [post]
func (ctrl *EtiquetasController) Create(c *gin.Context) {
	var req etiquetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e, err := ctrl.createUC.Execute(c.Request.Context(), req.Nombre)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

// List godoc
// @Summary Listar etiquetas con cantidad de vehículos
// @Tags etiquetas
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.Etiqueta
// @Router /etiquetas [get]
func (ctrl *EtiquetasController) List(c *gin.Context) {
	etiquetas, err := ctrl.listUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if etiquetas == nil {
		etiquetas = []entities.Etiqueta{}
	}
	c.JSON(http.StatusOK, etiquetas)
}

// Update godoc
// @Summary Editar etiqueta
// @Tags etiquetas
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID de la etiqueta"
// @Param body body etiquetaRequest true "Nuevo nombre"
// @Success 200 {object} map[string]string
// @Router /etiquetas/{id} [put]
func (ctrl *EtiquetasController) Update(c *gin.Context) {
	id := c.Param("id")
	var req etiquetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.updateUC.Execute(c.Request.Context(), id, req.Nombre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "etiqueta actualizada"})
}

// Delete godoc
// @Summary Eliminar etiqueta
// @Tags etiquetas
// @Security BearerAuth
// @Param id path string true "ID de la etiqueta"
// @Success 204
// @Router /etiquetas/{id} [delete]
func (ctrl *EtiquetasController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.deleteUC.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
