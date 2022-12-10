package asciify

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/lupinelab/asciicam/utils"

	"github.com/gdamore/tcell/v2"
	"gocv.io/x/gocv"
)


// var ascii_symbols = []rune(string(".:-~=+*#%@"))
var ascii_symbols = []rune(string(".,;!vlLFE$"))

// var ascii_symbols = []rune(string(".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"))



func Asciify(cam *utils.Camera, frame *gocv.Mat, canvas tcell.Screen, settings *utils.Settings, prev_frame_time time.Time) {
	style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	convFrame := gocv.NewMat()
	term_width, term_height := canvas.Size()
	scale := math.Min(cam.Cap_width/float64(term_width), cam.Cap_height/float64(term_height))
	size := image.Point{X: int(cam.Cap_width / scale), Y: int(cam.Cap_height / (scale * 2.5))}
	gocv.Resize(*frame, &convFrame, size, 0, 0, gocv.InterpolationArea)
	if settings.SingleColourMode {
		style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.NewRGBColor(settings.Colour["R"], settings.Colour["G"], settings.Colour["B"]))
		gocv.CvtColor(convFrame, &convFrame, gocv.ColorBGRToGray)
		for y := 0; y < convFrame.Rows(); y++ {
			for x := 0; x < convFrame.Cols(); x++ {
				sym := ascii_symbols[int(float64(convFrame.GetUCharAt(y, x))/26)]
				canvas.SetContent(x, y, sym, nil, style)
			}
		}
	} else {
		greyFrame := gocv.NewMat()
		gocv.CvtColor(convFrame, &greyFrame, gocv.ColorBGRToGray)
		for y := 0; y < convFrame.Rows(); y++ {
			for x := 0; x < convFrame.Cols(); x++ {
				pixvec := convFrame.GetVecbAt(y, x)
				pixelcolour := tcell.NewRGBColor(int32(pixvec[2]), int32(pixvec[1]), int32(pixvec[0]))
				sym := ascii_symbols[int(float64(greyFrame.GetUCharAt(y, x))/26)]
				pixstyle := tcell.StyleDefault.
					Background(tcell.ColorReset).
					Foreground(pixelcolour)
				canvas.SetContent(x, y, sym, nil, pixstyle)
			}
		}
	}

	if settings.ShowHelp {
		for y, l := range utils.Help {
			for x, r := range l {
				canvas.SetContent(x, (size.Y -len(utils.Help)) + y, r, nil, style)
			}		
		}
	}

	if settings.ShowInfo {
		fps := int(1 / (time.Since(prev_frame_time).Seconds()))
		for i, r := range fmt.Sprintf("FPS=%v Brightness=%v Contrast=%v Colour=[R]%v[G]%v[B]%v ",
										fps, 
										settings.Brightness, 
										settings.Contrast, 
										settings.Colour["R"], 
										settings.Colour["G"], 
										settings.Colour["B"]) {
			canvas.SetContent(i, 0, r, nil, style)
		}
	}
}
