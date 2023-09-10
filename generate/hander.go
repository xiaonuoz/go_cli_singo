package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateHandlerCode(structType []StructInfo) error {
	var text strings.Builder
	for _, st := range structType {
		text.WriteString("package api\n")

		text.WriteString(fmt.Sprintf(`
func get%sByID(c *gin.Context) {
	id, _ := util.ParseUint(c.Query("id"))
	// 校验ID是否为0
	if id <= 0 {
		c.JSON(200, Err(serializer.CodeBindJSONErr, errors.New("%s id <= 0")))
		return
	}

	res := service.%sService.GetByID(id)
	c.JSON(200, res)
}

`, st.Name, st.Name, st.Name))

		funcFormat := `func %s%s(c *gin.Context) {
	param := serializer.%s%sParam{}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, Err(serializer.CodeBindJSONErr, err))
		return
	}

	res := service.%sService.%s(param)
	c.JSON(200, res)
}

`

		text.WriteString(fmt.Sprintf(funcFormat, "search", st.Name, st.Name, "Search", st.Name, "Search"))
		text.WriteString(fmt.Sprintf(funcFormat, "list", st.Name, st.Name, "List", st.Name, "List"))
		text.WriteString(fmt.Sprintf(funcFormat, "modify", st.Name, st.Name, "Modify", st.Name, "Modify"))
		text.WriteString(fmt.Sprintf(funcFormat, "create", st.Name, st.Name, "Create", st.Name, "Create"))
		text.WriteString(fmt.Sprintf(funcFormat, "delete", st.Name, st.Name, "Delete", st.Name, "Delete"))

		path := filepath.Join(ProjectDir, "api")
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
		f, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("%s_handler.go", st.TableName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("GenerateHandlerCode err:%v", err)
		}
		defer f.Close()
		f.Write([]byte(text.String()))
	}
	return nil
}
