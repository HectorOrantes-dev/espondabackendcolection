package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/coleccion/application"
)

type DeleteColeccionController struct {
	useCase *application.DeleteColeccionUseCase
}

func NewDeleteColeccionController(uc *application.DeleteColeccionUseCase) *DeleteColeccionController {
	return &DeleteColeccionController{useCase: uc}
}

// Delete godoc
// @Summary Eliminar vehículo de colección
// @Tags coleccion
// @Security BearerAuth
// @Param id path string true "ID del vehículo"
// @Success 204
// @Router /coleccion/{id} [delete]
func (ctrl *DeleteColeccionController) Handle(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.useCase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
