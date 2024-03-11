package template

var ModelTemplate = `package ${TableName}

type ${Name} struct {
	${TableBody}
}

func (${Name}) TableName() string {
	return "${TableName}"
}

func InitTable() {
	table := &${Name}{}
	err := db.DB().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&${Name}{})
	if err != nil {
		return
	}
	db.DB().Exec("ALTER TABLE " + table.TableName() + " COMMENT '${Name}';")
}

func Get(id uint64) (*${Name}, error) {
	var o = &${Name}{}
	err := db.DB().Model(${Name}{}).Where("id = ?", id).Find(&o).Error
	return o, err
}

// 检测名称是否已存在
func NameIsExist(name string) (bool, error) {
	var c int64
	err := db.DB().Model(${Name}{}).Where("name=?", name).Count(&c).Error
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

// 分页查询
func GetList(param *List${Name}Param) (res []*${Name}, total int64, err error) {
	model := db.Model(&${Name}{}).Scopes(param.where(), param.order(), param.preload())

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

func Create(param *Create${Name}Param) (*${Name}, error) {
	data := &${Name}{
	${CreateBody}
	}

	err := db.DB().Create(data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Update(param *Update${Name}Param) error {
	if err := db.DB().Where("id = ?", param.Id).Updates(&${Name}{
	${CreateBody}
	}).Error; err != nil {
		return err
	}

	return nil
}

func Delete(param *Delete${Name}Param) (err error) {
	return db.DB().Where("id = ?", param.Id).Delete(&${Name}{}).Error
}
`

var ParamTemplate = `package ${TableName}

import "gorm.io/gorm"

type (
	List${Name}Param struct {
		${ParamBody}
		${Page}
	}

	List${Name}Resp struct {
		${RespData}
	}

	Create${Name}Param struct {
		${ParamBody}
	}

	Update${Name}Param struct {
		${ID}
		${ParamBody}
	}

	Delete${Name}Param struct {
		${ID}
	}
)

// where 条件
func (c *List${Name}Param) where() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		${whereBody}
		return db
	}
}

// order 处理排序规则
func (c *List${Name}Param) order() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order("created_at ASC")
		return db
	}
}

// preload 预加载表
func (c *List${Name}Param) preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// if c.PreloadProperty {
		// 	db = db.Preload("PropertyArr")
		// }
		return db
	}
}
`
