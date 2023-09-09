package generate

import (
	"fmt"
	"strings"
)

var (
	paramText = `
package serializer
	`
)

func GenerateParamCode(structType []StructInfo) {
	var text strings.Builder
	for _, st := range structType {
		// 生成create param
		text.WriteString(fmt.Sprintf("\n// %sCreateParam 创建参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sCreateParam struct {\n", st.Name))
		// 遍历字段
		for index, field := range st.Field {
			text.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t// %s\n", field, st.FieldType[index], fmt.Sprintf("`json:\"%s\" form:\"%s\"`", st.Tsgs[index], st.Tsgs[index]), st.Comments[index]))
		}
		text.WriteString("}\n")

		// 生成modify param
		text.WriteString(fmt.Sprintf("\n// %sModifyParam 创建参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sModifyParam struct {\n", st.Name))
		// text.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t// %s\n", field, st.FieldType[index], fmt.Sprintf("`json:\"%s\" form:\"%s\"`", st.Tsgs[index], st.Tsgs[index]), st.Comments[index]))
		// 遍历字段
		for index, field := range st.Field {
			text.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t// %s\n", field, st.FieldType[index], fmt.Sprintf("`json:\"%s\" form:\"%s\"`", st.Tsgs[index], st.Tsgs[index]), st.Comments[index]))
		}
		text.WriteString("}\n")
	}
	fmt.Println(text.String())
}
