package asciify

import (
	"image"
	"math"

	"github.com/lupinelab/asciicam/utils"

	"github.com/gdamore/tcell/v2"
	"gocv.io/x/gocv"
)

// var ascii_symbols = []rune(".:-~=+*#%@")
var ascii_symbols = []rune(".,;!vlLFE$")

func Asciify(cam *utils.Camera, frame *gocv.Mat, canvas tcell.Screen, settings *utils.Settings) {
	defstyle := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	
	pixstyle := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.NewRGBColor(settings.Colour["R"], settings.Colour["G"], settings.Colour["B"]))
	
	term_width, term_height := canvas.Size()
	scale := math.Min(cam.Cap_width/float64(term_width), cam.Cap_height/float64(term_height))
	size := image.Point{X: int(cam.Cap_width / scale), Y: int(cam.Cap_height / (scale * 1.8))}
	
	downframe := gocv.NewMat()
	gocv.Resize(*frame, &downframe, size, 0, 0, gocv.InterpolationArea)
	greyframe := gocv.NewMat()
	gocv.CvtColor(downframe, &greyframe, gocv.ColorBGRToGray)
	
	switch settings.SingleColourMode {
	case true:
		for y := 0; y < greyframe.Rows(); y++ {
			for x := 0; x < greyframe.Cols(); x++ {
				sym := ascii_symbols[int(float64(greyframe.GetUCharAt(y, x))/26)]
				canvas.SetContent(x, y, sym, nil, pixstyle)
			}
		}
	case false :
		for y := 0; y < greyframe.Rows(); y++ {
			for x := 0; x < greyframe.Cols(); x++ {
				pixvec := downframe.GetVecbAt(y, x)
				pixelcolour := tcell.NewRGBColor(int32(pixvec[2]), int32(pixvec[1]), int32(pixvec[0]))
				sym := ascii_symbols[int(float64(greyframe.GetUCharAt(y, x))/26)]
				pixstyle = tcell.StyleDefault.
					Background(tcell.ColorReset).
					Foreground(pixelcolour)
				canvas.SetContent(x, y, sym, nil, pixstyle)
			}
		}
	}

	if settings.ShowHelp {
		for y, l := range utils.Help {
			for x, r := range l {
				canvas.SetContent(x, (size.Y-len(utils.Help))+y, r, nil, defstyle)
			}
		}
	}
}
