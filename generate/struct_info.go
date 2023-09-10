package generate

// 项目路径
var ProjectDir string

type StructInfo struct {
	Name      string
	TableName string

	Field      []string
	FieldType  []string
	Tsgs       []string
	Comments   []string
	SourceFile string
}
