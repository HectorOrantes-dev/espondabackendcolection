package controllers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/coleccion/application"
)

type CreateColeccionController struct {
	useCase *application.CreateColeccionUseCase
}

func NewCreateColeccionController(uc *application.CreateColeccionUseCase) *CreateColeccionController {
	return &CreateColeccionController{useCase: uc}
}

// Create godoc
// @Summary Crear vehículo de colección
// @Tags coleccion
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param nombre formData string true "Nombre del vehículo"
// @Param marca formData string true "Marca del vehículo"
// @Param modelo formData string true "Modelo del vehículo"
// @Param imagenes formData file false "Imagen 1 (máx 3)"
// @Success 201 {object} entities.Vehiculo
// @Router /coleccion [post]
func (ctrl *CreateColeccionController) Handle(c *gin.Context) {
	nombre := c.PostForm("nombre")
	marca := c.PostForm("marca")
	modelo := c.PostForm("modelo")

	if nombre == "" || marca == "" || modelo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nombre, marca y modelo son requeridos"})
		return
	}

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

	// El precio es opcional; si no viene o es inválido, queda en 0.
	precio, _ := strconv.ParseFloat(c.PostForm("precio"), 64)

	input := application.CreateColeccionInput{
		Nombre:      nombre,
		Marca:       marca,
		Modelo:      modelo,
		Precio:      precio,
		Images:      images,
		EtiquetaIDs: form.Value["etiquetas"], // IDs de etiquetas (campo repetible)
	}

	v, err := ctrl.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, v)
}
