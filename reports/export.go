package reports

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/taufikmulyawan/ticketing-system/types"
)

// GenerateSalesSummaryPDF creates a PDF report for overall sales summary
func GenerateSalesSummaryPDF(summary *types.SalesSummary) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	// Title
	pdf.Cell(40, 10, "Sales Summary Report")
	pdf.Ln(15)
	
	// Overall Summary
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Overall Summary")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Total Events: %d", summary.TotalEvents))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total Tickets Sold: %d", summary.TotalTickets))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total Revenue: Rp %.2f", summary.TotalRevenue))
	pdf.Ln(15)
	
	// Event Details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Event Details")
	pdf.Ln(10)
	
	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(10, 10, "ID")
	pdf.Cell(80, 10, "Event Name")
	pdf.Cell(30, 10, "Tickets Sold")
	pdf.Cell(30, 10, "Revenue")
	pdf.Ln(8)
	
	// Table content
	pdf.SetFont("Arial", "", 10)
	for _, event := range summary.EventSummary {
		pdf.Cell(10, 10, fmt.Sprintf("%d", event.EventID))
		pdf.Cell(80, 10, event.EventName)
		pdf.Cell(30, 10, fmt.Sprintf("%d", event.TotalTickets))
		pdf.Cell(30, 10, fmt.Sprintf("Rp %.2f", event.TotalRevenue))
		pdf.Ln(8)
	}
	
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04:05")))
	
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// GenerateEventSalesPDF creates a PDF report for a specific event
func GenerateEventSalesPDF(summary *types.EventSalesSummary) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	// Title
	pdf.Cell(40, 10, fmt.Sprintf("Event Sales Report: %s", summary.EventName))
	pdf.Ln(15)
	
	// Event Details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Event Details")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Event ID: %d", summary.EventID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Event Name: %s", summary.EventName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total Tickets Sold: %d", summary.TotalTickets))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total Revenue: Rp %.2f", summary.TotalRevenue))
	pdf.Ln(15)
	
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04:05")))
	
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// GenerateSalesSummaryCSV creates a CSV report for overall sales summary
func GenerateSalesSummaryCSV(summary *types.SalesSummary) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	
	// Write headers
	headers := []string{"Event ID", "Event Name", "Tickets Sold", "Revenue (Rp)"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	// Write data rows
	for _, event := range summary.EventSummary {
		row := []string{
			strconv.FormatUint(uint64(event.EventID), 10),
			event.EventName,
			strconv.FormatInt(event.TotalTickets, 10),
			fmt.Sprintf("%.2f", event.TotalRevenue),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	
	// Write summary row
	writer.Write([]string{"", "", "", ""})
	writer.Write([]string{
		"TOTAL",
		fmt.Sprintf("%d events", summary.TotalEvents),
		strconv.FormatInt(summary.TotalTickets, 10),
		fmt.Sprintf("%.2f", summary.TotalRevenue),
	})
	
	writer.Flush()
	
	if err := writer.Error(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// GenerateEventSalesCSV creates a CSV report for a specific event
func GenerateEventSalesCSV(summary *types.EventSalesSummary) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	
	// Write headers and data for the event
	writer.Write([]string{"Event ID", "Event Name", "Tickets Sold", "Revenue (Rp)"})
	writer.Write([]string{
		strconv.FormatUint(uint64(summary.EventID), 10),
		summary.EventName,
		strconv.FormatInt(summary.TotalTickets, 10),
		fmt.Sprintf("%.2f", summary.TotalRevenue),
	})
	
	writer.Flush()
	
	if err := writer.Error(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
} 