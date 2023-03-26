package internal

import (
	"github.com/blackjack/webcam"
)

type Settings struct {
	FrameHeight    float64
	FrameWidth     float64
	BrightnessCaps map[string]float64
	ContrastCaps   map[string]float64
	Brightness     float64
	Contrast       float64
	// SupportedResolutions []string
	SingleColourMode bool
	Colour           map[string]int32
	ShowInfo         bool
	ShowControls     bool
}

func NewSettings(device string, fH float64, fW float64) (*Settings, error) {
	s := Settings{}

	s.FrameHeight = fH
	s.FrameWidth = fW

	cam_caps, err := webcam.Open(device)
	if err != nil {
		return nil, err
	}
	defer cam_caps.Close()

	capmap := cam_caps.GetControls()

	s.BrightnessCaps = make(map[string]float64)
	s.BrightnessCaps["min"] = float64(capmap[webcam.ControlID(0x00980900)].Min)
	s.BrightnessCaps["max"] = float64(capmap[webcam.ControlID(0x00980900)].Max)
	brightness, err := cam_caps.GetControl(webcam.ControlID(0x00980900))
	if err != nil {
		return nil, err
	}
	s.Brightness = float64(brightness)

	s.ContrastCaps = make(map[string]float64)
	s.ContrastCaps["min"] = float64(capmap[webcam.ControlID(0x00980901)].Min)
	s.ContrastCaps["max"] = float64(capmap[webcam.ControlID(0x00980901)].Max)
	contrast, err := cam_caps.GetControl(webcam.ControlID(0x00980901))
	if err != nil {
		return nil, err
	}
	s.Contrast = float64(contrast)

	// s.SupportedResolutions = []string{}
	// resolutions := cam_caps.GetSupportedFrameSizes(webcam.PixelFormat(1196444237))
	// if len(resolutions) == 0 {
	// 	resolutions = cam_caps.GetSupportedFrameSizes(webcam.PixelFormat(1448695129))
	// }
	// for _, fs := range resolutions {
	// 	s.SupportedResolutions = append(s.SupportedResolutions, fs.GetString())
	// }

	s.SingleColourMode = true

	s.Colour = make(map[string]int32)
	s.Colour["R"] = 13
	s.Colour["G"] = 188
	s.Colour["B"] = 121

	s.ShowInfo = true
	s.ShowControls = false

	return &s, err
}
