package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os/exec"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func FindVirtualCams() (string, error) {
	v4l2LoopbackDevices, err := exec.Command("ls", "-1", "/sys/devices/virtual/video4linux").Output()
	if err != nil {
		return "", err
	}

	v4l2LoopbackDevice := strings.Split(string(v4l2LoopbackDevices), "\n")[0]

	return fmt.Sprintf("/dev/%s", v4l2LoopbackDevice), err
}

func OutputToVirtualCam(asciiFrames chan image.Image, virtualcam string) error {
	for {
		newframe := <-asciiFrames
		fmt.Println("Received Frame")
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, newframe, nil)
		if err != nil {
			return err
		}
		byteframe := bufio.NewReader(buf)
		err = ffmpeg_go.Input("pipe:", ffmpeg_go.KwArgs{"format": "singlejpeg"}).
			WithInput(byteframe).
			Output(virtualcam, ffmpeg_go.KwArgs{"format": "v4l2"}).
			Run()
		if err != nil {
			return err
		}
	}
}
