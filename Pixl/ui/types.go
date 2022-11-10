package ui

import (
	"pradeep/pixl/apptype"
	"pradeep/pixl/pxcanvas"
	"pradeep/pixl/swatch"

	"fyne.io/fyne/v2"
)

type AppInit struct {
	PixlCanvas *pxcanvas.PxCanvas
	PixlWindow fyne.Window
	State      *apptype.State
	Swatches   []*swatch.Swatch
}
