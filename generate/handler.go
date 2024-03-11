package generate

import (
	"os"
	"strings"

	"github.com/xiaonuoz/go_cli_singo_generate_code/template"
)

func GenerateHandlerCode(structType []StructInfo) error {
	for _, st := range structType {
		// 读取模板文件，进行替换
		api := strings.ReplaceAll(template.HandlerTemplate, "${TableName}", st.TableName)
		api = strings.ReplaceAll(api, "${Name}", st.Name)

		handlerFile, err := os.Create("./handler.go")
		if err != nil {
			return err
		}
		defer handlerFile.Close()

		handlerFile.Write([]byte(api))
	}
	return nil
}
