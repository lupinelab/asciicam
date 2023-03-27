package internal

import (
	"image"

	"github.com/gdamore/tcell/v2"
	"gocv.io/x/gocv"
)

// var ascii_symbols = []rune(".:-~=+*#%@")
var ascii_symbols = []rune(".,;!vlLFE$")

func Asciify(frame *gocv.Mat, canvas tcell.Screen, settings *Settings, termWidth int, termHeight int, scale float64, scaledResolution image.Point, defStyle tcell.Style) {
	pixStyle := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.NewRGBColor(settings.Colour["R"], settings.Colour["G"], settings.Colour["B"]))

	downFrame := gocv.NewMat()
	gocv.Resize(*frame, &downFrame, scaledResolution, 0, 0, gocv.InterpolationArea)
	
	greyFrame := gocv.NewMat()
	gocv.CvtColor(downFrame, &greyFrame, gocv.ColorBGRToGray)

	switch settings.SingleColourMode {
	case true:
		for y := 0; y < greyFrame.Rows(); y++ {
			for x := 0; x < greyFrame.Cols(); x++ {
				sym := ascii_symbols[int(float64(greyFrame.GetUCharAt(y, x))/26)]
				canvas.SetContent(x, y, sym, nil, pixStyle)
			}
		}
	case false:
		for y := 0; y < greyFrame.Rows(); y++ {
			for x := 0; x < greyFrame.Cols(); x++ {
				pixvec := downFrame.GetVecbAt(y, x)
				pixelcolour := tcell.NewRGBColor(int32(pixvec[2]), int32(pixvec[1]), int32(pixvec[0]))
				sym := ascii_symbols[int(float64(greyFrame.GetUCharAt(y, x))/26)]
				pixStyle = tcell.StyleDefault.
					Background(tcell.ColorReset).
					Foreground(pixelcolour)
				canvas.SetContent(x, y, sym, nil, pixStyle)
			}
		}
	}
}
