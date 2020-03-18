package ddl

import (
	"fmt"
)

type Statement struct {
	Ddl []byte
}

const (
	TableNameRow    = 5  // テーブル名が記載されている行
	TableNameColumn = 8  // テーブル名が記載されている列
	Column          = 7  // テーブルのカラム名の記載位置
	DataType        = 14 // テーブルのデータ型の記載位置
	Length          = 17 // テーブルのデータ型の長さ記載位置
	Pk              = 23 // プライマリキーの記載位置
	NotNull         = 26 // NOT NULLの記載位置
	AutoIncrement   = 35 // AUTO INCREMENTの記載位置
	DefaultValue    = 38 // DEFAULT VALUEの記載位置
)

func (s *Statement) GenerateDDL(rows [][]string) {
	s.Ddl = append(s.Ddl, "CREATE TABLE "...)
	s.Ddl = append(s.Ddl, rows[TableNameRow][TableNameColumn]...)
	s.Ddl = append(s.Ddl, " (\n"...)

	for i := 10; i < len(rows); i++ {
		if rows[i] == nil || len(rows[i]) < 7 {
			break
		}
		if rows[i][Column] != "" {
			s.Ddl = s.GenerateColumnStruct(rows, i)
		}
	}
	s.Ddl = append(s.Ddl, ");\n\n"...)
}

func (s *Statement) GenerateColumnStruct(rows [][]string, i int) []byte { //, ddl []byte) []byte {
	s.Ddl = append(s.Ddl, "\t"...)
	s.Ddl = append(s.Ddl, rows[i][Column]...)

	switch rows[i][DataType] {
	case "varchar":
		if rows[i][Length] != "" {
			s.Ddl = append(s.Ddl, fmt.Sprintf(" varchar(%s)", rows[i][Length])...)
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
		s.Ddl = append(s.Ddl, " integer"...)
	case "bigint":
		s.Ddl = append(s.Ddl, " bigint"...)
	case "numeric":
		s.Ddl = append(s.Ddl, " numeric"...)
	case "date":
		s.Ddl = append(s.Ddl, " date"...)
	case "timestamp":
		s.Ddl = append(s.Ddl, " timestamp"...)
	default:
		return s.Ddl
	}
	if rows[i][NotNull] != "" {
		s.Ddl = append(s.Ddl, " NOT NULL "...)
	}
	if rows[i][Pk] != "" {
		s.Ddl = append(s.Ddl, " PRIMARY KEY"...)
	}
	if rows[i][DefaultValue] != "" {
		s.Ddl = append(s.Ddl, fmt.Sprintf(" DEFAULT %s", rows[i][DefaultValue])...)
	}
	if rows[i][AutoIncrement] != "" {
		s.Ddl = append(s.Ddl, " SERIAL"...)
	}

	// 最終行の次カラムが空なら終了とみなす
	if rows[i+1][Column] == "" {
		s.Ddl = append(s.Ddl, "\n"...)
	} else {
		s.Ddl = append(s.Ddl, ", \n"...)
	}
	return s.Ddl
}
