package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	orgwidget "github.com/tanabebe/go-excel-export-ddl/widget"
)

func main() {
	a := app.New()
	window := a.NewWindow("go-excel-export-ddl")
	window.Resize(fyne.NewSize(400, 150))

	item := fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		orgwidget.CreateImportButton(window),
	)
	window.SetContent(item)
	window.ShowAndRun()
}
