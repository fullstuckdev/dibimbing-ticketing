package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/service"
)

type ReportController interface {
	GetSalesReport(c *gin.Context)
	GetEventSalesReport(c *gin.Context)
	ExportSalesReportPDF(c *gin.Context)
	ExportEventSalesReportPDF(c *gin.Context)
	ExportSalesReportCSV(c *gin.Context)
	ExportEventSalesReportCSV(c *gin.Context)
}

type reportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) ReportController {
	return &reportController{
		reportService: reportService,
	}
}

// GetSalesReport godoc
// @Summary Get overall sales report
// @Description Get a summary of ticket sales and revenue across all events
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} types.SalesSummary
// @Failure 500 {object} map[string]interface{}
// @Router /reports/summary [get]
func (ctrl *reportController) GetSalesReport(c *gin.Context) {
	summary, err := ctrl.reportService.GetSalesSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetEventSalesReport godoc
// @Summary Get sales report for a specific event
// @Description Get a summary of ticket sales and revenue for a specific event
// @Tags reports
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Security BearerAuth
// @Success 200 {object} types.EventSalesSummary
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /reports/event/{id} [get]
func (ctrl *reportController) GetEventSalesReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	eventSummary, err := ctrl.reportService.GetEventSalesSummary(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, eventSummary)
}

// ExportSalesReportPDF godoc
// @Summary Export overall sales report as PDF
// @Description Export a PDF of ticket sales and revenue across all events
// @Tags reports
// @Accept json
// @Produce application/pdf
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 500 {object} map[string]interface{}
// @Router /reports/summary/pdf [get]
func (ctrl *reportController) ExportSalesReportPDF(c *gin.Context) {
	pdfBytes, err := ctrl.reportService.ExportSalesSummaryPDF()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for PDF download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=sales_report_"+time.Now().Format("2006-01-02")+".pdf")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportEventSalesReportPDF godoc
// @Summary Export event sales report as PDF
// @Description Export a PDF of ticket sales and revenue for a specific event
// @Tags reports
// @Accept json
// @Produce application/pdf
// @Param id path int true "Event ID"
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /reports/event/{id}/pdf [get]
func (ctrl *reportController) ExportEventSalesReportPDF(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	pdfBytes, err := ctrl.reportService.ExportEventSalesPDF(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for PDF download
	eventID := strconv.FormatUint(id, 10)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=event_"+eventID+"_report_"+time.Now().Format("2006-01-02")+".pdf")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ExportSalesReportCSV godoc
// @Summary Export overall sales report as CSV
// @Description Export a CSV of ticket sales and revenue across all events
// @Tags reports
// @Accept json
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 500 {object} map[string]interface{}
// @Router /reports/summary/csv [get]
func (ctrl *reportController) ExportSalesReportCSV(c *gin.Context) {
	csvBytes, err := ctrl.reportService.ExportSalesSummaryCSV()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for CSV download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=sales_report_"+time.Now().Format("2006-01-02")+".csv")
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.Data(http.StatusOK, "text/csv", csvBytes)
}

// ExportEventSalesReportCSV godoc
// @Summary Export event sales report as CSV
// @Description Export a CSV of ticket sales and revenue for a specific event
// @Tags reports
// @Accept json
// @Produce text/csv
// @Param id path int true "Event ID"
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /reports/event/{id}/csv [get]
func (ctrl *reportController) ExportEventSalesReportCSV(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	csvBytes, err := ctrl.reportService.ExportEventSalesCSV(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for CSV download
	eventID := strconv.FormatUint(id, 10)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=event_"+eventID+"_report_"+time.Now().Format("2006-01-02")+".csv")
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.Data(http.StatusOK, "text/csv", csvBytes)
} 