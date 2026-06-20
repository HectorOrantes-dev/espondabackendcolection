package controllers

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/coleccion/application"
)

type UpdateColeccionController struct {
	useCase *application.UpdateColeccionUseCase
}

func NewUpdateColeccionController(uc *application.UpdateColeccionUseCase) *UpdateColeccionController {
	return &UpdateColeccionController{useCase: uc}
}

// Update godoc
// @Summary Actualizar vehículo de colección
// @Tags coleccion
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID del vehículo"
// @Param nombre formData string false "Nombre del vehículo"
// @Param marca formData string false "Marca del vehículo"
// @Param modelo formData string false "Modelo del vehículo"
// @Param imagenes formData file false "Nuevas imágenes (reemplazan las anteriores, máx 3)"
// @Success 200 {object} map[string]string
// @Router /coleccion/{id} [put]
func (ctrl *UpdateColeccionController) Handle(c *gin.Context) {
	id := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error leyendo formulario"})
		return
	}

	files := form.File["imagenes"]
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "máximo 3 imágenes permitidas"})
		return
	}

	var images []application.ImageInput
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error leyendo imagen"})
			return
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error leyendo contenido de imagen"})
			return
		}

		images = append(images, application.ImageInput{
			Filename: filepath.Base(fh.Filename),
			Content:  content,
			MimeType: fh.Header.Get("Content-Type"),
		})
	}

	input := application.UpdateColeccionInput{
		ID:     id,
		Nombre: c.PostForm("nombre"),
		Marca:  c.PostForm("marca"),
		Modelo: c.PostForm("modelo"),
		Images: images,
		// Si el campo "etiquetas" viene en el form, se reemplazan; si no, nil
		// y se conservan las actuales.
		EtiquetaIDs: form.Value["etiquetas"],
	}

	if err := ctrl.useCase.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehículo actualizado correctamente"})
}
