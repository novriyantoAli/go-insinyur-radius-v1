package helper

import (
	"net/http"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/signintech/gopdf"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type ResponseErrorMessage struct {
	Message string `json:"error"`
}

type ResponseSuccessMessage struct {
	Message string `json:"message"`
}

// GroupByProperty groups a slice of structs by a specific property.
func GroupByProperty[T any, K comparable](items []T, getProperty func(T) K) map[K][]T {
	grouped := make(map[K][]T)
	for _, item := range items {
		key := getProperty(item)
		grouped[key] = append(grouped[key], item)
	}
	return grouped
}

// TranslateError ...
func TranslateError(err error) int {
	switch err {
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func rectFillColor(pdf *gopdf.GoPdf, text string, fontSize int, x, y, w, h float64, r, g, b uint8) (err error) {
	pdf.SetLineWidth(2.1)
	pdf.SetFillColor(r, g, b) //setup fill color
	pdf.RectFromUpperLeftWithStyle(x, y, w, h, "FD")
	pdf.SetFillColor(0, 0, 0)

	qrc, err := qrcode.New(text)
	if err != nil {
		return
	}

	w0, err := standard.New("./res/imgs/generates/"+text+".png",
		standard.WithHalftone("./res/imgs/pngegg.png"),
		standard.WithQRWidth(21),
	)
	if err != nil {
		return
	}

	err = qrc.Save(w0)
	if err != nil {
		return
	}

	pdf.Image("./res/imgs/generates/"+text+".png", x, y, &gopdf.Rect{
		W: w,
		H: (h - (float64(fontSize) * 2)),
	})

	// originalX := x
	textw, _ := pdf.MeasureTextWidth(text)
	x = x + (w / 2) - (textw / 2)
	// set y
	y = y + h - (float64(fontSize) * 2)
	pdf.SetXY(x, y)
	pdf.Cell(nil, text)

	return
}

func PrintVocuhers(code string, data []domain.Vcr) (err error) {

	pageWidth := 595.28
	pageHeight := 841.89

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pageWidth, H: pageHeight}}) //595.28, 841.89 = A4
	pdf.AddPage()

	err = pdf.AddTTFFont("Ubuntu-L", "./res/fonts/Ubuntu-L.ttf")
	if err != nil {
		return
	}
	fontSize := 18
	err = pdf.SetFont("Ubuntu-L", "U", fontSize)
	if err != nil {
		return
	}

	rectWidth := 100.0
	rectHeight := 100.0

	xStart := 10.0
	yStart := 10.0

	for i := 0; i < len(data); i++ {
		if (xStart + rectWidth) < pageWidth {
			rectFillColor(&pdf, data[i].Username, fontSize, xStart, yStart, rectWidth, rectHeight, 255, 255, 255)
			xStart += (rectWidth + 10)
		} else {
			if (yStart + rectHeight) > pageHeight {
				pdf.AddPage()
				yStart = 10
			} else {
				yStart += (rectHeight + 10)
			}
			xStart = 10
			rectFillColor(&pdf, data[i].Username, fontSize, xStart, yStart, rectWidth, rectHeight, 255, 255, 255)
			xStart += (rectWidth + 10)
		}
	}

	pdf.WritePdf("./res/pdfs/" + code + ".pdf")

	return
}
