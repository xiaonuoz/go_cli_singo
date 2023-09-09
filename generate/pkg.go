package generate

type StructInfo struct {
	Name      string
	TableName string

	Field      []string
	FieldType  []string
	Tsgs       []string
	Comments   []string
	SourceFile string
}
