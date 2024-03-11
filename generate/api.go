package generate

import (
	"os"
	"strings"

	"github.com/xiaonuoz/go_cli_singo_generate_code/template"
)

func GenerateApiCode(structType []StructInfo) error {
	for _, st := range structType {
		// 读取模板文件，进行替换
		api := strings.ReplaceAll(template.ApiTemplate, "${TableName}", st.TableName)
		api = strings.ReplaceAll(api, "${Name}", st.Name)
		api = strings.ReplaceAll(api, "${LocalName}", st.LocalName)

		paramFile, err := os.Create("./api.go")
		if err != nil {
			return err
		}
		defer paramFile.Close()
		paramFile.Write([]byte(api))
	}
	return nil
}
