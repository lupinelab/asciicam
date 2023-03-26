package internal

import (
	"github.com/gdamore/tcell/v2"
)

var Controls = []string{
	"Brightness    [Down]/[Up] ",
	"Contrast      [Left]/[Right] ",
	"Red           [e]/[r] ",
	"Green         [f]/[g] ",
	"Blue          [v]/[b] ",
	"Colour Mode   [m]",	
	"Info          [i] ",
	"Show Controls [h] ",
	"Quit          [Esc] ",
	}

func PrintControls(canvas tcell.Screen, style tcell.Style) {
	// Show the controls
	ready_screen := []string{}
	ready_screen = append(ready_screen, "")
	ready_screen = append(ready_screen, "Controls")
	ready_screen = append(ready_screen, "--------------------------")
	ready_screen = append(ready_screen, Controls...)
	ready_screen = append(ready_screen, "--------------------------")
	ready_screen = append(ready_screen, "")
	ready_screen = append(ready_screen, "Press Enter key when ready...")
	for i, l := range ready_screen {
		for n, r := range l {
			canvas.SetContent(n, i, r, nil, style)
		}
	}
	canvas.Show()
}