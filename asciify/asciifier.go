package asciify

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/lupinelab/asciicam/internal"
)

func Asciifier(work <-chan Row, canvas tcell.Screen, settings *internal.Settings, wg sync.WaitGroup) {
	for row := range work {
		pixStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(
			tcell.NewRGBColor(
				settings.Colour["R"],
				settings.Colour["G"],
				settings.Colour["B"]),
		)

		for x := row.rowData.Bounds().Min.X; x < row.rowData.Bounds().Max.X; x++ {
			lum := row.rowData.GrayAt(x, 0).Y
			sym := ascii_symbols[int(lum/26)]
			canvas.SetContent(x, row.y, sym, nil, pixStyle)
		}
	}
	wg.Done()
}
