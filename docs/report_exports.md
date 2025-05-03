# Report Export Functionality

The ticketing system now supports exporting reports in PDF and CSV formats. This document describes the available export endpoints and their usage.

## Available Export Endpoints

### PDF Export Endpoints

1. **Export Overall Sales Report as PDF**

   - **URL**: `/reports/summary/pdf`
   - **Method**: `GET`
   - **Authentication**: Bearer Token (Admin only)
   - **Description**: Generates a PDF document with sales data across all events
   - **Response**: PDF file with comprehensive sales data

2. **Export Event-Specific Sales Report as PDF**
   - **URL**: `/reports/event/{id}/pdf`
   - **Method**: `GET`
   - **Authentication**: Bearer Token (Admin only)
   - **URL Parameters**: `id` - The ID of the event
   - **Description**: Generates a PDF document with sales data for a specific event
   - **Response**: PDF file with event-specific sales data

### CSV Export Endpoints

1. **Export Overall Sales Report as CSV**

   - **URL**: `/reports/summary/csv`
   - **Method**: `GET`
   - **Authentication**: Bearer Token (Admin only)
   - **Description**: Generates a CSV file with sales data across all events
   - **Response**: CSV file with tabular sales data

2. **Export Event-Specific Sales Report as CSV**
   - **URL**: `/reports/event/{id}/csv`
   - **Method**: `GET`
   - **Authentication**: Bearer Token (Admin only)
   - **URL Parameters**: `id` - The ID of the event
   - **Description**: Generates a CSV file with sales data for a specific event
   - **Response**: CSV file with tabular event-specific sales data

## File Format Details

### PDF Reports

The PDF reports include:

- Title with report name and date
- Summary section with key statistics
- Detailed data in tabular format
- Footer with generation timestamp

### CSV Reports

The CSV reports include:

- Header row with column names
- Data rows for each event or ticket
- Summary row at the bottom with totals

## Usage Examples

### Using with Postman

1. Select the desired export endpoint from the Reports folder
2. Ensure you have a valid admin token in the `admin_token` variable
3. Send the request and save the response as a file

### Using with cURL

```bash
# Export overall sales report as PDF
curl -X GET "http://localhost:8080/reports/summary/pdf" \
     -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
     --output sales_report.pdf

# Export event sales report as CSV
curl -X GET "http://localhost:8080/reports/event/1/csv" \
     -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
     --output event_1_report.csv
```

## Implementation Details

The export functionality is implemented using:

- PDF Generation: `github.com/jung-kurt/gofpdf` package
- CSV Generation: Go's built-in `encoding/csv` package
- Types: Custom types defined in the `types` package
- Controllers: Export endpoints in `controller/report_controller.go`
- Services: Export methods in `service/report_service.go`
- Generators: Export generation in `reports/export.go`
