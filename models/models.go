package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// 声明 db 全局变量
var db *gorm.DB

// 生明 Model 结构体, 用于继承, 使其他结构体拥有公共字段
// 相当于修改了 gorm.Model 结构体
type Model struct {
	// id 主键
	ID int `gorm:"primary_key" json:"id"`

	// 创建时间
	CreatedOn int `json:"created_on"`

	// 修改时间
	ModifiedOn int `json:"modified_on"`

	// 软删除
	DeletedOn int `json:"deleted_on"`
}

// Setup initializes the database instance
func Setup() {
	var err error

	// connect to db
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	// 设置表名 前缀 + 表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	// 设置表名为复数形式
	db.SingularTable(true)

	// 设置 create, update, delete 回调函数
	// 用于设置 CreatedOn, ModifiedOn, DeletedOn 字段
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	// 设置最大连接数
	db.DB().SetMaxIdleConns(10)

	// 设置最大打开连接数
	db.DB().SetMaxOpenConns(100)
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
// 用于设置 CreatedOn, ModifiedOn 字段, 用于记录创建时间, 修改时间
func updateTimeStampForCreateCallback(scope *gorm.Scope) {

	if !scope.HasError() {

		// 获取 现在的时间
		nowTime := time.Now().Unix()

		// 判断是否存在 CreatedOn 字段
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				// 设置 CreatedOn 字段
				createTimeField.Set(nowTime)
			}
		}

		// 判断是否存在 ModifiedOn 字段
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				// 设置 ModifiedOn 字段
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
// 用于设置 ModifiedOn 字段, 用于记录修改时间
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {

	// 判断是否存在 ModifiedOn 字段
	if _, ok := scope.Get("gorm:update_column"); !ok {

		// 获取 现在的时间, 并设置 ModifiedOn 字段
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// deleteCallback will set `DeletedOn` where deleting
// 用于设置 DeletedOn 字段, 用于记录删除时间
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string

		// 判断是否存在 delete_option 字段
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {

			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
// 用于添加空格
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
