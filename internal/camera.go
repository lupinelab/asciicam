package internal

import (
	"fmt"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

type Camera struct {
	Capture  v4l2.Device
}

func GetSupportedResolutions(camDevice string) ([]string, v4l2.FourCCType, error) {
	cam_caps, err := device.Open(camDevice)
	if err != nil { 
		return nil, v4l2.PixelFmtMJPEG, err
	}
	defer cam_caps.Close()

	// supportedFormats, err := v4l2.GetAllFormatDescriptions()
	// if err != nil {
	// 	return nil, err
	// }

	var pixelFormat v4l2.FourCCType
	var resolutions []v4l2.FrameSizeEnum
	var supportedResolutions []string
	
	resolutions, _ = v4l2.GetFormatFrameSizes(cam_caps.Fd(), v4l2.PixelFmtMJPEG)
	pixelFormat = v4l2.PixelFmtMJPEG
	if len(resolutions) == 0 {
		resolutions, err = v4l2.GetFormatFrameSizes(cam_caps.Fd(), v4l2.PixelFmtYUYV)
		if err != nil {
			return nil, v4l2.PixelFmtMPEG, err
		}
		pixelFormat = v4l2.PixelFmtYUYV
	}
	for _, fs := range resolutions {
		supportedResolutions = append(supportedResolutions, fmt.Sprintf("%vx%v", fs.Size.MaxWidth, fs.Size.MaxHeight))
	}

	return supportedResolutions, pixelFormat, err
}

func NewCamera(camDevice string, pixelFormat v4l2.FourCCType, settings Settings) (*device.Device, error) {
	cam, err := device.Open(
		camDevice,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: pixelFormat, Width: uint32(settings.FrameWidth), Height: uint32(settings.FrameHeight), Field: v4l2.FieldAny}),
		device.WithBufferSize(1),
	)
	if err != nil {
		return nil, err
	}

	return cam, err
}
