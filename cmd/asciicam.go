package cmd

import (
	"bufio"
	"fmt"
	"image"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lupinelab/asciicam/internal"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

var asciicamCmd = &cobra.Command{
	Use:   "asciicam [device]",
	Short: "Turn your camera into ASCII",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get camera capabilities
		supportedResolutions, err := internal.GetSupportedResolutions(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the supported reolutions
		fmt.Println("Supported resolutions:")
		for i, fs := range supportedResolutions {
			fmt.Printf("%v) %v\n", i, fs)
		}

		var resolution string
	inputloop:
		for {
			fmt.Print("\nSelect: ")
			reader := bufio.NewReader(os.Stdin)
			// ReadString will block until the delimiter is entered
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				goto inputloop
			}
			// remove the delimeter from the string
			input = strings.TrimSuffix(input, "\n")
			// Check for a list element whose position matches the input
			for i, fs := range supportedResolutions {
				value, err := strconv.Atoi(input)
				if err != nil {
					fmt.Printf("Invalid selection: %v\n", input)
					goto inputloop
				} else if value == i {
					resolution = fs
					break inputloop
				}
			}
			fmt.Printf("Invalid selection: %v\n", input)
			goto inputloop
		}

		// Parse the selected resolution for the height and width
		fWH := strings.Split(resolution, "x")

		// Frame height
		fW, err := strconv.ParseFloat(fWH[0], 32)
		if err != nil {
			fmt.Println(err)
		}
		// Frame width
		fH, err := strconv.ParseFloat(fWH[1], 32)
		if err != nil {
			fmt.Println(err)
		}

		// Settings
		settings, err := internal.NewSettings(args[0], fH, fW)
		if err != nil {
			fmt.Println(err)
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

		// Use the default terminal style
		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		canvas.SetStyle(defStyle)

		// Show the controls
		internal.PrintControls(canvas, defStyle)

		// Setup the camera
		cam, err := internal.NewCamera(args[0], *settings)
		if err != nil {
			canvas.Fini()
			fmt.Println(err)
		}
		defer cam.Capture.Close()

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
						if settings.Brightness < settings.BrightnessCaps["max"] {
							settings.Brightness += 1
							cam.Capture.Set(gocv.VideoCaptureBrightness, settings.Brightness)
						}
					}
					if control.Key() == tcell.KeyDown {
						if settings.Brightness > settings.BrightnessCaps["min"] {
							settings.Brightness -= 1
							cam.Capture.Set(gocv.VideoCaptureBrightness, settings.Brightness)
						}
					}
					// Contrast controls
					if control.Key() == tcell.KeyRight {
						if settings.Contrast < settings.ContrastCaps["max"] {
							settings.Contrast += 1
							cam.Capture.Set(gocv.VideoCaptureContrast, settings.Contrast)
						}
					}
					if control.Key() == tcell.KeyLeft {
						if settings.Contrast > settings.ContrastCaps["min"] {
							settings.Contrast -= 1
							cam.Capture.Set(gocv.VideoCaptureContrast, settings.Contrast)
						}
					}
					// SingleColourMode control
					if string(control.Rune()) == "m" {
						if settings.SingleColourMode {
							settings.SingleColourMode = false
						} else if !settings.SingleColourMode {
							settings.SingleColourMode = true
						}
					}
					// Colour controls
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
						if settings.ShowInfo {
							settings.ShowInfo = false
						} else if !settings.ShowInfo {
							settings.ShowInfo = true
						}
					}
					// ShowHelp control
					if string(control.Rune()) == "h" {
						if !settings.ShowControls {
							settings.ShowControls = true
						} else if settings.ShowControls {
							settings.ShowControls = false
						}
					}
				}
			}
		}()

		// FPS counter
		var fps int
		var frameCount int

		go func() {
			for {
				frameCount = 0
				for settings.ShowInfo {
					fps = frameCount
					frameCount = 0
					time.Sleep(time.Second)
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
				if cam.Capture.Read(&frame) {
					termWidth, termHeight := canvas.Size()
					scale := math.Min(settings.FrameWidth/float64(termWidth), settings.FrameHeight/float64(termHeight))
					scaledResolution := image.Point{X: int(settings.FrameWidth / scale), Y: int(settings.FrameHeight / (scale * 1.8))}

					canvas.Clear()
					internal.Asciify(&frame, canvas, settings, termWidth, termHeight, scale, scaledResolution, defStyle)

					// Show info
					if settings.ShowInfo {
						for i, r := range fmt.Sprintf("Terminal=%vx%v Capture=%vx%v Scaled=%vx%v Scale=%v ",
							termWidth,
							termHeight,
							settings.FrameWidth,
							settings.FrameHeight,
							scaledResolution.X,
							scaledResolution.Y,
							1/scale) {
							canvas.SetContent(i, 0, r, nil, defStyle)
						}
						for i, r := range fmt.Sprintf("FPS=%v Brightness=%v Contrast=%v Colour=[R]%v[G]%v[B]%v ",
							fps,
							settings.Brightness,
							settings.Contrast,
							settings.Colour["R"],
							settings.Colour["G"],
							settings.Colour["B"]) {
							canvas.SetContent(i, 1, r, nil, defStyle)
						}

						// Show controls
						if settings.ShowControls {
							for y, l := range internal.Controls {
								for x, r := range l {
									canvas.SetContent(x, (scaledResolution.Y-len(internal.Controls))+y, r, nil, defStyle)
								}
							}
						}
					}
					canvas.Show()
					frameCount += 1
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
