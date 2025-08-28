package orders

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)



func addHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Moroccan Beauty Store")
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(5)
	pdf.Cell(0, 10, "123 Business St, Las Vegas, NV 12345")
	pdf.Ln(5)
	pdf.Cell(0, 10, "Email: contact@moroccanbeautystore.com | Phone: (858) 456-7890")
	pdf.Ln(5)
}

func drawBorder(pdf *gofpdf.Fpdf, yPosition float64) {
	pdf.Ln(10)
	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(10, yPosition, 200, yPosition)

}

func addOrderTable(pdf *gofpdf.Fpdf, orSum *OrderSummary) {
	// table column widths
	colWidths := []float64{15, 100, 10, 30, 20}
	tableWidth := 0.0
	for _, width := range colWidths {
		tableWidth += width
	}

	//initial position
	pdf.SetXY(pdf.GetX(), pdf.GetY())

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)

	headers := []string{"Prd ID", "Name", "Qty", "Price", "Discount"}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "B", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	for _, product := range orSum.ProductList {
		pdf.CellFormat(colWidths[0], 6, product.Product_ID, "T", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[1], 6, product.Product_Name, "T", 0, "C", false, 0, "") // Replace with real product name if available
		pdf.CellFormat(colWidths[2], 6, product.Quantity, "T", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], 6, fmt.Sprintf("$%.2f", product.Price), "T", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[4], 6, product.Discount, "T", 1, "C", false, 0, "")
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

	pdf.SetDrawColor(0, 0, 0)
	pdf.Rect(boxX, boxY, boxWidth, boxHeight, "D")

	pdf.SetFont("Arial", "", 10)

	leftColumnY := boxY + 5
	pdf.SetXY(boxX+2, leftColumnY)
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