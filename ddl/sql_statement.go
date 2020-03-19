package ddl

import (
	"fmt"

	"github.com/tanabebe/go-excel-export-ddl/constant"
)

// Statement DDLの実行文を定義
type Statement struct {
	Ddl []byte // DDLの実行文
}

// GenerateDDL Excel内の1シート毎のDDLの全てを生成する
func (s *Statement) GenerateDDL(rows [][]string) error {
	sql := make([]byte, 0)
	sql = append(sql, fmt.Sprintf("DROP TABLE %s;\n", rows[constant.TableNameRow][constant.TableNameColumn])...)
	sql = append(sql, "CREATE TABLE "...)
	sql = append(sql, rows[constant.TableNameRow][constant.TableNameColumn]...)
	sql = append(sql, " (\n"...)

	for i := 10; i < len(rows); i++ {
		if rows[i] == nil || len(rows[i]) < 7 {
			break
		}
		if rows[i][constant.Column] != "" {
			result, err := GenerateSQLColumn(rows, i)
			if err != nil {
				return err
			}
			sql = append(sql, result...)
		}
	}
	sql = append(sql, ");\n\n"...)
	s.Ddl = append(s.Ddl, sql...)
	return nil
}

// GenerateSQLColumn DDLのカラム定義を生成
func GenerateSQLColumn(rows [][]string, i int) ([]byte, error) {
	sql := make([]byte, 0)

	sql = append(sql, fmt.Sprintf("%4s", "")...)
	sql = append(sql, rows[i][constant.Column]...)

	switch rows[i][constant.DataType] {
	case "varchar":
		if rows[i][constant.Length] != "" {
			sql = append(sql, fmt.Sprintf(" varchar(%s)", rows[i][constant.Length])...)
		} else {
			sql = append(sql, " text"...)
		}
	case "char":
		sql = append(sql, " char"...)
	case "text":
		sql = append(sql, " text"...)
	case "smallint":
		sql = append(sql, " smallint"...)
	case "integer":
		sql = append(sql, " integer"...)
	case "bigint":
		sql = append(sql, " bigint"...)
	case "numeric":
		sql = append(sql, " numeric"...)
	case "date":
		sql = append(sql, " date"...)
	case "timestamp":
		sql = append(sql, " timestamp"...)
	// 該当しないデータ型はエラーにする
	default:
		return sql, fmt.Errorf("%s", "Unknown table definition.")
	}
	if rows[i][constant.NotNull] != "" {
		sql = append(sql, " NOT NULL "...)
	}
	if rows[i][constant.Pk] != "" {
		sql = append(sql, " PRIMARY KEY"...)
	}
	if rows[i][constant.DefaultValue] != "" {
		sql = append(sql, fmt.Sprintf(" DEFAULT %s", rows[i][constant.DefaultValue])...)
	}
	if rows[i][constant.AutoIncrement] != "" {
		sql = append(sql, " SERIAL"...)
	}

	if rows[i+1][constant.Column] == "" {
		sql = append(sql, "\n"...)
	} else {
		sql = append(sql, ", \n"...)
	}
	return sql, nil
}
