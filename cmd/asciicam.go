package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lupinelab/asciicam/asciify"
	"github.com/lupinelab/asciicam/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var asciicamCmd = &cobra.Command{
	Use:   "asciicam [device]",
	Short: "Turn your camera into ASCII",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get camera and capabilities
		settings, err := utils.NewSettings(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Create a tcell screen to use as a canvas
		canvas, err := tcell.NewScreen()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = canvas.Init()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer canvas.Fini()

		// Default terminal style
		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		canvas.SetStyle(defStyle)

		// Write a slice of strings onto the canvas as the message select screen
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

		// Wait for the user to choose a resolution or exit
		var resolution string
	inputloop:
		for {
			input := canvas.PollEvent()
			switch input := input.(type) {
			case *tcell.EventKey:
				if input.Key() == tcell.KeyEsc || input.Key() == tcell.KeyCtrlC {
					canvas.Fini()
					os.Exit(0)
				}
				for i, fs := range settings.Supported_resolutions {
					// Is there a better way to see the human value of a key press Rune?
					value, err := strconv.Atoi(string(input.Rune()))
					if err != nil {
						goto inputloop
					}
					// Compare the converted rune to the numbered list items
					if value == i {
						resolution = fs
						break inputloop
					}
				}
				goto inputloop
			}
		}
		
		// Parse the selected resolution for the height and width
		fWH := strings.Split(resolution, "x")

		// Frame height
		fW, err := strconv.ParseFloat(fWH[0], 32)
		if err != nil {
			canvas.Fini()
			fmt.Println(err)
		}
		// Frame width
		fH, err := strconv.ParseFloat(fWH[1], 32)
		if err != nil {
			canvas.Fini()
			fmt.Println(err)
		}

		// Clear the canvas and show the controls
		canvas.Clear()
		ready_screen := []string{}
		ready_screen = append(ready_screen, fmt.Sprint(""))
		ready_screen = append(ready_screen, fmt.Sprint("Controls"))
		ready_screen = append(ready_screen, fmt.Sprint("--------------------------"))
		for _, l := range utils.Help {
			ready_screen = append(ready_screen, l)
		}
		ready_screen = append(ready_screen, fmt.Sprint("--------------------------"))
		ready_screen = append(ready_screen, fmt.Sprint(""))
		ready_screen = append(ready_screen, fmt.Sprint("Press Enter key when ready..."))
		for i, l := range ready_screen {
			for n, r := range l {
				canvas.SetContent(n, i, r, nil, defStyle)
			}
		}
		canvas.Show()

		// wait for user to kick things off
	ready:
		for {
			input := canvas.PollEvent()
			switch input := input.(type) {
			case *tcell.EventKey:
				if input.Key() == tcell.KeyEsc || input.Key() == tcell.KeyCtrlC {
					canvas.Fini()
					os.Exit(0)
				} else if input.Key() == tcell.KeyEnter {
					break ready
				}
			}
		}

		// The camera
		cam, err := utils.NewCamera(args[0])
		if err != nil {
			canvas.Fini()
			fmt.Println(err)
		}
		defer cam.Cap.Close()
		
		cam.Cap.Set(gocv.VideoCaptureFrameWidth, fW)
		cam.Cap.Set(gocv.VideoCaptureFrameHeight, fH)

		// Listen for control keypresses (non blocking)
		quit := make(chan struct{})
		go func() {
			for {
				control := canvas.PollEvent()
				switch control := control.(type) {
				case *tcell.EventKey:
					if control.Key() == tcell.KeyEsc || control.Key() == tcell.KeyCtrlC {
						close(quit)
						return
					}
					// Brightness controls
					if control.Key() == tcell.KeyUp {
						if settings.Brightness < settings.Brightness_caps["max"] {
							settings.Brightness += 1
							cam.Cap.Set(gocv.VideoCaptureBrightness, settings.Brightness)
						}
					}
					if control.Key() == tcell.KeyDown {
						if settings.Brightness > settings.Brightness_caps["min"] {
							settings.Brightness -= 1
							cam.Cap.Set(gocv.VideoCaptureBrightness, settings.Brightness)
						}
					}
					// Contrast controls
					if control.Key() == tcell.KeyRight {
						if settings.Contrast < settings.Contrast_caps["max"] {
							settings.Contrast += 1
							cam.Cap.Set(gocv.VideoCaptureContrast, settings.Contrast)
						}
					}
					if control.Key() == tcell.KeyLeft {
						if settings.Contrast > settings.Contrast_caps["min"] {
							settings.Contrast -= 1
							cam.Cap.Set(gocv.VideoCaptureContrast, settings.Contrast)
						}
					}
					// SingleColourMode control
					if string(control.Rune()) == "m" {
						if settings.SingleColourMode == true {
							settings.SingleColourMode = false
						} else if settings.SingleColourMode == false {
							settings.SingleColourMode = true
						}
					}
					// Colour Controls
					if string(control.Rune()) == "r" {
						if settings.Colour["R"] < 255 {
							settings.Colour["R"] += 1
						}
					}
					if string(control.Rune()) == "e" {
						if settings.Colour["R"] > 0 {
							settings.Colour["R"] -= 1
						}
					}
					if string(control.Rune()) == "g" {
						if settings.Colour["G"] < 255 {
							settings.Colour["G"] += 1
						}
					}
					if string(control.Rune()) == "f" {
						if settings.Colour["G"] > 0 {
							settings.Colour["G"] -= 1
						}
					}
					if string(control.Rune()) == "b" {
						if settings.Colour["B"] < 255 {
							settings.Colour["B"] += 1
						}
					}
					if string(control.Rune()) == "v" {
						if settings.Colour["B"] > 0 {
							settings.Colour["B"] -= 1
						}
					}
					// ShowInfo control
					if string(control.Rune()) == "i" {
						if settings.ShowInfo == true {
							settings.ShowInfo = false
						} else if settings.ShowInfo == false {
							settings.ShowInfo = true
						}
					}
					// ShowHelp control
					if string(control.Rune()) == "h" {
						if settings.ShowHelp == false {
							settings.ShowHelp = true
						} else if settings.ShowHelp == true {
							settings.ShowHelp = false
						}
					}
				}
			}
		}()

		// Do the business
		frame := gocv.NewMat()
	mainloop:
		for {
			select {
			case <-quit:
				canvas.Fini()
				break mainloop
			default:
				prev_frame_time := time.Now()
				if cam.Cap.Read(&frame) {
					asciify.Asciify(cam, &frame, canvas, settings, prev_frame_time)
					canvas.Show()
				}
			}
		}
	},
}

func Execute() error {
	return asciicamCmd.Execute()
}

func init() {
	asciicamCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	cobra.EnableCommandSorting = false
	asciicamCmd.CompletionOptions.DisableDefaultCmd = true
}
