package mysql_driver

import (
	"fmt"
	"ByzoroAC/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "ByzoroAC/common"
)

type localConf struct {
	mySQLUserName string
	mySQLPassword string
	mySQLHost string
	mySQLPort string
	mySQLDataBase string
	mySQLCharset string
}
type mySQLHandles struct {
	Insert func(data interface{}) error
	Delete func (data interface{}) error
	Update func (where interface{}, value map[string]interface{}) error
	FindFirst func (where interface{}) (result interface{}, err error)
	FindAll func (where interface{}) (result interface{}, err error)
	FindAllPage func (where interface{}, pageSize int, page int) (result interface{}, err error)
	FindLike func (table interface{}, key string, value string) (result interface{}, err error)
	FindAllCount func (where interface{}) (count int, err error)
}
var sqlConf localConf
var db *gorm.DB
var Handle mySQLHandles
/*
	List of Tables
		*TblGroup
		*TblDevInfo
		*TblDevStat
		*TblWlan
		*TblLED
		*TblLoadBalanceGroup
		*TblPolicy
		*TblPolicy
		*TblDevGPONStat
*/
/**
 * insert - insert a new line to table
 *
 * @data: the table struct include the data that want to insert.
 */
func insert(data interface{}) error {
	switch data.(type) {
	case TblGroup:
		temp := data.(TblGroup)
		fmt.Println(temp.Name)
		return db.Create(&temp).Error
	case TblDevInfo:
		temp := data.(TblDevInfo)
		return db.Create(&temp).Error
	case TblDevStat:
		temp := data.(TblDevStat)
		return db.Create(&temp).Error
	case TblWlan:
		temp := data.(TblWlan)
		return db.Create(&temp).Error
	case TblLED:
		temp := data.(TblLED)
		return db.Create(&temp).Error
	case TblLoadBalanceGroup:
		temp := data.(TblLoadBalanceGroup)
		return db.Create(&temp).Error
	case TblPolicy:
		temp := data.(TblPolicy)
		return db.Create(&temp).Error
	case TblBWList:
		temp := data.(TblBWList)
		return db.Create(&temp).Error
	default:
		return  nil
	}
}
/**
 * @data: the table struct include id.
 */
func delete(data interface{}) error {
	switch data.(type) {
	case TblGroup:
		temp := data.(TblGroup)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblDevInfo:
		temp := data.(TblDevInfo)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblDevStat:
		temp := data.(TblDevStat)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblWlan:
		temp := data.(TblWlan)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblLED:
		temp := data.(TblLED)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblLoadBalanceGroup:
		temp := data.(TblLoadBalanceGroup)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblPolicy:
		temp := data.(TblPolicy)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	case TblBWList:
		temp := data.(TblBWList)
		if temp.Id <= 0 {
			return nil
		}
		return db.Delete(&temp).Error
	default:
		return  nil
	}
}
/**
 * @where: has to exist
 */
