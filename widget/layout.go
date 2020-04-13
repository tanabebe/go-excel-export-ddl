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

		var targetList []string
		var schema []string
		for i := 5; i < len(idxRows); i++ {
			if idxRows[i][constant.SheetIndexTableName] != "" && idxRows[i][constant.Exclude] != "ON" {
				targetList = append(targetList, idxRows[i][constant.TargetSheetName])
				schema = append(schema, idxRows[i][constant.SchemaName])
			}
		}

		progress := dialog.NewProgress("start create", "please wait...", window)
		progress.Show()

		errTxt = make([]byte, 0)
		var schemaIndex int
		stm := ddl.Statement{}
		mergeSlice := make([]string, 0)

		// スキーマの重複を削除してスライスの作成し直し
		for _, v := range schema {
			if v != "" && !SchemaContains(mergeSlice, v) {
				mergeSlice = append(mergeSlice, v)
			}
		}

		// schema指定がなくてもエラーとはさせない
		stm.SchemaStatement(mergeSlice)

		// 以降はエラーがあっても処理は中断させない
		for _, list := range targetList {
			rows, err := readFile.GetRows(list)
			if err != nil {
				errTxt = append(errTxt, fmt.Sprintf("%s\n", err)...)
			}
			if rows != nil && list != rows[constant.TableNameRow][constant.TableNameColumn] {
				if err := stm.DropStatement(rows, schema[schemaIndex]); err != nil {
					errTxt = append(errTxt, ErrMsg(rows[constant.TableNameRow][constant.TableNameColumn], err)...)
				}
				if err := stm.CreateStatement(rows, schema[schemaIndex]); err != nil {
					errTxt = append(errTxt, ErrMsg(rows[constant.TableNameRow][constant.TableNameColumn], err)...)
				}
				pk := ddl.PrimaryKeyStatement(rows)
				if err := stm.ColumnStatement(rows, pk); err != nil {
					errTxt = append(errTxt, ErrMsg(rows[constant.TableNameRow][constant.TableNameColumn], err)...)
				}
				stm.CommentsStatement(rows, schema[schemaIndex])
			}
			schemaIndex += 1
			counter += 1.0
			bar = counter / float64(len(targetList))
			progress.SetValue(bar)
		}

		time.Sleep(time.Millisecond * 100)
		progress.Hide()

		if len(errTxt) != 0 {
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

// ErrMsg　画面表示用のメッセージを返却する
func ErrMsg(tableName string, err error) string {
	return fmt.Sprintf("%s : %s\n", tableName, err)
}

// SchemaContains 重複するスキーマ名を除去する
func SchemaContains(sl []string, str string) bool {
	for _, v := range sl {
		if str == v {
			return true
		}
	}
	return false
}
