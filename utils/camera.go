package utils

import (

	"gocv.io/x/gocv"
)

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
	cam.Cap.Set(gocv.VideoCaptureFOURCC, cam.Cap.ToCodec("MJPG"))

	return &cam, err
}
