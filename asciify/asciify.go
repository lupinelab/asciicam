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

var ascii_symbols = []rune(string(".,-~:;=!*#$@"))

func Asciify(cam *utils.Camera, frame *gocv.Mat, canvas tcell.Screen, settings *utils.Settings) {
	style := tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.NewRGBColor(settings.Colour["R"], settings.Colour["G"], settings.Colour["B"]))
	canvas.Clear()
	cam.Cap.Set(gocv.VideoCaptureBrightness, settings.Brightness)
	cam.Cap.Set(gocv.VideoCaptureContrast, settings.Contrast)
	prev_frame_time := time.Now()
	greyFrame := gocv.NewMat()
	_ = greyFrame
	gocv.CvtColor(*frame, &greyFrame, gocv.ColorBGRToGray)
	downFrame := gocv.NewMat()
	term_width, term_height := canvas.Size()
	scale := math.Min(cam.Cap_width/float64(term_width), cam.Cap_height/float64(term_height))
	size := image.Point{X: int(cam.Cap_width / scale), Y: int(cam.Cap_height / (scale * 2.5))}
	gocv.Resize(greyFrame, &downFrame, size, 0, 0, gocv.InterpolationArea)
	for y := 0; y < downFrame.Rows(); y++ {
		for x := 0; x < downFrame.Cols(); x++ {
			sym := ascii_symbols[int(downFrame.GetUCharAt(y, x)/22)]
			canvas.SetContent(x, y, sym, nil, style)
		}
	}
	if settings.ShowStats {
		new_frame_time := time.Now()
		fps := int(1 / (time.Since(prev_frame_time).Seconds()))
		prev_frame_time = new_frame_time
		for i, r := range fmt.Sprintf("FPS=%v Brightness=%v Contrast=%v Colour=[R]%v[G]%v[B]%v ",
										fps, settings.Brightness, settings.Contrast, settings.Colour["R"], settings.Colour["G"], settings.Colour["B"]) {
			canvas.SetContent(i, 0, r, nil, style)
		}
	}
	if settings.ShowHelp {
		for y, l := range utils.Help {
			for x, r := range l {
				canvas.SetContent(x, (size.Y -len(utils.Help)) + y, r, nil, style)
		}
			
}
	}
	canvas.Show()
}
