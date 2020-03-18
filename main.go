package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	orgwidget "github.com/tanabebe/go-excel-export-ddl/widget"
)

func main() {
	a := app.New()
	window := a.NewWindow("go-excel-export-ddl")
	window.Resize(fyne.NewSize(400, 150))
	window.SetContent(orgwidget.CreateImportButton(window))
	window.ShowAndRun()
}
