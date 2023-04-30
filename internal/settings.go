package internal

import (
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

type Settings struct {
	FrameHeight      float64
	FrameWidth       float64
	BrightnessCaps   map[string]int32
	ContrastCaps     map[string]int32
	Brightness       int32
	Contrast         int32
	SingleColourMode bool
	Colour           map[string]int32
	ShowInfo         bool
	ShowControls     bool
	VirtualCam       bool
}

func NewSettings(camDevice string, fH float64, fW float64) (*Settings, error) {
	s := Settings{}

	s.FrameHeight = fH
	s.FrameWidth = fW

	cam_caps, err := device.Open(camDevice)
	if err != nil {
		return nil, err
	}
	defer cam_caps.Close()

	s.BrightnessCaps = make(map[string]int32)
	brightnessCapInfo, err := cam_caps.GetControl(v4l2.CtrlBrightness)
	s.BrightnessCaps["min"] = brightnessCapInfo.Minimum
	s.BrightnessCaps["max"] = brightnessCapInfo.Maximum
	s.Brightness = brightnessCapInfo.Value

	s.ContrastCaps = make(map[string]int32)
	contrastCapInfo, err := cam_caps.GetControl(v4l2.CtrlContrast)
	s.ContrastCaps["min"] = contrastCapInfo.Minimum
	s.ContrastCaps["max"] = contrastCapInfo.Maximum
	s.Contrast = contrastCapInfo.Value

	s.SingleColourMode = true

	s.Colour = make(map[string]int32)
	s.Colour["R"] = 13
	s.Colour["G"] = 188
	s.Colour["B"] = 121

	s.ShowInfo = true
	s.ShowControls = false
	s.VirtualCam = true

	return &s, err
}
