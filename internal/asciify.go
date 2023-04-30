package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/gdamore/tcell/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// var ascii_symbols = []rune(".:-~=+*#%@")
var ascii_symbols = []rune(".,;!vlLFE$")

func Asciify(frame image.Image, canvas tcell.Screen, settings *Settings, termWidth int, termHeight int, scale float64, scaledResolution image.Point, defStyle tcell.Style, framesForVirtualCam *chan image.Image) {

	pixStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(
		tcell.NewRGBColor(
			settings.Colour["R"],
			settings.Colour["G"],
			settings.Colour["B"]),
	)

	downFrame := imaging.Resize(
		frame,
		scaledResolution.X,
		scaledResolution.Y,
		imaging.NearestNeighbor)

	greyFrame := image.NewGray(
		image.Rect(
			downFrame.Bounds().Min.X,
			downFrame.Bounds().Min.Y,
			downFrame.Bounds().Max.X,
			downFrame.Bounds().Max.Y,
		),
	)

	var frameForVirtualCam draw.Image

	switch settings.VirtualCam {
	case true:
		frameForVirtualCam = image.NewRGBA(
			image.Rect(
				frame.Bounds().Min.X,
				frame.Bounds().Min.Y,
				frame.Bounds().Max.X,
				frame.Bounds().Max.Y,
			),
		)
	}

	for y := greyFrame.Bounds().Min.Y; y < greyFrame.Bounds().Max.Y; y++ {
		for x := greyFrame.Bounds().Min.X; x < greyFrame.Bounds().Max.X; x++ {
			greyFrame.Set(x, y, color.GrayModel.Convert(downFrame.At(x, y)))
		}
	}

	switch settings.SingleColourMode {
	case true:
		for y := greyFrame.Bounds().Min.Y; y < greyFrame.Bounds().Max.Y; y++ {
			for x := greyFrame.Bounds().Min.X; x < greyFrame.Bounds().Max.X; x++ {
				lum := greyFrame.GrayAt(x, y).Y
				sym := ascii_symbols[int(lum/26)]
				switch settings.VirtualCam {
				case false:
					canvas.SetContent(x, y, sym, nil, pixStyle)
				case true:
					d := font.Drawer{
						Dst: frameForVirtualCam,
						Src: image.NewUniform(color.RGBA{
							R: uint8(settings.Colour["R"]),
							G: uint8(settings.Colour["G"]),
							B: uint8(settings.Colour["B"]),
							A: 0xFF,
						}),
						Face: basicfont.Face7x13,
						Dot: fixed.P(
							int(int(float64(y)*scale)),
							int(int(float64(x)*scale)),
						),
					}
					symbol := []byte{byte(sym)}
					d.DrawBytes(symbol)
				}
			}
		}
		fmt.Println("Sending Frame")
		*framesForVirtualCam <- frameForVirtualCam
		fmt.Println("Sent Frame")

	case false:
		for y := downFrame.Bounds().Min.Y; y < downFrame.Bounds().Max.Y; y++ {
			for x := downFrame.Bounds().Min.X; x < downFrame.Bounds().Max.X; x++ {
				pixelcolour := tcell.FromImageColor(downFrame.At(x, y))
				lum := greyFrame.GrayAt(x, y).Y
				sym := ascii_symbols[int(lum/26)]
				pixStyle = tcell.StyleDefault.
					Background(tcell.ColorReset).
					Foreground(pixelcolour)
				canvas.SetContent(x, y, sym, nil, pixStyle)
			}
		}
	}
}
