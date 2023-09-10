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
		text.WriteString(fmt.Sprintf(`package sdk
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

)

const (
	get%sByIDURL = "http://%%s/api/v1/%s?id=%%d"
	%sURL        = "http://%%s/api/%s"
	search%sURL  = "http://%%s/api/%s/search"
	list%sURL    = "http://%%s/api/%s/list"
)

`, st.Name, st.TableName, st.TableName, st.TableName, st.Name, st.TableName, st.Name, st.TableName))
		modelStruct := fmt.Sprintf("*%s.%s", st.TableName, st.Name)
		structString := `result := &struct {
		Code  int          %s
		Data  %s %s
		Msg   string       %s
		Error string       %s
	}{}
	
	`
		cudSting := fmt.Sprintf(structString, "`json:\"code\"`", modelStruct, "`json:\"data,omitempty\"`", "`json:\"msg\"`", "`json:\"error,omitempty\"`")

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

`, st.Name, modelStruct, st.Name, cudSting))

		reqString := `func %s%s(param serializer.%s%sParam) (%s, error) {
	%s
		p, err := json.Marshal(param)
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

		// 		text.WriteString(fmt.Sprintf(`func Create%s(param serializer.%sCreateParam) (%s, error) {
		// 	%s
		// 	request, err := json.Marshal(param)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	resp, err := http.Post(fmt.Sprintf(%sURL, host), "application/json", bytes.NewBuffer(request))
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	defer resp.Body.Close()
		// 	body, err := io.ReadAll(resp.Body)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	err = json.Unmarshal(body, result)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	if result.Code != 0 {
		// 		return nil, errors.New(result.Error)
		// 	}

		// 	return result.Data, nil
		// }`, st.Name, st.Name, modelStruct, resultString, st.TableName))

		text.WriteString(fmt.Sprintf(reqString, "Create", st.Name, st.Name, "Create", modelStruct, cudSting, "POST", st.TableName))
		text.WriteString(fmt.Sprintf(reqString, "Modify", st.Name, st.Name, "Modify", modelStruct, cudSting, "PUT", st.TableName))

		text.WriteString(fmt.Sprintf(`func %s%s(param serializer.%s%sParam) error {
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
			
`, "Delete", st.Name, st.Name, "Delete", cudSting, "DELETE", st.TableName))

		text.WriteString(fmt.Sprintf(`type %s%sResp struct {
	%s      []%s.%s          %s
	Pagination *serializer.Pagination %s
}

`, "Search", st.Name, st.Name, st.TableName, st.Name, "`json:\"array\" form:\"array\"`", "`json:\"pagination\" form:\"pagination\"`"))

		searchStruct := "*Search" + st.Name + "Resp"
		// search
		searchSting := fmt.Sprintf(structString, "`json:\"code\"`", searchStruct, "`json:\"data,omitempty\"`", "`json:\"msg\"`", "`json:\"error,omitempty\"`")

		text.WriteString(fmt.Sprintf(reqString, "Search", st.Name, st.Name, "Search", searchStruct, searchSting, "GET", st.TableName))

		text.WriteString(fmt.Sprintf(reqString, "List", st.Name, st.Name, "List", modelStruct, cudSting, "GET", st.TableName))

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
