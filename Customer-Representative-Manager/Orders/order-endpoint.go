package orders

import (
	"fmt"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func OrderHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("hey!")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// Add Header
	addHeader(pdf)

	// Draw Border Below Header
	drawBorder(pdf, 30)

	// Add Order Table
	addOrderTable(pdf)

	addOrderSummary(pdf)
	// Add Footer
	addFooter(pdf)

	// Save PDF
	err := pdf.OutputFileAndClose("order_summary.pdf")
	if err != nil {
		panic(err)
	}
}

func addHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Souhaila and Alex")
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(5)
	pdf.Cell(0, 10, "123 Business St, Business City, BC 12345")
	pdf.Ln(5)
	pdf.Cell(0, 10, "Email: contact@company.com | Phone: (123) 456-7890")
	pdf.Ln(5)
}

func drawBorder(pdf *gofpdf.Fpdf, yPosition float64) {
	pdf.Ln(10)
	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(10, yPosition, 200, yPosition)
	
}

func addOrderTable(pdf *gofpdf.Fpdf) {
	// Add Table Header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(40, 10, "Order ID", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Discount", "1", 1, "C", true, 0, "")

	// Add Table Rows
	pdf.SetFont("Arial", "", 12)
	orderData := []struct {
		ID       string
		Quantity int
		Price    float64
		Discount float64
	}{
		{"12345", 2, 49.99, 5.00},
		{"12346", 1, 19.99, 2.00},
		{"12347", 5, 99.95, 10.00},
	}

	for _, order := range orderData {
		pdf.CellFormat(40, 10, order.ID, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, "45", "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, "$55.55", "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, "25%", "1", 1, "C", false, 0, "")
	}

	pdf.Ln(10)
}

func addFooter(pdf *gofpdf.Fpdf) {
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 12)
	pdf.Cell(10, 10, "Thank you for your business!")
	pdf.Ln(5)
	pdf.Cell(10, 10, "If you have any questions, please contact us.")
}


func addOrderSummary(pdf *gofpdf.Fpdf) {
	boxX := 10.0
    boxY := pdf.GetY()
    boxWidth := 190.0
    boxHeight := 90.0
    columnGap := 5.0
    leftWidth := (boxWidth - columnGap) / 2
    rightX := boxX + leftWidth + columnGap

    // Draw the outer box
    pdf.SetDrawColor(0, 0, 0)
    pdf.Rect(boxX, boxY, boxWidth, boxHeight, "D")

    // Set smaller font
    pdf.SetFont("Arial", "", 10)

    // Left Column: Customer & Order Details
    leftColumnY := boxY + 5
    pdf.SetXY(boxX + 2, leftColumnY)
    pdf.CellFormat(leftWidth-4, 5, "Shipping Address:", "B", 1, "L", false, 0, "")
    pdf.MultiCell(leftWidth-4, 5, "John Doe\n123 Elm Street\nApt 456\nBusiness City, BC 12345", "", "L", false)
    pdf.Ln(2)

    pdf.SetX(boxX + 2)
    pdf.CellFormat(leftWidth-4, 5, "Customer Name:", "B", 1, "L", false, 0, "")
    pdf.CellFormat(leftWidth-4, 5, "John Doe", "", 1, "L", false, 0, "")
    pdf.Ln(2)

    pdf.SetX(boxX + 2)
    pdf.CellFormat(leftWidth-4, 5, "Payment Method:", "B", 1, "L", false, 0, "")
    pdf.CellFormat(leftWidth-4, 5, "Credit Card (Visa ****1234)", "", 1, "L", false, 0, "")
    pdf.Ln(2)

    pdf.SetX(boxX + 2)
    pdf.CellFormat(leftWidth-4, 5, "Order Name:", "B", 1, "L", false, 0, "")
    pdf.CellFormat(leftWidth-4, 5, "Order #12345", "", 1, "L", false, 0, "")
    pdf.Ln(2)

    pdf.SetX(boxX + 2)
    pdf.CellFormat(leftWidth-4, 5, "Order Date & Time:", "B", 1, "L", false, 0, "")
    pdf.CellFormat(leftWidth-4, 5, "2024-11-16 14:35:00", "", 1, "L", false, 0, "")

    // Right Column: Order Totals
    pdf.SetXY(rightX, leftColumnY)
    pdf.CellFormat(leftWidth-4, 5, "Total:", "B", 1, "R", false, 0, "")
    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "$199.99", "", 1, "R", false, 0, "")

    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "Taxes:", "B", 1, "R", false, 0, "")
    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "$15.00", "", 1, "R", false, 0, "")

    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "Shipping Method:", "B", 1, "R", false, 0, "")
    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "Standard Shipping", "", 1, "R", false, 0, "")

    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "Shipping Cost:", "B", 1, "R", false, 0, "")
    pdf.SetX(rightX)
    pdf.CellFormat(leftWidth-4, 5, "$5.99", "", 1, "R", false, 0, "")
}

// // Helper functions for formatting
// func formatInt(num int) string {
// 	return gofpdf.IntToStr(num)
// }

// func formatFloat(num float64) string {
// 	return gofpdf.Sprintf("$%.2f", num)
// }

