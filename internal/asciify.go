package internal

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/gdamore/tcell/v2"
)

// var ascii_symbols = []rune(".:-~=+*#%@")
var ascii_symbols = []rune(".,;!vlLFE$")

func Asciify(frame image.Image, canvas tcell.Screen, settings *Settings, termWidth int, termHeight int, scaledResolution image.Point, defStyle tcell.Style) {
	pixStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(
		tcell.NewRGBColor(
			settings.Colour["R"],
			settings.Colour["G"],
			settings.Colour["B"]),
	)

	downFrame := imaging.Resize(frame, scaledResolution.X, scaledResolution.Y, imaging.NearestNeighbor)

	greyPixel := image.NewGray(image.Rect(0, 0, 1, 1))

	switch settings.SingleColourMode {
	case true:
		for y := downFrame.Bounds().Min.Y; y < downFrame.Bounds().Max.Y; y++ {
			for x := downFrame.Bounds().Min.X; x < downFrame.Bounds().Max.X; x++ {
				greyPixel.Set(0, 0, color.GrayModel.Convert(downFrame.At(x, y)))
				lum := greyPixel.GrayAt(0, 0).Y
				sym := ascii_symbols[int((lum)/26)]
				canvas.SetContent(x, y, sym, nil, pixStyle)
			}
		}
	case false:
		for y := downFrame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
			for x := downFrame.Bounds().Min.X; x < frame.Bounds().Max.X; x++ {
				pixelcolour := tcell.FromImageColor(downFrame.At(x, y))
				greyPixel.Set(0, 0, color.GrayModel.Convert(downFrame.At(x, y)))
				lum := greyPixel.GrayAt(0, 0).Y
				sym := ascii_symbols[int((lum)/26)]
				pixStyle = tcell.StyleDefault.
					Background(tcell.ColorReset).
					Foreground(pixelcolour)
				canvas.SetContent(x, y, sym, nil, pixStyle)
			}
		}
	}
}
