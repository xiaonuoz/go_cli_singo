package generate

import (
	"fmt"
	"os"
	"strings"
)

func GenerateParamCode(structType []StructInfo) {
	var text strings.Builder
	for _, st := range structType {
		// 遍历字段
		var rangeField strings.Builder
		for index, field := range st.Field {
			rangeField.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t%s\n", field, st.FieldType[index], fmt.Sprintf("`json:\"%s\" form:\"%s\"`", st.Tsgs[index], st.Tsgs[index]), st.Comments[index]))
		}

		text.WriteString("package serializer\n")

		// 生成create param
		text.WriteString(fmt.Sprintf("\n// %sCreateParam 创建参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sCreateParam struct {\n", st.Name))

		text.WriteString(rangeField.String())
		text.WriteString("}\n")

		// 生成modify param
		text.WriteString(fmt.Sprintf("\n// %sModifyParam 修改参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sModifyParam struct {\n", st.Name))
		text.WriteString("\tID\tuint\t`json:\"id\" form:\"id\"`\n")
		text.WriteString(rangeField.String())
		text.WriteString("}\n")

		// 生成search param
		text.WriteString(fmt.Sprintf("\n// %sSearchParam 查询参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sSearchParam struct {\n", st.Name))
		text.WriteString(rangeField.String())
		text.WriteString("\n")
		text.WriteString("\tPageSize\tuint\t`json:\"pageSize\" form:\"pageSize\"`\n")
		text.WriteString("\tPage\tuint\t`json:\"page\" form:\"page\"`\n")
		text.WriteString("}\n")

		// 生成delete param
		text.WriteString(fmt.Sprintf("\n// %sDeletehParam 删除参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sDeletehParam struct {\n", st.Name))
		text.WriteString("\tID\tuint\t`json:\"id\" form:\"id\"`\n")
		text.WriteString("}\n")

		f, err := os.OpenFile("generate/param_ex.go", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return
		}
		defer f.Close()
		f.Write([]byte(text.String()))
	}

}
