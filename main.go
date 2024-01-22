package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/pdf417"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter barcode generation: ")
		input, _ := reader.ReadString('\n')
		// Trim all leading and trailing whitespace, including \r
		input = strings.TrimSpace(input)

		// delete all previous barcode files before generation
		err := deleteBarcodeFiles("./")
		if err != nil {
			log.Fatalln("error deleten previous barcode files:", err)
		}

		// Call your existing barcode PDF generation function
		generateBarcodePDF(input)
	}
}

func deleteBarcodeFiles(dir string) error {
	pattern := "barcode_*.pdf"

	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateBarcodePDF(barcodeData string) {
	// Check for invalid characters in barcodeData
	if strings.ContainsAny(barcodeData, `< > : " / \ | ? *`) {
		log.Fatalf("Invalid characters in barcode data: %s", barcodeData)
	}

	pdf := newPdfLabel()
	addBarcode(pdf, barcodeData)

	filename := "barcode_" + barcodeData + ".pdf"

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		log.Fatalf("error generating pdf: %v", err)
	}

	fmt.Printf("barcode_%s.pdf successfully generated\n", barcodeData)
	fmt.Println("---")
}

func newPdfLabel() *gofpdf.Fpdf {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "p",
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: 50, Ht: 30},
	})
	pdf.SetMargins(0, 0, 0)
	pdf.SetCellMargin(0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 10)
	pdf.SetY(10)

	return pdf
}

func addBarcode(pdf *gofpdf.Fpdf, barcodeData string) {
	pdf.RegisterImageOptionsReader("barcode.png", gofpdf.ImageOptions{ImageType: "png"}, createBarcodeReader(barcodeData))

	pdf.ImageOptions("barcode.png", 2.5, 0, 45, 0, true, gofpdf.ImageOptions{ImageType: "png"}, 0, "")
	// pdf.ImageOptions("barcode.png", 2.5, 0, 45, 0, true, gofpdf.ImageOptions{ImageType: "png"}, 0, "")
	pdf.SetY(pdf.GetY() + 1)
	pdf.CellFormat(50, 3, barcodeData, "0", 1, "C", false, 0, "")
}

func createBarcodeReader(data string) io.Reader {
	// bc, _ := code128.Encode(data)
	// bcScaled, err := barcode.Scale(bc, 600, 100)
	bc, _ := pdf417.Encode(data, 3)
	bcScaled, err := barcode.Scale(bc, 600, 100)
	if err != nil {
		log.Fatalf("error scaling barcode: %v", err)
	}

	// convert from 16bit to 8bit color depth for compatibility with pdf library
	bc8Bit := convertTo8Bit(bcScaled)

	var buf8Bit bytes.Buffer
	if err := png.Encode(&buf8Bit, bc8Bit); err != nil {
		panic(err)
	}

	return bytes.NewReader(buf8Bit.Bytes())
}

func convertTo8Bit(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	newImg := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			newImg.SetNRGBA(x, y, c)
		}
	}
	return newImg
}
