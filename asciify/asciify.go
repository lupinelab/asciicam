package asciify

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/blackjack/webcam"
	"github.com/gdamore/tcell/v2"
	"gocv.io/x/gocv"
)

var colour tcell.Color
var ascii_symbols = []rune(string(".,-~:;=!*#$@"))

type Settings struct {
	Brightness_caps map[string]float64
	Contrast_caps   map[string]float64
	Brightness      float64
	Contrast        float64
	Supported_resolutions	[]string
}

type Camera struct {
	Cap        *gocv.VideoCapture
	Cap_width  float64
	Cap_height float64
}

func Newcam(device string) (*Camera, *Settings, error) {
	s := Settings{}

	cam_caps, err := webcam.Open(device)
	if err != nil {
		return nil, nil, err
	}
	defer cam_caps.Close()

	capmap := cam_caps.GetControls()

	s.Brightness_caps = make(map[string]float64)
	s.Brightness_caps["min"] = float64(capmap[webcam.ControlID(0x00980900)].Min)
	s.Brightness_caps["max"] = float64(capmap[webcam.ControlID(0x00980900)].Max)
	brightness, err := cam_caps.GetControl(webcam.ControlID(0x00980900))
	if err != nil {
		return nil, nil, err
	}
	s.Brightness = float64(brightness)

	s.Contrast_caps = make(map[string]float64)
	s.Contrast_caps["min"] = float64(capmap[webcam.ControlID(0x00980901)].Min)
	s.Contrast_caps["max"] = float64(capmap[webcam.ControlID(0x00980901)].Max)
	contrast, err := cam_caps.GetControl(webcam.ControlID(0x00980901))
	if err != nil {
		return nil, nil, err
	}
	s.Contrast = float64(contrast)

	s.Supported_resolutions = []string{}
	resolutions := cam_caps.GetSupportedFrameSizes(webcam.PixelFormat(1196444237))
	for _, fs := range resolutions {
		s.Supported_resolutions = append(s.Supported_resolutions, fs.GetString())
	}

	cam := Camera{}

	cam.Cap, err = gocv.OpenVideoCaptureWithAPI(device, 200)
	if err != nil {
		return nil, nil, err
	}
	cam.Cap_width = cam.Cap.Get(gocv.VideoCaptureFrameWidth)
	cam.Cap_height = cam.Cap.Get(gocv.VideoCaptureFrameHeight)

	return &cam, &s, err
}

func CamToAscii(cam *Camera, frame *gocv.Mat, canvas tcell.Screen, settings *Settings) {
	colour = tcell.NewRGBColor(18, 181, 131)
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(colour)
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
	new_frame_time := time.Now()
	fps := int(1 / (time.Since(prev_frame_time).Seconds()))
	prev_frame_time = new_frame_time
	for i, r := range fmt.Sprintf("%vFPS, Brightness=%v, Contrast=%v", fps, settings.Brightness, settings.Contrast) {
		canvas.SetContent(i, 0, r, nil, style)
	}
	canvas.Show()
}

// func AsciiFy(cam *Camera, settings *Settings) {
// 	err := termbox.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer termbox.Close()
// 	termbox.SetOutputMode(termbox.OutputRGB)
// 	event_queue := make(chan termbox.Event, 1)
// 	go func() {
// 		for {
// 			event_queue <- termbox.PollEvent()
// 		}
// 	}()
// mainloop:
// 	for {
// 		select {
// 		case control := <-event_queue:
// 			if control.Type == termbox.EventKey && control.Key == termbox.KeyEsc || control.Key == termbox.KeyCtrlC {
// 				break mainloop
// 			}
// 			if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowUp {
// 				if settings.Brightness < settings.Brightness_caps["max"] {
// 					settings.Brightness += 1
// 				}
// 			}
// 			if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowDown {
// 				if settings.Brightness > settings.Brightness_caps["min"] {
// 					settings.Brightness -= 1
// 				}
// 			}
// 			if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowRight {
// 				if settings.Contrast < settings.Contrast_caps["max"] {
// 					settings.Contrast += 1
// 				}
// 			}
// 			if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowLeft {
// 				if settings.Contrast > settings.Contrast_caps["min"] {
// 					settings.Contrast -= 1
// 				}
// 			}
// 		default:
// 			frame := gocv.NewMat()
// 			cam.Cap.Read(&frame)
// 			camtoascii(cam, &frame, settings)
// 		}
// 	}
// }
