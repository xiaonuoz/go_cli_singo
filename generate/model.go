package generate

import (
	"fmt"
	"os"
	"strings"
)

var (
	modelText = `

	`
)

func GenerateModelCode(structType []StructInfo) {
	var text strings.Builder
	for _, st := range structType {

		// 遍历字段
		var rangeField strings.Builder
		var rangeField1 strings.Builder
		var rangeField2 strings.Builder
		var rangeFieldWhere strings.Builder

		for index, field := range st.Field {
			rangeField.WriteString(fmt.Sprintf("\t\t%v:  param.%v,\n", field, field))
			rangeField1.WriteString(fmt.Sprintf("\tobj.%v = param.%v\n", field, field))
			rangeField2.WriteString(fmt.Sprintf("\t%v\t%v\t%v\n", field, st.FieldType[index], st.Comments[index]))
			switch st.FieldType[index] {
			case "string":
				rangeFieldWhere.WriteString(fmt.Sprintf(`
		if len(c.%v) > 0 {
			db = db.Where("%v like ?", "%%"+c.%v+"%%")
		}
`, field, st.Tsgs[index], field))
			case "int", "uint", "int64", "uint64":
				rangeFieldWhere.WriteString(fmt.Sprintf(`
		if c.%v > 0 {
			db = db.Where("%v = ?", c.%v)
		}
`, field, st.Tsgs[index], field))
			}
		}

		// if err := genModelSQL(st.TableName, st, text, rangeField, rangeField1); err != nil {
		// 	panic(fmt.Errorf("genModelSQL err:%v", err))
		// }

		if err := genModelQuery(text, st.TableName, rangeField2, rangeFieldWhere); err != nil {
			panic(fmt.Errorf("genModelQuery err:%v", err))
		}
	}

}

func genModelQuery(text strings.Builder, tableName string, rangeField2 strings.Builder, rangeFieldWhere strings.Builder) error {
	text.WriteString(fmt.Sprintf(`package model

import "gorm.io/gorm"

// %sQuery %s查询条件
type %sQuery struct {
	ID           uint
%v

	Limit  uint
	Offset uint
}

`, tableName, tableName, tableName, rangeField2.String()))

	text.WriteString(fmt.Sprintf(`// where 条件
func (c %sQuery) where() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if c.ID > 0 {
			db = db.Where("id = ?", c.ID)
		}
%v
		return db
	}
}

`, tableName, rangeFieldWhere.String()))

	text.WriteString(fmt.Sprintf(`// order 处理排序规则
func (c %vQuery) order() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order("created_at DESC")
		return db
	}
}

// preload 预加载表
func (c %vQuery) preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// if c.PreloadProperty {
		// 	db = db.Preload("PropertyArr")
		// }
		return db
	}
}
`, tableName, tableName))
	fmt.Println(text.String())

	f, err := os.OpenFile("generate/model.ex.go", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(text.String()))
	return nil
}

func genModelSQL(tableName string, st StructInfo, text strings.Builder, rangeField strings.Builder, rangeField1 strings.Builder) error {
	text.WriteString(fmt.Sprintf("type %sSQLRepo struct {\n\tdb *gorm.DB\n}\n", tableName))

	text.WriteString(fmt.Sprintf("func (repo %sSQLRepo) GetByID(id uint) (*%s, error) {\n\tq := %sQuery{\n\t\tID: id,\n\t}\n\treturn repo.get(q)\n}\n\n", tableName, st.Name, tableName))

	text.WriteString(fmt.Sprintf("func (repo %vSQLRepo) Create(param serializer.%vCreateParam) (*%v, error) {\n", tableName, st.Name, st.Name))
	text.WriteString(fmt.Sprintf(`	obj := &%v{
%v
	}
	if err := repo.db.Create(obj).Error; err != nil {
		return nil, fmt.Errorf("Create %v err:%%v", err)
	}

	return repo.GetByID(obj.ID)
}

`, st.Name, rangeField.String(), st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) Modify(param serializer.%sModifyParam) (*%s, error) {

	query := %sQuery{
		ID: param.ID,
	}
	obj, err := repo.get(query)
	if err != nil {
		return nil, err
	}

%v

	if err := repo.db.Save(obj).Error; err != nil {
		return nil, fmt.Errorf("Modify %v err: %%v", err)
	}
	return obj, nil
}

`, tableName, st.Name, st.Name, tableName, rangeField1.String(), st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) Search(param serializer.%sSearchParam) ([]%s, uint, error) {

	query := %sQuery{
%s
		Limit:        param.PageSize,
		Offset: (param.Page - 1) * param.PageSize,
	}

	count, err := repo.count(query)
	if err != nil {
		return nil, 0, err
	}

	var objArr []%s
	db := repo.db.Scopes(query.where(), query.preload(), query.order()).Offset(int(query.Offset))
	if query.Limit > 0 {
		db = db.Limit(int(query.Limit))
	}

	if err := db.Find(&objArr).Error; err != nil {
		return nil, 0, fmt.Errorf("Search %s err: %%v", err)
	}
	return objArr, count, nil
}

`, tableName, st.Name, st.Name, tableName, rangeField.String(), st.Name, st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) Delete(param serializer.%sDeleteParam) error {

	query := %sQuery{
		ID: param.ID,
	}
	obj, err := repo.get(query)
	if err != nil {
		return err
	}
	if err := repo.db.Delete(obj).Error; err != nil {
		return fmt.Errorf("Delete %s err: %%v", err)
	}
	return nil
}

`, tableName, st.Name, tableName, st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) get(query %sQuery) (*%s, error) {

	obj := &%s{}
	if err := repo.db.Scopes(query.where(), query.preload()).First(obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("get %s err: %v", err)
		}
		return nil, fmt.Errorf("get %s err: %%v", err)
	}
	return obj, nil
}

`, tableName, tableName, st.Name, st.Name, st.Name, "%v", st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) count(query %sQuery) (uint, error) {

		var count int64
		if err := repo.db.Model(&%s{}).Scopes(query.where()).Count(&count).Error; err != nil {
			return 0, fmt.Errorf("count %s err: %%v", err)
		}
		return uint(count), nil
	}
	`, tableName, tableName, st.Name, st.Name))

	f, err := os.OpenFile("generate/model.ex.go", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(text.String()))
	return nil
}
