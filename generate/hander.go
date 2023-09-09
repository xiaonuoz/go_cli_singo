package generate

import (
	"strings"
)

func GenerateHanderCode(structType []StructInfo) {
	var text strings.Builder
	for _, st := range structType {
		text.WriteString("package api\n")

		text.WriteString(`func get%sByID(c *gin.Context) {
	%sID, _ := util.ParseUint(c.Query("%sID"))
	// 校验ID是否为0
	if %sID <= 0 {
		c.JSON(200, Err(serializer.CodeBindJSONErr, errors.New("%s id <= 0")))
		return
	}

	res := service.%sService.GetByID(%sID)
	c.JSON(200, res)
		}`)
	}
}
