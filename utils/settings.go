package utils

import (
	"github.com/blackjack/webcam"
)

type Settings struct {
	Brightness_caps       map[string]float64
	Contrast_caps         map[string]float64
	Brightness            float64
	Contrast              float64
	Supported_resolutions []string
	Colour                map[string]int32
	ShowStats             bool
	ShowHelp              bool
}

func NewSettings(device string) (*Settings, error) {
	s := Settings{}

	cam_caps, err := webcam.Open(device)
	if err != nil {
		return nil, err
	}
	defer cam_caps.Close()

	capmap := cam_caps.GetControls()

	s.Brightness_caps = make(map[string]float64)
	s.Brightness_caps["min"] = float64(capmap[webcam.ControlID(0x00980900)].Min)
	s.Brightness_caps["max"] = float64(capmap[webcam.ControlID(0x00980900)].Max)
	brightness, err := cam_caps.GetControl(webcam.ControlID(0x00980900))
	if err != nil {
		return nil, err
	}
	s.Brightness = float64(brightness)

	s.Contrast_caps = make(map[string]float64)
	s.Contrast_caps["min"] = float64(capmap[webcam.ControlID(0x00980901)].Min)
	s.Contrast_caps["max"] = float64(capmap[webcam.ControlID(0x00980901)].Max)
	contrast, err := cam_caps.GetControl(webcam.ControlID(0x00980901))
	if err != nil {
		return nil, err
	}
	s.Contrast = float64(contrast)

	s.Supported_resolutions = []string{}
	resolutions := cam_caps.GetSupportedFrameSizes(webcam.PixelFormat(1196444237))
	for _, fs := range resolutions {
		s.Supported_resolutions = append(s.Supported_resolutions, fs.GetString())
	}

	s.Colour = make(map[string]int32)
	s.Colour["R"] = 13
	s.Colour["G"] = 188
	s.Colour["B"] = 123

	s.ShowStats = true
	s.ShowHelp = false

	return &s, err
}
