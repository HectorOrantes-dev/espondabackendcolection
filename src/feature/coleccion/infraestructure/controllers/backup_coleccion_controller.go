package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"coleccionbackend/src/feature/coleccion/application"
	"coleccionbackend/src/feature/coleccion/domain"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

// maxDescargasConcurrentes limita cuántas imágenes se bajan de Drive a la vez.
const maxDescargasConcurrentes = 5

type BackupColeccionController struct {
	listUC       *application.ListColeccionUseCase
	imageService domain.ImageService
}

func NewBackupColeccionController(l *application.ListColeccionUseCase, img domain.ImageService) *BackupColeccionController {
	return &BackupColeccionController{listUC: l, imageService: img}
}

// Backup godoc
// @Summary Descargar respaldo completo (ZIP)
// @Description Genera un .zip con el Excel de datos y todas las imágenes organizadas por vehículo
// @Tags coleccion
// @Security BearerAuth
// @Produce application/zip
// @Success 200 {file} file
// @Router /coleccion/backup [get]
func (ctrl *BackupColeccionController) Handle(c *gin.Context) {
	vehiculos, err := ctrl.listUC.Execute(c.Request.Context(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("respaldo_coleccion_%s.zip", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/zip")

	zw := zip.NewWriter(c.Writer)
	defer zw.Close()

	// 1. Agregar el Excel con los datos del formulario
	if err := ctrl.writeExcel(zw, vehiculos); err != nil {
		// El header ya se envió; solo podemos registrar el error en el ZIP parcial.
		return
	}

	// 2. Imágenes: se descargan en paralelo (red, lo lento) pero se escriben
	//    al ZIP de forma secuencial y en orden, para no corromper el archivo.

	// Construir la lista plana de tareas, preservando el orden.
	type imgTask struct {
		entryName string // ruta dentro del ZIP, ej: "imagenes/01_Batmobile/2.jpg"
		fileID    string
	}
	var tasks []imgTask
	for idx, v := range vehiculos {
		folder := fmt.Sprintf("imagenes/%02d_%s", idx+1, sanitizeName(v.Nombre))
		for i, id := range v.ImageIDs {
			tasks = append(tasks, imgTask{
				entryName: fmt.Sprintf("%s/%d.jpg", folder, i+1),
				fileID:    id,
			})
		}
	}

	// FASE 1: descargar todo en paralelo (máx maxDescargasConcurrentes a la vez).
	contenidos := make([][]byte, len(tasks))
	sem := make(chan struct{}, maxDescargasConcurrentes)
	var wg sync.WaitGroup
	for i, t := range tasks {
		wg.Add(1)
		sem <- struct{}{} // bloquea si ya hay 5 descargas en curso
		go func(i int, fileID string) {
			defer wg.Done()
			defer func() { <-sem }()
			if data, err := ctrl.imageService.Download(fileID); err == nil {
				contenidos[i] = data
			}
		}(i, t.fileID)
	}
	wg.Wait()

	// FASE 2: escribir al ZIP en orden (secuencial → estructura siempre correcta).
	for i, t := range tasks {
		if contenidos[i] == nil {
			continue // descarga fallida, se omite esa imagen
		}
		w, err := zw.Create(t.entryName)
		if err != nil {
			continue
		}
		_, _ = w.Write(contenidos[i])
	}
}

// writeExcel escribe el archivo Excel dentro del ZIP.
func (ctrl *BackupColeccionController) writeExcel(zw *zip.Writer, vehiculos []entities.Vehiculo) error {
	f := buildVehiculosExcel(vehiculos)
	defer f.Close()

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}

	w, err := zw.Create("coleccion.xlsx")
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

// sanitizeName limpia el nombre para usarlo como nombre de carpeta/archivo.
func sanitizeName(name string) string {
	replacer := strings.NewReplacer(
		"/", "-", "\\", "-", ":", "-", "*", "-",
		"?", "-", "\"", "-", "<", "-", ">", "-", "|", "-",
	)
	name = replacer.Replace(strings.TrimSpace(name))
	if name == "" {
		name = "vehiculo"
	}
	return name
}
