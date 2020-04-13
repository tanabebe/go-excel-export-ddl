package ddl

import (
	"fmt"
	"strings"

	"github.com/tanabebe/go-excel-export-ddl/constant"
)

// Statement DDLの実行文を定義
type Statement struct {
	Ddl []byte // DDLの実行文
}

// SchemaStatement Schemaの作成を行う
func (s *Statement) SchemaStatement(schemaList []string) {
	for _, v := range schemaList {
		s.Ddl = append(s.Ddl, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;\n", v)...)
	}
}

// DropStatement DROP文を作成する
func (s *Statement) DropStatement(rows [][]string, schema string) error {
	if rows[constant.TableNameRow][constant.TableNameColumn] != "" {
		if schema == "" {
			s.Ddl = append(s.Ddl, fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", rows[constant.TableNameRow][constant.TableNameColumn])...)
		} else {
			s.Ddl = append(s.Ddl, fmt.Sprintf("DROP TABLE IF EXISTS %s.%s;\n", schema, rows[constant.TableNameRow][constant.TableNameColumn])...)
		}
	} else {
		return fmt.Errorf("%s", "drop statement error.\nUnknown table definition.")
	}
	return nil
}

// CreateStatement CREATE文を作成する
func (s *Statement) CreateStatement(rows [][]string, schema string) error {
	if rows[constant.TableNameRow][constant.TableNameColumn] != "" {
		s.Ddl = append(s.Ddl, "CREATE TABLE "...)
		if schema == "" {
			s.Ddl = append(s.Ddl, fmt.Sprintf("%s.%s", "public", rows[constant.TableNameRow][constant.TableNameColumn])...)
		} else {
			s.Ddl = append(s.Ddl, fmt.Sprintf("%s.%s", schema, rows[constant.TableNameRow][constant.TableNameColumn])...)
		}
		s.Ddl = append(s.Ddl, " (\n"...)
	} else {
		return fmt.Errorf("%s", "create statement error.\nUnknown table definition.")
	}
	return nil
}

// ColumnStatement Columnの定義文を作成する
func (s *Statement) ColumnStatement(rows [][]string, pk []string) error {
	for i := 10; i < len(rows); i++ {
		if rows[i] == nil || len(rows[i]) < 7 {
			break
		}
		if rows[i][constant.Column] != "" {
			err := s.GenerateColumn(rows, i, pk)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// PrimaryKeyStatement PKを生成する
func PrimaryKeyStatement(rows [][]string) []string {
	pk := make([]string, 0)
	for i := 10; i < len(rows); i++ {
		if rows[i] == nil || len(rows[i]) < 7 {
			break
		}
		if rows[i][constant.Column] != "" && rows[i][constant.Pk] != "" {
			pk = append(pk, rows[i][constant.Column])
		}
	}
	return pk
}

// GenerateColumn カラムの定義を全て作成する
func (s *Statement) GenerateColumn(rows [][]string, i int, pk []string) error {
	s.Ddl = append(s.Ddl, fmt.Sprintf("%4s", "")...)
	s.Ddl = append(s.Ddl, rows[i][constant.Column]...)

	switch rows[i][constant.DataType] {
	case "varchar":
		if rows[i][constant.Length] != "" {
			s.Ddl = append(s.Ddl, fmt.Sprintf(" varchar(%s)", rows[i][constant.Length])...)
		} else {
			s.Ddl = append(s.Ddl, " text"...)
		}
	case "char":
		s.Ddl = append(s.Ddl, " char"...)
	case "text":
		s.Ddl = append(s.Ddl, " text"...)
	case "smallint":
		s.Ddl = append(s.Ddl, " smallint"...)
	case "integer":
		if rows[i][constant.AutoIncrement] != "" {
			s.Ddl = append(s.Ddl, " SERIAL"...)
		} else {
			s.Ddl = append(s.Ddl, " integer"...)
		}
	case "bigint":
		s.Ddl = append(s.Ddl, " bigint"...)
	case "numeric":
		s.Ddl = append(s.Ddl, " numeric"...)
	case "date":
		s.Ddl = append(s.Ddl, " date"...)
	case "timestamp":
		s.Ddl = append(s.Ddl, " timestamp"...)
	default:
		return fmt.Errorf("%s", "Unknown table definition.")
	}
	if rows[i][constant.NotNull] != "" {
		s.Ddl = append(s.Ddl, " NOT NULL "...)
	}
	if rows[i][constant.DefaultValue] != "" {
		s.Ddl = append(s.Ddl, fmt.Sprintf(" DEFAULT %s", rows[i][constant.DefaultValue])...)
	}
	if rows[i+1] == nil || rows[i+1][constant.Column] == "" && pk != nil && len(pk) != 0 {
		s.Ddl = append(s.Ddl, ",\n"...)
		s.Ddl = append(s.Ddl, fmt.Sprintf("%4sPRIMARY KEY (%s)\n", "", strings.Join(pk, ","))...)
		s.Ddl = append(s.Ddl, ");\n"...)
	} else if rows[i+1][constant.Column] == "" && pk != nil && len(pk) == 0 {
		s.Ddl = append(s.Ddl, "\n);\n"...)
	} else {
		s.Ddl = append(s.Ddl, ",\n"...)
	}
	return nil
}

// CommentsStatement コメント用のSQLを作成する.
// 当処理まで来たら問題ないと判断しエラーは返却しない
func (s *Statement) CommentsStatement(rows [][]string, schema string) {
	if schema != "" {
		s.Ddl = append(s.Ddl, fmt.Sprintf("COMMENT ON TABLE %s.%s IS '%s';\n", schema, rows[constant.TableNameRow][constant.TableNameColumn], rows[4][constant.TableNameColumn])...)
	} else {
		s.Ddl = append(s.Ddl, fmt.Sprintf("COMMENT ON TABLE %s.%s IS '%s';\n", "public", rows[constant.TableNameRow][constant.TableNameColumn], rows[4][constant.TableNameColumn])...)
	}

	for i := 10; i < len(rows); i++ {
		if schema != "" {
			s.Ddl = append(s.Ddl, fmt.Sprintf("COMMENT ON COLUMN %s.%s.%s IS '%s';\n", schema, rows[constant.TableNameRow][constant.TableNameColumn], rows[i][constant.Column], rows[i][0])...)
		} else {
			s.Ddl = append(s.Ddl, fmt.Sprintf("COMMENT ON COLUMN %s.%s.%s IS '%s';\n", "public", rows[constant.TableNameRow][constant.TableNameColumn], rows[i][constant.Column], rows[i][0])...)
		}
		if rows[i+1][constant.Column] == "" && rows[i+1][0] == "" {
			break
		}
	}
}
