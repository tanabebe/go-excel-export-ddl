package orgwidget

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/360EntSecGroup-Skylar/excelize"
	filedialog "github.com/sqweek/dialog"
	"github.com/tanabebe/go-excel-export-ddl/constant"
	"github.com/tanabebe/go-excel-export-ddl/ddl"
)

// CreateImportButton Excelファイルを選択するfyneのボタンを返却
func CreateImportButton(window fyne.Window) *fyne.Container {

	importBtn := widget.NewButton("Please select an Excel file.", func() {
		importFile, err := filedialog.File().Filter("Excel files", "xlsx").Title("ファイルを選択してください").Load()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		readFile, err := excelize.OpenFile(importFile)
		if err != nil {
			log.Fatal(err)
		}
		idxRows, err := readFile.GetRows("目次")
		if err != nil {
			log.Fatal(err)
		}
		var noTargetList []string
		for i := 5; i < len(idxRows); i++ {
			if idxRows[i][52] == "ON" {
				noTargetList = append(noTargetList, idxRows[i][22])
			}
		}

		stm := ddl.Statement{}

		// Excel内の除外シート以外を対象とする
		for _, sheet := range readFile.GetSheetMap() {
			for _, list := range noTargetList {
				rows, err := readFile.GetRows(sheet)
				if err != nil {
					log.Fatal(err)
				}
				if list != rows[constant.TableNameRow][constant.TableNameColumn] {
					err := stm.GenerateDDL(rows)
					if err != nil {
						// errorがあれば抜けたい
						fmt.Printf("\n%s : %s", rows[constant.TableNameRow][constant.TableNameColumn], err)
					}
				}
			}
		}

		filename, err := filedialog.File().Filter("SQL files", "sql").Title("保存する先を選択して下さい").Save()
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
