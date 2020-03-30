package constant

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
	Exclude         = 52 // DDLの除外対象記載位置
	TargetSheetName = 3  // DDLの除外対象だった場合に使用するシート名
)
