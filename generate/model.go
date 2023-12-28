package generate

import (
	"fmt"
	"os"
	"strings"
)

func GenerateModelCode(structType []StructInfo) error {
	var text strings.Builder
	for _, st := range structType {
		var rangeField strings.Builder
		var paramField strings.Builder
		var updateField strings.Builder
		var whereField strings.Builder
		for index, field := range st.Field {
			commentValue := strings.TrimSpace(strings.Trim(st.Comments[index], "/"))
			tagValue := fmt.Sprintf("`gorm:\"column:%s;type:varchar(0) default '' comment '%s'\" db:\"%s\" json:\"%s\" form:\"%s\"`", st.Tsgs[index], commentValue, st.Tsgs[index], st.Tsgs[index], st.Tsgs[index])
			rangeField.WriteString(fmt.Sprintf("\t%s\t%s\t%s\t%s\n", field, st.FieldType[index], tagValue, st.Comments[index]))

			if st.Tsgs[index] != "id" {
				updateField.WriteString(fmt.Sprintf("\t\t%v:  param.%v,\n", field, field))
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

		// 增加parama
		where := genWhere(paramField, whereField, st.TableName)
		fw, _ := os.Create("./where.go")
		fw.Write([]byte(where.String()))

		// 为结构体增加tag
		// 增加CRUD
		text = genModel(text, st, rangeField, updateField)
		f, _ := os.Create("./db.go")
		f.Write([]byte(text.String()))
	}

	return nil
}

func genModel(text strings.Builder, st StructInfo, rangeField strings.Builder, updateField strings.Builder) strings.Builder {
	text.WriteString(fmt.Sprintf("type %s struct {\n", st.Name))
	text.WriteString(rangeField.String())
	text.WriteString("}\n")

	text.WriteString(fmt.Sprintf(`func (%s) TableName() string {
				return "%s"
			}
		
		`, st.Name, st.TableName))

	text.WriteString(fmt.Sprintf(`func InitTable() {
			table := &%s{}
			err := db.DB().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&%s{})
			if err != nil {
				return
			}
			db.DB().Exec("ALTER TABLE " + table.TableName() + " COMMENT '%s';")
		}
		
`, st.Name, st.Name, st.Name))

	// get func
	text.WriteString(fmt.Sprintf(`func Get(id uint64) (*%s, error) {
		var o = &%s{}
		err := db.DB().Model(%s{}).Where("id = ?", id).Find(&o).Error
		return o, err
	}
	
`, st.Name, st.Name, st.Name))

	// NameIsExist
	text.WriteString(fmt.Sprintf(`// 检测名称是否已存在
	func NameIsExist(name string) (bool, error) {
		var c int64
		err := db.DB().Model(%s{}).Where("name=?", name).Count(&c).Error
		if err != nil {
			return false, err
		}
		return c > 0, nil
	}
	
`, st.Name))

	// list func
	text.WriteString(fmt.Sprintf(`// 分页查询
func GetList(param *ListParam) (res []*%s, total int64, err error) {
	model := db.Model(&%s{}).Scopes(param.where(), param.order(), param.preload())

	// 分页
	if param.PageNum <= 0 {
		param.PageNum = 1
	}
	if param.PageSize < 1 {
		param.PageSize = 10
	}
	err = model.Count(&total).Error
	if err != nil {
		return
	}

	err = model.Offset(param.PageSize * (param.PageNum - 1)).Limit(param.PageSize).Find(&res).Error
	return
}

`, st.Name, st.Name))

	text.WriteString(fmt.Sprintf(`func Create(param *CreateParam, tx *gorm.DB) (err error) {
		if tx == nil {
			tx = db.DB()
		}
	
		err = tx.Create(&%s{
		%s
		}).Error
		if err != nil {
			return err
		}
		return
	}
	
	`, st.Name, updateField.String()))

	text.WriteString(fmt.Sprintf(`func Delete(param *DelParam, tx *gorm.DB) (err error) {
		if tx == nil {
			tx = db.DB()
		}
		return tx.Where("id = ?", param.Id).Delete(&%s{}).Error
	}
	
	`, st.Name))

	// update
	text.WriteString(fmt.Sprintf(`func Update(param *UpdateParam, tx *gorm.DB) error {
		if tx == nil {
			tx = db.DB()
		}
		if err := tx.Where("id=?", param.Id).Updates(&NetworkSegmentTable{
			%s
		}).Error; err != nil {
			return err
		}
	
		return nil
	}
	
	`, updateField.String()))

	return text
}

func genWhere(paramField strings.Builder, whereField strings.Builder, tableName string) strings.Builder {
	var where strings.Builder
	where.WriteString(fmt.Sprintf(`package %s

import "gorm.io/gorm"

`, tableName))
	where.WriteString("type (\n")

	paramString := `%s struct{
					%s
		}
			
`
	where.WriteString(fmt.Sprintf(`%s struct{
		%s
PageSize int   %s// 查询条数
PageNum  int   %s// 页码
	}

`, "ListParam", paramField.String(), "`json:\"page_size\" form:\"page_size\"`", "`json:\"page_num\" form:\"page_num\"`"))
	where.WriteString(fmt.Sprintf(paramString, "CreateParam", paramField.String()))
	where.WriteString(fmt.Sprintf(`%s struct{
		%s
		%s
}

`, "UpdateParam", "Id     uint64 `json:\"id\" form:\"id\"`", paramField.String()))
	where.WriteString(fmt.Sprintf(paramString, "DelParam", "Id     uint64 `json:\"id\" form:\"id\"`"))
	where.WriteString(")\n")

	where.WriteString(fmt.Sprintf(`// where 条件
func (c *ListParam) where() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
%v

		return db
	}
}

`, whereField.String()))

	where.WriteString(`// order 处理排序规则
func (c *ListParam) order() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order("created_at DESC")
		return db
	}
}

// preload 预加载表
func (c *ListParam) preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// if c.PreloadProperty {
		// 	db = db.Preload("PropertyArr")
		// }
		return db
	}
}
`)

	return where
}
