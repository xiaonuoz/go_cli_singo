package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateModelCode(structType []StructInfo) error {
	var text strings.Builder
	for _, st := range structType {

		// 遍历字段
		var rangeField strings.Builder
		var rangeField1 strings.Builder
		var rangeField2 strings.Builder

		var rangeFieldWhere strings.Builder

		for index, field := range st.Field {
			rangeField.WriteString(fmt.Sprintf("\t\t%v:  param.%v,\n", field, field))
			rangeField2.WriteString(fmt.Sprintf("\t%v\t%v\t%v\n", field, st.FieldType[index], st.Comments[index]))
			switch st.FieldType[index] {
			case "string":
				rangeField1.WriteString(fmt.Sprintf(`
	if len(param.%v) > 0 {
		obj.%v = param.%v
	}				
`, field, field, field))
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
				rangeField1.WriteString(fmt.Sprintf(`
	if param.%v > 0 {
		obj.%v = param.%v
	}				
`, field, field, field))
			}
		}

		// 创建ormPath文件
		ormPath := filepath.Join(ProjectDir, "model", "orm")
		err := os.MkdirAll(ormPath, 0755)
		if err != nil && os.IsExist(err) {
			return err
		}

		fs, err := os.Create(filepath.Join(ormPath, "type_orm.go"))
		if err != nil {
			return fmt.Errorf("create orm file err:%v", err)
		}
		defer fs.Close()
		fs.WriteString(fmt.Sprintf(`package orm

import (
	"time"

	"gorm.io/gorm"
)

// GormModel base model
type GormModel struct {
	ID        uint           %s
	CreatedAt time.Time      %s
	UpdatedAt time.Time      %s
	DeletedAt gorm.DeletedAt %s
}
		`, "`json:\"id\" gorm:\"primaryKey\"`", "`json:\"created_at\" gorm:\"index\"`", "`json:\"updated_at\"`", "`json:\"-\" gorm:\"index\"`"))

		// 创建文件夹
		path := filepath.Join(ProjectDir, "model", st.TableName)
		err = os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}

		// 创建查询文件
		f, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("query_%s.go", st.TableName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("OpenFile err:%v", err)
		}
		defer f.Close()

		query := genModelQuery(text, st.TableName, rangeField2, rangeFieldWhere)
		f.Write([]byte(query.String()))

		f1, err := os.OpenFile(filepath.Join(path, fmt.Sprintf("sql_repo_%s.go", st.TableName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("OpenFile err:%v", err)
		}
		defer f1.Close()

		sql := genModelSQL(st, text, rangeField, rangeField1)
		f1.Write([]byte(sql.String()))

		// 创建init
		initPath := filepath.Join(path, "init.go")
		_, err = os.Stat(initPath)
		if os.IsNotExist(err) {
			fi, err := os.Create(initPath)
			if err != nil {
				return fmt.Errorf("create init file err:%v", err)
			}
			defer fi.Close()

			fi.WriteString(fmt.Sprintf(`package %s

var %sRepo %sSQLRepo

func Init(db *gorm.DB) {
	%sRepo = %sSQLRepo{
		db: db,
	}
}

`, st.TableName, st.Name, st.LocalName, st.Name, st.LocalName))
		}
	}
	return nil
}

func genModelQuery(text strings.Builder, tableName string, rangeField2 strings.Builder, rangeFieldWhere strings.Builder) strings.Builder {
	text.WriteString(fmt.Sprintf(`package %s

import "gorm.io/gorm"

// %sQuery %s查询条件
type %sQuery struct {
	ID           uint
%v

	Limit  uint
	Offset uint
}

`, tableName, tableName, tableName, tableName, rangeField2.String()))

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

	return text
}

func genModelSQL(st StructInfo, text strings.Builder, rangeField strings.Builder, rangeField3 strings.Builder) strings.Builder {
	text.WriteString(fmt.Sprintf("package %s\n\n", st.TableName))

	text.WriteString(fmt.Sprintf("type %sSQLRepo struct {\n\tdb *gorm.DB\n}\n\n", st.LocalName))

	text.WriteString(fmt.Sprintf("func (repo %sSQLRepo) GetByID(id uint) (*%s, error) {\n\tq := %sQuery{\n\t\tID: id,\n\t}\n\treturn repo.get(q)\n}\n\n", st.LocalName, st.Name, st.TableName))

	text.WriteString(fmt.Sprintf("func (repo %vSQLRepo) Create(param *serializer.%vCreateParam) error {\n", st.LocalName, st.Name))
	text.WriteString(fmt.Sprintf(`	obj := &%v{
%v
	}
	if err := repo.db.Create(obj).Error; err != nil {
		return fmt.Errorf("Create %v err:%%v", err)
	}

	return nil
}

`, st.Name, rangeField.String(), st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) Modify(param *serializer.%sModifyParam) (*%s, error) {

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

`, st.LocalName, st.Name, st.Name, st.TableName, rangeField3.String(), st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) Search(param *serializer.%sSearchParam) ([]%s, uint, error) {

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

`, st.LocalName, st.Name, st.Name, st.TableName, rangeField.String(), st.Name, st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) Delete(param *serializer.%sDeleteParam) error {

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

`, st.LocalName, st.Name, st.TableName, st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) List(param *serializer.%sListParam) ([]%s, error) {
		query := %sQuery{
	%s
		}
	
		var objArr []%s
		db := repo.db.Scopes(query.where(), query.preload(), query.order())
	
		if err := db.Find(&objArr).Error; err != nil {
			return nil, fmt.Errorf("List %s err: %%v", err)
		}
		return objArr, nil
	}
	
	`, st.LocalName, st.Name, st.Name, st.TableName, rangeField.String(), st.Name, st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) get(query %sQuery) (*%s, error) {

	obj := &%s{}
	if err := repo.db.Scopes(query.where(), query.preload()).First(obj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil,nil
		}
		return nil, fmt.Errorf("get %s err: %%v", err)
	}
	return obj, nil
}

`, st.LocalName, st.TableName, st.Name, st.Name, st.Name))

	text.WriteString(fmt.Sprintf(`func (repo %sSQLRepo) count(query %sQuery) (uint, error) {

		var count int64
		if err := repo.db.Model(&%s{}).Scopes(query.where()).Count(&count).Error; err != nil {
			return 0, fmt.Errorf("count %s err: %%v", err)
		}
		return uint(count), nil
	}
	`, st.LocalName, st.TableName, st.Name, st.Name))

	return text
}
