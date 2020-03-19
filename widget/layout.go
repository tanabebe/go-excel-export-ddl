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

var counter float64
var errTxt []byte

// CreateImportButton Excelファイルを選択するfyneのボタンを返却
func CreateImportButton(window fyne.Window) *widget.Box {

	lbl := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

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
			if idxRows[i][constant.Exclude] == "ON" {
				noTargetList = append(noTargetList, idxRows[i][constant.ExcludeTableName])
			}
		}

		stm := ddl.Statement{}

		prog := dialog.NewProgress("start create", "please wait...", window)
		prog.Show()

		for _, sheet := range readFile.GetSheetMap() {
			for _, list := range noTargetList {

				rows, err := readFile.GetRows(sheet)
				if err != nil {
					log.Fatal(err)
				}
				if list != rows[constant.TableNameRow][constant.TableNameColumn] {
					err := stm.GenerateDDL(rows)
					if err != nil {
						fmt.Sprintf("%s : %s\n", rows[constant.TableNameRow][constant.TableNameColumn], err)
						errTxt = append(errTxt, fmt.Sprintf("%s : %s\n", rows[constant.TableNameRow][constant.TableNameColumn], err)...)
					}
				}
			}
			counter += (counter + 1) / float64(readFile.SheetCount)
			prog.SetValue(counter)
		}
		prog.Hide()

		if errTxt != nil {
			lbl.SetText(string(errTxt))
		} else {
			lbl.SetText("All created successfully!!")
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

	btn := fyne.NewContainerWithLayout(layout.NewMaxLayout(), importBtn)

	return widget.NewVBox(btn, lbl)
}
