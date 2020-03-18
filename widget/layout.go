package orgwidget

import (
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/360EntSecGroup-Skylar/excelize"
	filedialog "github.com/sqweek/dialog"
	"github.com/tanabebe/go-excel-export-ddl/ddl"
)

const (
	TableNameRow    = 5 // テーブル名が記載されている行
	TableNameColumn = 8 // テーブル名が記載されている列
)

func CreateImportButton(window fyne.Window) *fyne.Container {
	importBtn := widget.NewButton("Please select an Excel file.", func() {
		importFile, err := filedialog.File().Filter("", "xlsx").Title("ファイルを選択してください").Load()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		readFile, err := excelize.OpenFile(importFile)
		if err != nil {
			log.Fatal(err)
		}
		idxRows, _ := readFile.GetRows("目次")
		var noTargetList []string
		for i := 5; i < len(idxRows); i++ {
			// DDL除外がONなら除外対象,参照出来ないシートは無視
			if idxRows[i][52] == "ON" {
				noTargetList = append(noTargetList, idxRows[i][22])
			}
		}

		stm := ddl.Statement{}
		for _, sheet := range readFile.GetSheetMap() {
			for _, list := range noTargetList {
				rows, _ := readFile.GetRows(sheet)
				if list != rows[TableNameRow][TableNameColumn] {
					stm.GenerateDDL(rows)
				}
			}
		}

		filename, err := filedialog.File().Filter("", "sql").Title("保存する先を選択して下さい").Save()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		err = ioutil.WriteFile(filename, stm.Ddl, os.ModePerm)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
	})
	return fyne.NewContainerWithLayout(layout.NewMaxLayout(), importBtn)
}
