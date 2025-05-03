package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/service"
)

type ReportController interface {
	GetSalesReport(c *gin.Context)
	GetEventSalesReport(c *gin.Context)
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
// @Success 200 {object} service.SalesSummary
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
// @Success 200 {object} service.EventSalesSummary
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