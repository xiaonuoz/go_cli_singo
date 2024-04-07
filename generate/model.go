package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/xiaonuoz/go_cli_singo_generate_code/template"
)

func GenerateModelCode(structType []StructInfo) error {
	for _, st := range structType {
		var tableBody strings.Builder
		var paramField strings.Builder
		var createBody strings.Builder
		var whereField strings.Builder
		for index, field := range st.Field {
			commentValue := strings.TrimSpace(strings.Trim(st.Comments[index], "/"))
			tagValue := fmt.Sprintf("`gorm:\"column:%s;type:varchar(0) default '' comment '%s'\" db:\"%s\" json:\"%s\" form:\"%s\"`", st.Tsgs[index], commentValue, st.Tsgs[index], st.Tsgs[index], st.Tsgs[index])
			tableBody.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t%s\n", field, st.FieldType[index], tagValue, st.Comments[index]))

			if st.Tsgs[index] != "id" {
				createBody.WriteString(fmt.Sprintf("\t\t%v:  param.%v,\n", field, field))
				paramField.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t%s\n", field, st.FieldType[index], fmt.Sprintf("`json:\"%s\" form:\"%s\"`", st.Tsgs[index], st.Tsgs[index]), st.Comments[index]))
				switch st.FieldType[index] {
				case "string":
					whereField.WriteString(fmt.Sprintf(`
		if len(c.%v) > 0 {
			db = db.Where("%v like ?", "%%"+c.%v+"%%")
		}
			`, field, st.Tsgs[index], field))
				case "int", "uint", "int64", "uint64":
					whereField.WriteString(fmt.Sprintf(`
		if c.%v > 0 {
			db = db.Where("%v = ?", c.%v)
		}
			`, field, st.Tsgs[index], field))
				}
			}
		}

		// 读取模板文件，进行替换
		dbgo := strings.ReplaceAll(template.ModelTemplate, "${TableBody}", tableBody.String())
		dbgo = strings.ReplaceAll(dbgo, "${CreateBody}", createBody.String())
		dbgo = strings.ReplaceAll(dbgo, "${TableName}", st.TableName)
		dbgo = strings.ReplaceAll(dbgo, "${Name}", st.Name)

		dbFile, _ := os.Create("./db.go")
		dbFile.Write([]byte(dbgo))
		dbFile.Close()

		paramgo := strings.ReplaceAll(template.ParamTemplate, "${TableName}", st.TableName)
		paramgo = strings.ReplaceAll(paramgo, "${Name}", st.Name)
		paramgo = strings.ReplaceAll(paramgo, "${ParamBody}", paramField.String())
		paramgo = strings.ReplaceAll(paramgo, "${Page}", fmt.Sprintf(`Page int   %s// 查询条数
	Size  int   %s// 页码
	`, "`json:\"page\" form:\"page\"`", "`json:\"size\" form:\"size\"`"))
		paramgo = strings.ReplaceAll(paramgo, "${ID}", fmt.Sprintf(`Id uint64   %s`, "`json:\"id\" form:\"id\"`"))
		paramgo = strings.ReplaceAll(paramgo, "${whereBody}", whereField.String())
		paramgo = strings.ReplaceAll(paramgo, "${RespData}", fmt.Sprintf(`		Data     []*%s %s      //列表
		Page int                     %s //查询条数
		Size  int                     %s  //页码
		Total    int                     %s     //总数`, st.Name, "`json:\"data\" form:\"data\"`", "`json:\"page\" form:\"page\"`", "`json:\"size\" form:\"size\"`", "`json:\"total\" form:\"total\"`"))

		paramFile, _ := os.Create("./param.go")
		paramFile.Write([]byte(paramgo))
		paramFile.Close()

	}

	return nil
}
