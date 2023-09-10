package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateServiceCode(structType []StructInfo) error {
	var text strings.Builder
	for _, st := range structType {
		text.WriteString(fmt.Sprintf("package %v\n", st.TableName))
		text.WriteString(fmt.Sprintf(`// %s Service
type %s struct {
}

`, st.Name, st.Name))

		text.WriteString(fmt.Sprintf(`func (s *%s) GetByID(id uint) *serializer.Response {
	value, err := model.%sRepo.GetByID(id)
	if err != nil {
		return serializer.Err(serializer.CodeHandlerErr, err)
	}
	return serializer.ResponseOk(value)
}

`, st.Name, st.Name))

		text.WriteString(fmt.Sprintf(`func (s *%s) %s(param serializer.%s%sParam) *serializer.Response {
	pagination := serializer.GetPagination(param.Page, param.PageSize)
	param.Page = pagination.Current
	param.PageSize = pagination.PageSize

	value, count, err := model.%sRepo.%s(param)
	if err != nil {
		return serializer.Err(serializer.CodeHandlerErr, err)
	}

	pagination.Total = count

	return serializer.ResponseOk(struct {
			%s []%s.%s %s
			Pagination *serializer.Pagination %s
		}{value, pagination})
}

`, st.Name, "Search", st.Name, "Search", st.Name, "Search", st.Name, st.TableName, st.Name, "`json:\"array\" form:\"array\"`", "`json:\"pagination\" form:\"pagination\"`"))

		funcFormat := `func (s *%s) %s(param serializer.%s%sParam) *serializer.Response {
	value, err := model.%sRepo.%s(param)
	if err != nil {
		return serializer.Err(serializer.CodeHandlerErr, err)
	}
	return serializer.ResponseOk(value)
}

`
		text.WriteString(fmt.Sprintf(funcFormat, st.Name, "List", st.Name, "List", st.Name, "List"))
		text.WriteString(fmt.Sprintf(funcFormat, st.Name, "Create", st.Name, "Create", st.Name, "Create"))
		text.WriteString(fmt.Sprintf(funcFormat, st.Name, "Modify", st.Name, "Modify", st.Name, "Modify"))
		text.WriteString(fmt.Sprintf(`func (s *%s) %s(param serializer.%s%sParam) *serializer.Response {
	err := model.%sRepo.%s(param)
	if err != nil {
		return serializer.Err(serializer.CodeHandlerErr, err)
	}
	return serializer.ResponseOk(nil)
}

`, st.Name, "Delete", st.Name, "Delete", st.Name, "Delete"))

		path := filepath.Join(ProjectDir, "service", st.TableName)
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}

		f, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("%s_service.go", st.TableName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("GenerateServiceCode err:%v", err)
		}
		defer f.Close()
		f.Write([]byte(text.String()))
	}
	return nil
}
