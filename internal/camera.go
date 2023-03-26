package internal

import (
	"github.com/blackjack/webcam"
	"gocv.io/x/gocv"
)

type Camera struct {
	Capture *gocv.VideoCapture
}

func GetSupportedResolutions(device string) ([]string, error) {
	cam_caps, err := webcam.Open(device)
	if err != nil {
		return nil, err
	}
	defer cam_caps.Close()

	var supportedResolutions []string
	resolutions := cam_caps.GetSupportedFrameSizes(webcam.PixelFormat(1196444237))
	if len(resolutions) == 0 {
		resolutions = cam_caps.GetSupportedFrameSizes(webcam.PixelFormat(1448695129))
	}
	for _, fs := range resolutions {
		supportedResolutions = append(supportedResolutions, fs.GetString())
	}

	return supportedResolutions, err
}

func NewCamera(device string, settings Settings) (*Camera, error) {
	capture, err := gocv.OpenVideoCaptureWithAPI(device, 200)
	if err != nil {
		return nil, err
	}

	cam := Camera{
		Capture: capture,
	}
	
	cam.Capture.Set(gocv.VideoCaptureFOURCC, cam.Capture.ToCodec("MJPG"))
	cam.Capture.Set(gocv.VideoCaptureFrameHeight, settings.FrameHeight)
	cam.Capture.Set(gocv.VideoCaptureFrameWidth, settings.FrameWidth)

	return &cam, err
}
