package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lupinelab/asciicam/asciify"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var asciicamCmd = &cobra.Command{
	Use:   "asciicam [device]",
	Short: "Turn your camera into ASCII",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO add regex to match expected video device path
		// Get camera and capabilities
		cam, settings, err := asciify.Newcam(args[0])
		if err != nil {
			panic(err.Error())
		}
		defer cam.Cap.Close()
		
		canvas, err := tcell.NewScreen()
		if err != nil {
			panic(err)
		}
		err = canvas.Init()
		if err != nil {
			panic(err)
		}
		defer canvas.Fini()
		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		canvas.SetStyle(defStyle)
		// err = termbox.Init()
		// if err != nil {
		// 	panic(err)
		// }
		// defer termbox.Close()
		// termbox.SetOutputMode(termbox.OutputRGB)
		res_sel_screen := []string{}
		res_sel_screen = append(res_sel_screen, fmt.Sprintf("Press a number key to choose a resolution:"))
		for i, fs := range settings.Supported_resolutions {
			res_sel_screen = append(res_sel_screen, fmt.Sprintf("%v) %v", i, fs))
		}
		for i, l := range res_sel_screen {
			for n, r := range l {
				canvas.SetContent(n, i, r, nil, defStyle)
			}
		}
		canvas.Show()
		var resolution string
	inputloop:
		for {
			input := canvas.PollEvent()
			switch input := input.(type) {
			case *tcell.EventKey:
				for i, fs := range settings.Supported_resolutions {
					value, err := strconv.Atoi(string(input.Rune()))
					if err != nil {
						panic(err)
					}
					if value == i {
						resolution = fs
						break inputloop
					}
			// input := termbox.PollEvent()
			// switch input.Type {
			// case termbox.EventKey:
			// 	for i, fs := range settings.Supported_resolutions {
			// 		if int(input.Key) == i {
			// 			resolution = fs
			// 			break inputloop
			// 		}
			// 	}
			// }
				}
			}
		}
		fWH := strings.Split(resolution, "x")
		fW, err := strconv.ParseFloat(fWH[0], 32)
		if err != nil {
			panic(err)
		}
		fH, err :=  strconv.ParseFloat(fWH[1], 32)
		if err != nil {
			panic(err)
		}
		cam.Cap.Set(gocv.VideoCaptureFrameWidth, fW)
		cam.Cap.Set(gocv.VideoCaptureFrameHeight, fH)
		// event_queue := make(chan tcell.Event, 1)
		// go func() {
		// 	for {
		// 		event_queue <- canvas.PollEvent()
		// 	}
		// }()
		quit := make(chan struct{})
		go func () {
			for {
				control := canvas.PollEvent()
				// select {
				// case control := <-event_queue:
				switch control := control.(type) {
				case *tcell.EventKey:
					if control.Key() == tcell.KeyEsc || control.Key() == tcell.KeyCtrlC {
						close(quit)
						return
					}
					// if control.Type == tcell.EventKey && control.Key == termbox.KeyEsc || control.Key == termbox.KeyCtrlC {
					// 	break mainloop
					// }
					if control.Key() == tcell.KeyUp {
						if settings.Brightness < settings.Brightness_caps["max"] {
							settings.Brightness += 1
						}
					}
					// if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowUp {
					// 	if settings.Brightness < settings.Brightness_caps["max"] {
					// 		settings.Brightness += 1
					// 	}
					// }
					if control.Key() == tcell.KeyDown {
						if settings.Brightness > settings.Brightness_caps["min"] {
							settings.Brightness -= 1
						}
					}
					// if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowDown {
					// 	if settings.Brightness > settings.Brightness_caps["min"] {
					// 		settings.Brightness -= 1
					// 	}
					// }
					if control.Key() == tcell.KeyRight {
						if settings.Contrast < settings.Contrast_caps["max"] {
							settings.Contrast += 1
						}
					}
					// if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowRight {
					// 	if settings.Contrast < settings.Contrast_caps["max"] {
					// 		settings.Contrast += 1
					// 	}
					// }
					if control.Key() == tcell.KeyLeft {
						if settings.Contrast > settings.Contrast_caps["min"] {
							settings.Contrast -= 1
						}
					}
					// if control.Type == termbox.EventKey && control.Key == termbox.KeyArrowLeft {
					// 	if settings.Contrast > settings.Contrast_caps["min"] {
					// 		settings.Contrast -= 1
					// 	}
					// }
				// default:
				}
			}
		}()
	mainloop:
		for {
			select {
			case <-quit:
				break mainloop
			default:
				frame := gocv.NewMat()
				cam.Cap.Read(&frame)
				asciify.CamToAscii(cam, &frame, canvas, settings)
			}
		}
	},
}

func Execute() error {
	return asciicamCmd.Execute()
}

func init() {
	asciicamCmd.PersistentFlags().BoolP("list-res", "l", false, "List available resolutions of a device")
	asciicamCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	asciicamCmd.MarkFlagRequired("")
	asciicamCmd.PersistentFlags().Lookup("help").Hidden = true
	cobra.EnableCommandSorting = false
	asciicamCmd.CompletionOptions.DisableDefaultCmd = true
}