func update(where interface{}, value map[string]interface{}) error {
	switch where.(type) {
	case TblGroup:
		a := where.(TblGroup)
		b := TblGroup{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblDevInfo:
		a := where.(TblDevInfo)
		b := TblDevInfo{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblDevStat:
		a := where.(TblDevStat)
		b := TblDevStat{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblWlan:
		a := where.(TblWlan)
		b := TblWlan{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblLED:
		a := where.(TblLED)
		b := TblLED{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblLoadBalanceGroup:
		a := where.(TblLoadBalanceGroup)
		b := TblLoadBalanceGroup{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblPolicy:
		a := where.(TblPolicy)
		b := TblPolicy{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	case TblBWList:
		a := where.(TblBWList)
		b := TblBWList{}
		if a == b {
			fmt.Println("no where")
			return nil
		}
		return db.Model(&a).Where(a).Updates(value).Error
	default:
		return  nil
	}
}
/**
 * findFirst - find the first lines match where. get all lines if where is empty
 */
func findFirst(where interface{}) (result interface{}, err error) {
	switch where.(type) {
	case TblGroup:
		var r TblGroup
		a := where.(TblGroup)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblDevInfo:
		var r TblDevInfo
		a := where.(TblDevInfo)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblDevStat:
		var r TblDevStat
		a := where.(TblDevStat)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblWlan:
		var r TblWlan
		a := where.(TblWlan)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblLED:
		var r TblLED
		a := where.(TblLED)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblLoadBalanceGroup:
		var r TblLoadBalanceGroup
		a := where.(TblLoadBalanceGroup)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblPolicy:
		var r TblPolicy
		a := where.(TblPolicy)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	case TblBWList:
		var r TblBWList
		a := where.(TblBWList)
		err := db.Model(&a).Where(a).First(&r).Error
		return r, err
	default:
		return  nil, nil
	}
}
func findAll(where interface{}) (result interface{}, err error) {
	switch where.(type) {
	case TblGroup:
		var r []TblGroup
		a := where.(TblGroup)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblDevInfo:
		var r []TblDevInfo
		a := where.(TblDevInfo)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblDevStat:
		var r []TblDevStat
		a := where.(TblDevStat)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblWlan:
		var r []TblWlan
		a := where.(TblWlan)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblLED:
		var r []TblLED
		a := where.(TblLED)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblLoadBalanceGroup:
		var r []TblLoadBalanceGroup
		a := where.(TblLoadBalanceGroup)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblPolicy:
		var r []TblPolicy
		a := where.(TblPolicy)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	case TblBWList:
		var r []TblBWList
		a := where.(TblBWList)
		err := db.Model(&a).Where(a).Find(&r).Error
		return r, err
	default:
		return  nil, nil
	}
}
func findAllCount(where interface{}) (count int, err error) {
	switch where.(type) {
	case TblGroup:
		a := where.(TblGroup)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblDevInfo:
		a := where.(TblDevInfo)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblDevStat:
		a := where.(TblDevStat)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblWlan:
		a := where.(TblWlan)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblLED:
		a := where.(TblLED)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblLoadBalanceGroup:
		a := where.(TblLoadBalanceGroup)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblPolicy:
		a := where.(TblPolicy)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	case TblBWList:
		a := where.(TblBWList)
		err := db.Model(&a).Where(a).Count(&count).Error
		return count, err
	default:
		return  0, nil
	}
}
/**
 * findAllPage - find a page from database
 *
 * @where: the condition to mach
 * @pageSize: the size of page
 * @page: which page do you want. it is >=0
 */
func findAllPage(where interface{}, pageSize int, page int) (result interface{}, err error) {
	switch where.(type) {
	case TblGroup:
		var r []TblGroup
		a := where.(TblGroup)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblDevInfo:
		var r []TblDevInfo
		a := where.(TblDevInfo)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblDevStat:
		var r []TblDevStat
		a := where.(TblDevStat)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblWlan:
		var r []TblWlan
		a := where.(TblWlan)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblLED:
		var r []TblLED
		a := where.(TblLED)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblLoadBalanceGroup:
		var r []TblLoadBalanceGroup
		a := where.(TblLoadBalanceGroup)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblPolicy:
		var r []TblPolicy
		a := where.(TblPolicy)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	case TblBWList:
		var r []TblBWList
		a := where.(TblBWList)
		err := db.Model(&a).Where(a).Limit(pageSize).Offset(page*pageSize).Find(&r).Error
		return r, err
	default:
		return  nil, nil
	}
}
/**
 * findLike - find the results by a vague condition
 *
 * @table: which table do you want to find in database. it is can be empty
 * @key: one column of table
 * @value: the value of key. example. "test%", "%test%", "%test"
 */
func findLike(table interface{}, key string, value string) (result interface{}, err error) {
	switch table.(type) {
	case TblGroup:
		var r []TblGroup
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblDevInfo:
		var r []TblDevInfo
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblDevStat:
		var r []TblDevStat
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblWlan:
		var r []TblWlan
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblLED:
		var r []TblLED
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblLoadBalanceGroup:
		var r []TblLoadBalanceGroup
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblPolicy:
		var r []TblPolicy
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	case TblBWList:
		var r []TblBWList
		k := key + " like ? "
		err = db.Where(k, value).Find(&r).Error
		return r, err
	default:
		return  nil, nil
	}
}
func checkTables() {
	if !db.HasTable(&TblGroup{}) {
		db.CreateTable(&TblGroup{})
		defaultGroup := TblGroup{Name:"default", Id:1, Url:"/root"}
		insert(defaultGroup)
	}
	if !db.HasTable(&TblDevInfo{}) {
		db.CreateTable(&TblDevInfo{})
	}
	if !db.HasTable(&TblDevStat{}) {
		db.CreateTable(&TblDevStat{})
	}
	if !db.HasTable(&TblDevInfo{}) {
		db.CreateTable(&TblDevInfo{})
	}
	if !db.HasTable(&TblWlan{}) {
		db.CreateTable(&TblWlan{})
	}
	if !db.HasTable(&TblLED{}) {
		db.CreateTable(&TblLED{})
	}
	if !db.HasTable(&TblLoadBalanceGroup{}) {
		db.CreateTable(&TblLoadBalanceGroup{})
	}
	if !db.HasTable(&TblPolicy{}) {
		db.CreateTable(&TblPolicy{})
	}
	if !db.HasTable(&TblBWList{}) {
		db.CreateTable(&TblBWList{})
	}
}
func initHandles() {
	Handle.Insert = insert
	Handle.Delete = delete
	Handle.Update = update
	Handle.FindFirst = findFirst
	Handle.FindAll = findAll
	Handle.FindAllPage = findAllPage
	Handle.FindLike = findLike
	Handle.FindAllCount = findAllCount
}
func ModuleInit() int {
	// ready
	if conf.GlobalConf["MySQLUserName"] != "" {
		sqlConf.mySQLUserName = conf.GlobalConf["MySQLUserName"]
	} else {
		sqlConf.mySQLUserName = "root"
	}
	if conf.GlobalConf["MySQLPassword"] != "" {
		sqlConf.mySQLPassword = conf.GlobalConf["MySQLPassword"]
	} else {
		sqlConf.mySQLPassword = "12345678"
	}
	if conf.GlobalConf["MySQLHost"] != "" {
		sqlConf.mySQLHost = conf.GlobalConf["MySQLHost"]
	} else {
		sqlConf.mySQLHost = "127.0.0.1"
	}
	if conf.GlobalConf["MySQLPort"] != "" {
		sqlConf.mySQLPort = conf.GlobalConf["MySQLPort"]
	} else {
		sqlConf.mySQLPort = "3306"
	}
	if conf.GlobalConf["MySQLDataBase"] != "" {
		sqlConf.mySQLDataBase = conf.GlobalConf["MySQLDataBase"]
	} else {
		sqlConf.mySQLDataBase = "ac"
	}
	if conf.GlobalConf["MySQLCharset"] != "" {
		sqlConf.mySQLCharset = conf.GlobalConf["MySQLCharset"]
	} else {
		sqlConf.mySQLCharset = "utf8"
	}
	// root:123456@tcp(localhost:3306)/ac?charset=utf8
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", sqlConf.mySQLUserName, sqlConf.mySQLPassword,
		sqlConf.mySQLHost, sqlConf.mySQLPort, sqlConf.mySQLDataBase, sqlConf.mySQLCharset)
	fmt.Println(dbDSN)

	// 打开连接失败
	var err error
	db, err = gorm.Open("mysql", dbDSN)

	if err != nil {
		fmt.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + err.Error())
	}
	db.DB().SetConnMaxLifetime(0)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(20)
	db.SingularTable(true)
	checkTables()
	initHandles()
	return 0;
}

