package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/phpdave11/gofpdf"
	"strconv"
	"strings"
)

type PageStyle struct {
	ML                    float64 // margen izq
	MT                    float64 // margen sup
	MR                    float64 // margen der
	MB                    float64 // margen inf
	WW                    float64 // ancho area trabajo
	HW                    float64 // alto area trabajo
	HH                    float64 // alto header
	HB                    float64 // alto body
	HF                    float64 // alto footer
	BaseColorRGB          [3]int
	SecondaryColorRGB     [3]int
	ComplementaryColorRGB [3]int
}

// EncodePDF convierte pdf a base64
func EncodePDF(pdf *gofpdf.Fpdf) string {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("plan.pdf") // para guardar el archivo localmente
	pdf.Output(writer)
	writer.Flush()
	encodedFile := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encodedFile
}

// AddImage agrega imagen de archivo a pdf, w o h en cero autoajusta segun ratio imagen
func AddImage(pdf *gofpdf.Fpdf, image string, x, y, w, h float64) *gofpdf.Fpdf {
	//The ImageOptions method takes a file path, x, y, width, and height parameters, and an ImageOptions struct to specify a couple of options.
	imageSplit := strings.Split(image, ".")
	imageType := imageSplit[len(imageSplit)-1]

	pdf.ImageOptions(image, x, y, w, h, false, gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
	return pdf
}

func FontStyle(pdf *gofpdf.Fpdf, style string, size float64, bw int, fontFamily string) {
	pdf.SetTextColor(bw, bw, bw)

	if fontFamily == "" {
		fontFamily = "Arial"
	}
	pdf.SetFont(fontFamily, style, size)
}

func Hex2RGB(hexString string) ([3]int, error) {
	hexString = strings.Replace(hexString, "#", "", 1)

	var rgb [3]int
	values, err := strconv.ParseInt(hexString, 16, 32)

	if err != nil {
		return [3]int{}, err
	}

	rgb[0] = int(values >> 16)
	rgb[1] = int((values >> 8) & 0xFF)
	rgb[2] = int(values & 0xFF)

	return rgb, err
}
