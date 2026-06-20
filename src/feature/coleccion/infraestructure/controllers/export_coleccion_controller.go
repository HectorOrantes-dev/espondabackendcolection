package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"

	"coleccionbackend/src/feature/coleccion/application"
	"coleccionbackend/src/feature/coleccion/domain/entities"
)

type ExportColeccionController struct {
	listUC *application.ListColeccionUseCase
}

func NewExportColeccionController(l *application.ListColeccionUseCase) *ExportColeccionController {
	return &ExportColeccionController{listUC: l}
}

// Export godoc
// @Summary Exportar vehículos a Excel
// @Description Descarga un archivo .xlsx con los datos del formulario (nombre, marca, modelo)
// @Tags coleccion
// @Security BearerAuth
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file
// @Router /coleccion/export [get]
func (ctrl *ExportColeccionController) Handle(c *gin.Context) {
	vehiculos, err := ctrl.listUC.Execute(c.Request.Context(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	f := buildVehiculosExcel(vehiculos)
	defer f.Close()

	filename := fmt.Sprintf("coleccion_%s.xlsx", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generando el archivo Excel"})
		return
	}
}

// buildVehiculosExcel arma el archivo Excel con los datos del formulario.
func buildVehiculosExcel(vehiculos []entities.Vehiculo) *excelize.File {
	f := excelize.NewFile()

	const sheet = "Vehiculos"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"#", "Nombre", "Marca", "Modelo", "Etiquetas"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheet, "A1", "E1", headerStyle)

	for idx, v := range vehiculos {
		row := idx + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), idx+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Nombre)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Marca)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Modelo)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), etiquetasToString(v.Etiquetas))
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "D", 25)
	f.SetColWidth(sheet, "E", "E", 35)

	return f
}

// etiquetasToString une los nombres de las etiquetas separados por coma.
func etiquetasToString(etiquetas []entities.Etiqueta) string {
	nombres := make([]string, len(etiquetas))
	for i, e := range etiquetas {
		nombres[i] = e.Nombre
	}
	return strings.Join(nombres, ", ")
}
