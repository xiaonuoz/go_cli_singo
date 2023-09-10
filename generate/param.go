package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateParamCode(structType []StructInfo) {
	path := filepath.Join(ProjectDir, "serializer")
	// 创建文件夹
	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	pagePath := filepath.Join(path, "pagination.go")
	// 生成页面生成器工具类
	_, err = os.Stat(pagePath)
	if err != nil && os.IsNotExist(err) {
		f, err := os.Create(pagePath)
		if err != nil {
			panic(fmt.Errorf("create page file err:%v", err))
		}
		f.WriteString(fmt.Sprintf(`
package serializer

// Pagination 提供给前端的页码生成器
type Pagination struct {
	Total    uint %s
	PageSize uint %s
	Current  uint %s
}

// Page PageSize 每页数据条数
const (
	Page     = 1
	PageSize = 10
)

func GetPagination(page, pageSize uint) *Pagination {
	var p uint = Page
	if page > 0 {
		p = page
	}
	var ps uint = PageSize
	if pageSize > 0 {
		ps = pageSize
	}
	return &Pagination{
		Total:    0,
		PageSize: ps,
		Current:  p,
	}
}

`, "`json:\"total\"`", "`json:\"pageSize\"`", "`json:\"current\"`"))
	}
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

		// 生成list param
		text.WriteString(fmt.Sprintf("\n// %sListParam 查询参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sListParam struct {\n", st.Name))
		text.WriteString(rangeField.String())
		text.WriteString("}\n")

		// 生成delete param
		text.WriteString(fmt.Sprintf("\n// %sDeleteParam 删除参数\n", st.Name))
		text.WriteString(fmt.Sprintf("type %sDeleteParam struct {\n", st.Name))
		text.WriteString("\tID\tuint\t`json:\"id\" form:\"id\"`\n")
		text.WriteString("}\n")

		f, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("%s.go", st.TableName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Errorf("GenerateParamCode err:%v", err))
		}
		defer f.Close()
		f.Write([]byte(text.String()))
	}

}
