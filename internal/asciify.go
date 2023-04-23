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
				canvas.SetContent(x, y, sym, nil, pixStyle)
			}
		}
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
