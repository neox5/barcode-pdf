package main

import (
	"log"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// pdf := gofpdf.New("P", "mm", "A4", "")
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "p",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: 50, Ht: 30},
	})
	pdf.SetMargins(5, 5, 5)
	pdf.SetCellMargin(0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetY(8)

	pdf.Rect(0, 0, 5, 5, "D")
	pdf.Rect(45, 0, 5, 5, "D")
	pdf.Rect(5, 5, 40, 20, "F")
	pdf.Rect(0, 25, 5, 5, "D")
	pdf.Rect(45, 25, 5, 5, "D")

	err := pdf.OutputFileAndClose("rect.pdf")
	if err != nil {
		log.Fatalf("error generating pdf: %v", err)
	}
}
