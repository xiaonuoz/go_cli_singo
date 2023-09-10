package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fatih/camelcase"
)

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

func GetStructInfoArr(sourceFilePath string) []StructInfo {
	sourceFile, err := os.ReadFile(sourceFilePath)
	if err != nil {
		fmt.Println("无法读取源代码文件:", err)
		return nil
	}

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, sourceFilePath, sourceFile, parser.ParseComments)
	if err != nil {
		fmt.Println("语法解析错误:", err)
		return nil
	}

	var StructInfoArr []StructInfo
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						var structInfo StructInfo
						structInfo.Name = ts.Name.String()
						if len(structInfo.Name) < 2 {
							panic("结构体名不能单字符")
						}
						structInfo.TableName = strings.ToLower(structInfo.Name[:1]) + structInfo.Name[1:]

						if st, ok := ts.Type.(*ast.StructType); ok {

							for _, field := range st.Fields.List {
								var comment string

								if field.Doc != nil {
									for _, c := range field.Doc.List {
										comment = c.Text
									}
								}

								if field.Comment != nil {
									for _, c := range field.Comment.List {
										comment = c.Text
									}
								}

								fieldComment := comment
								for _, fieldName := range field.Names {

									entries := camelcase.Split(fieldName.Name)
									var tagNameArr []string
									for _, v := range entries {
										tagNameArr = append(tagNameArr, strings.ToLower(v))
									}
									tagName := strings.Join(tagNameArr, "_")

									structInfo.Field = append(structInfo.Field, fmt.Sprint(fieldName.Name))
									structInfo.FieldType = append(structInfo.FieldType, fmt.Sprint(field.Type))
									structInfo.Tsgs = append(structInfo.Tsgs, tagName)
									structInfo.Comments = append(structInfo.Comments, strings.ReplaceAll(fieldComment, "\n", ""))
								}
							}
						}
						structInfo.SourceFile = sourceFilePath
						StructInfoArr = append(StructInfoArr, structInfo)
					}
				}
			}
		}
	}

	if len(StructInfoArr) == 0 {
		fmt.Println("文件中不存在结构体！")
		return StructInfoArr
	}
	return StructInfoArr
}
