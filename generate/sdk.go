package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateSDKCode(structType []StructInfo) error {
	var text strings.Builder
	for _, st := range structType {
		var rangeField strings.Builder
		for index, field := range st.Field {
			switch st.FieldType[index] {
			case "string":
				rangeField.WriteString(fmt.Sprintf(`
	if len(param.%s) > 0 {
		up.Add("%s", param.%s)
	}
`, field, st.Tsgs[index], field))
			case "int", "uint", "int64", "uint64":
				rangeField.WriteString(fmt.Sprintf(`
		if param.%s > 0 {
			up.Add("%s", strconv.FormatUint(uint64(param.%s), 10))
		}
`, field, st.Tsgs[index], field))
			}
		}

		text.WriteString(fmt.Sprintf(`package sdk
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

)

var (
	get%sByIDURL = "http://%%s/api/v1/%s?id=%%d"
	%sURL        = "http://%%s/api/v1/%s"
	search%sURL  = "http://%%s/api/v1/%s/search"
	list%sURL    = "http://%%s/api/v1/%s/list"
)

`, st.Name, st.TableName, st.LocalName, st.TableName, st.Name, st.TableName, st.Name, st.TableName))
		modelStruct := fmt.Sprintf("*%s.%s", st.TableName, st.Name)
		listStruct := fmt.Sprintf("[]%s.%s", st.TableName, st.Name)
		structString := `result := &struct {
		Code  int          %s
		Data  %s %s
		Msg   string       %s
		Error string       %s
	}{}
	`
		rudStructSting := fmt.Sprintf(structString, "`json:\"code\"`", modelStruct, "`json:\"data,omitempty\"`", "`json:\"msg\"`", "`json:\"error,omitempty\"`")
		createStructSting := fmt.Sprintf(structString, "`json:\"code\"`", "interface{}", "`json:\"data,omitempty\"`", "`json:\"msg\"`", "`json:\"error,omitempty\"`")
		listStructSting := fmt.Sprintf(structString, "`json:\"code\"`", listStruct, "`json:\"data,omitempty\"`", "`json:\"msg\"`", "`json:\"error,omitempty\"`")

		text.WriteString(fmt.Sprintf(`func Get%sByID(id uint) (%s, error) {
	url := fmt.Sprintf(get%sByIDURL, host, id)
	%s
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, errors.New(result.Error)
	}

	return result.Data, nil
}

`, st.Name, modelStruct, st.Name, rudStructSting))

		reqString := `func %s%s(param *serializer.%s%sParam) (%s, error) {
	%s
	%s
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, errors.New(result.Error)
	}

	return result.Data, nil
}

`

		crudHttpString := `p, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("%s", fmt.Sprintf(%sURL, host), bytes.NewBuffer(p))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
`

		text.WriteString(fmt.Sprintf(`func %s%s(param *serializer.%s%sParam) error {
			%s
			%s
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	
	if result.Code != 0 {
		return errors.New(result.Error)
	}
		
	return nil
}
		
`, "Create", st.Name, st.Name, "Create", createStructSting, fmt.Sprintf(`p, err := json.Marshal(param)
if err != nil {
	return err
}

req, err := http.NewRequest("%s", fmt.Sprintf(%sURL, host), bytes.NewBuffer(p))
if err != nil {
	return err
}
req.Header.Set("Content-Type", "application/json")

resp, err := http.DefaultClient.Do(req)
if err != nil {
	return err
}
`, "POST", st.LocalName)))
		text.WriteString(fmt.Sprintf(reqString, "Modify", st.Name, st.Name, "Modify", modelStruct, rudStructSting, fmt.Sprintf(crudHttpString, "PUT", st.LocalName)))

		text.WriteString(fmt.Sprintf(`func %s%s(param *serializer.%s%sParam) error {
	%s
	p, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("%s", fmt.Sprintf(%sURL, host), bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Code != 0 {
		return errors.New(result.Error)
	}

	return nil
}
			
`, "Delete", st.Name, st.Name, "Delete", rudStructSting, "DELETE", st.LocalName))

		text.WriteString(fmt.Sprintf(`type %s%sResp struct {
	%s      []%s.%s          %s
	Pagination *serializer.Pagination %s
}

`, "Search", st.Name, st.Name, st.TableName, st.Name, "`json:\"array\" form:\"array\"`", "`json:\"pagination\" form:\"pagination\"`"))

		searchStruct := "*Search" + st.Name + "Resp"
		// search & list
		searchSting := fmt.Sprintf(structString, "`json:\"code\"`", searchStruct, "`json:\"data,omitempty\"`", "`json:\"msg\"`", "`json:\"error,omitempty\"`")

		searchListHttpString := `baseURL, err := url.Parse(fmt.Sprintf(%sURL, host))
	if err != nil {
		return nil, err
	}

	up := url.Values{}
	%s

	%s
	baseURL.RawQuery = up.Encode()
	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}
`
		text.WriteString(fmt.Sprintf(reqString, "Search", st.Name, st.Name, "Search", searchStruct, searchSting, fmt.Sprintf(searchListHttpString, "search"+st.Name, rangeField.String(), `	up.Add("page", strconv.FormatUint(uint64(param.Page), 10))
		up.Add("pageSize", strconv.FormatUint(uint64(param.PageSize), 10))`)))

		text.WriteString(fmt.Sprintf(reqString, "List", st.Name, st.Name, "List", listStruct, listStructSting, fmt.Sprintf(searchListHttpString, "list"+st.Name, rangeField.String(), "")))

		path := filepath.Join(ProjectDir, "sdk")
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
		f, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("%s.go", st.TableName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("GenerateSDKCode err:%v", err)
		}
		defer f.Close()
		f.Write([]byte(text.String()))
	}
	return nil
}
