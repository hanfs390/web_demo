package api

import (
	. "ByzoroAC/common"
	"ByzoroAC/models/mysql_driver"
	"fmt"
)

/**
 * the cmd of group
 */

func HandleCreateGroup(new TblGroup) error {
	fmt.Println(new.Name)
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdateGroup(id uint64, value map[string]interface{}) error {
	where := TblGroup{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeleteGroup(id uint64) error {
	where := TblGroup{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindGroupById(id uint64) (result TblGroup, err error) {
	where := TblGroup{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblGroup)
	return result, err
}
func HandleFindGroupsByUrl(url string) (result []TblGroup, err error) {
	r, err := mysql_driver.Handle.FindLike(TblGroup{}, "url", url)
	result = r.([]TblGroup)
	return result, err
}
func HandleFindGroupsAll() (result []TblGroup, err error) {
	r, err := mysql_driver.Handle.FindAll(TblGroup{})
	result = r.([]TblGroup)
	return result, err
}