package utils

import "gocv.io/x/gocv"

type Camera struct {
	Cap        *gocv.VideoCapture
	Cap_width  float64
	Cap_height float64
}

func NewCamera(device string) (*Camera, error) {
	cam := Camera{}

	newcam, err := gocv.OpenVideoCaptureWithAPI(device, 200)
	if err != nil {
		return nil, err
	}
	cam.Cap = newcam

	cam.Cap_width = cam.Cap.Get(gocv.VideoCaptureFrameWidth)
	cam.Cap_height = cam.Cap.Get(gocv.VideoCaptureFrameHeight)

	return &cam, err
}
