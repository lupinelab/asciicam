package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lupinelab/asciicam/utils"
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
		settings, err := utils.NewSettings(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
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
		
		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		canvas.SetStyle(defStyle)

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
				if input.Key() == tcell.KeyEsc || input.Key() == tcell.KeyCtrlC {
					canvas.Fini()
					os.Exit(0)
				} else {
					for i, fs := range settings.Supported_resolutions {
						value, err := strconv.Atoi(string(input.Rune()))
						if err != nil {
							goto inputloop
						}
						if value == i {
							resolution = fs
							break inputloop
						} else {
							goto inputloop
						}
					}
				}
			}
		}

		canvas.Clear()
		ready_screen := []string{}
		ready_screen = append(ready_screen, fmt.Sprint("Press Enter key when ready..."))
		ready_screen = append(ready_screen, fmt.Sprint(""))
		ready_screen = append(ready_screen, fmt.Sprint("Help"))
		ready_screen = append(ready_screen, fmt.Sprint("===="))
		for _, l := range utils.Help {
			ready_screen = append(ready_screen, l)
		}
		for i, l := range ready_screen {
			for n, r := range l {
				canvas.SetContent(n, i, r, nil, defStyle)
			}
		}
		canvas.Show()
	
	ready:
		for {
			input := canvas.PollEvent()
			switch input := input.(type) {
			case *tcell.EventKey:
				if input.Key() == tcell.KeyEsc || input.Key() == tcell.KeyCtrlC {
					canvas.Fini()
					os.Exit(0)
				} else if input.Key() == tcell.KeyEnter{
					break ready
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

		cam, err := utils.NewCamera(args[0])
		if err != nil {
			panic(err)
		}
		defer cam.Cap.Close()
		cam.Cap.Set(gocv.VideoCaptureFrameWidth, fW)
		cam.Cap.Set(gocv.VideoCaptureFrameHeight, fH)

		quit := make(chan struct{})
		go func () {
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
						}
					}
					if control.Key() == tcell.KeyDown {
						if settings.Brightness > settings.Brightness_caps["min"] {
							settings.Brightness -= 1
						}
					}
					// Contrast controls
					if control.Key() == tcell.KeyRight {
						if settings.Contrast < settings.Contrast_caps["max"] {
							settings.Contrast += 1
						}
					}
					if control.Key() == tcell.KeyLeft {
						if settings.Contrast > settings.Contrast_caps["min"] {
							settings.Contrast -= 1
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
					// ShowStats control
					if string(control.Rune()) == "i" {
						if settings.ShowStats == true {
							settings.ShowStats = false	
						} else if settings.ShowStats == false {
							settings.ShowStats = true
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

		frame := gocv.NewMat()
	mainloop:
		for {
			select {
			case <-quit:
				break mainloop
			default:
				cam.Cap.Read(&frame)
				asciify.Asciify(cam, &frame, canvas, settings)
				canvas.Show()
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
