package asciify

import (
	"image"

	"github.com/lupinelab/asciicam/internal"

	"github.com/gdamore/tcell/v2"
	"gocv.io/x/gocv"
)

// var ascii_symbols = []rune(".:-~=+*#%@")
var ascii_symbols = []rune(".,;!vlLFE$")

func Asciify(frame *gocv.Mat, canvas tcell.Screen, settings *internal.Settings, termWidth int, termHeight int, scale float64, defStyle tcell.Style) (image.Point){
	pixStyle := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.NewRGBColor(settings.Colour["R"], settings.Colour["G"], settings.Colour["B"]))

	newSize := image.Point{X: int(settings.FrameWidth / scale), Y: int(settings.FrameHeight / (scale * 1.8))}

	downFrame := gocv.NewMat()
	gocv.Resize(*frame, &downFrame, newSize, 0, 0, gocv.InterpolationArea)
	
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

	return newSize
}
