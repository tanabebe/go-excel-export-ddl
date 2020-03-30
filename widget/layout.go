package orgwidget

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

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
var bar float64
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

		var TargetList []string
		for i := 5; i < len(idxRows); i++ {
			if idxRows[i][constant.Exclude] != "ON" {
				TargetList = append(TargetList, idxRows[i][constant.TargetSheetName])
			}
		}

		stm := ddl.Statement{}

		progress := dialog.NewProgress("start create", "please wait...", window)
		progress.Show()

		// 以降はエラーがあっても処理は中断させない
		for _, list := range TargetList {
			rows, err := readFile.GetRows(list)
			if err != nil {
				errTxt = append(errTxt, fmt.Sprintf("%s\n", err)...)
			}
			if rows != nil && list != rows[constant.TableNameRow][constant.TableNameColumn] {
				err := stm.DropStatement(rows)
				if err != nil {
					errTxt = append(errTxt, ErrMsg(rows[constant.TableNameRow][constant.TableNameColumn], err)...)
				}
				err = stm.CreateStatement(rows)
				if err != nil {
					errTxt = append(errTxt, ErrMsg(rows[constant.TableNameRow][constant.TableNameColumn], err)...)
				}
				pk := ddl.PrimaryKeyStatement(rows)
				err = stm.ColumnStatement(rows, pk)
				if err != nil {
					errTxt = append(errTxt, ErrMsg(rows[constant.TableNameRow][constant.TableNameColumn], err)...)
				}
			}
			counter += 1.0
			bar = counter / float64(len(TargetList))
			progress.SetValue(bar)
		}

		time.Sleep(time.Millisecond * 100)
		progress.Hide()

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

func ErrMsg(tableName string, err error) string {
	return fmt.Sprintf("%s : %s\n", tableName, err)
}
