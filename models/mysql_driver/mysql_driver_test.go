package mysql_driver

import (
	"testing"
	"ByzoroAC/conf"
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"fmt"
)

func TestMain(m *testing.M) {
	aclog.Init()
	conf.ReadGlobalConf()
	ModuleInit()
	m.Run()
}

func tblGroupInsert(t *testing.T) {
	a := common.TblGroup{Name:"default",Uuid:"12341415"}
	insert(a)
	a.Name = "test"
	insert(a)
}
func tblGroupDelete(t *testing.T) {
	a := common.TblGroup{Id:2}
	delete(a)
}
func tblGroupUpdate(t *testing.T) {
	a := common.TblGroup{}
	m := make(map[string]interface{})
	m["name"] = "han"
	update(a, m)
}

func tblSelectOneLine (t *testing.T) {
	w := common.TblGroup{Name:"han"}
	r, err := findFirst(w)
	if err != nil {
		fmt.Println(err)
	}
	group := r.(common.TblGroup)
	fmt.Println(group.Id)
}
func tblSelectLines(t *testing.T) {
	w := common.TblGroup{Name:"han"}
	r, err := findAll(w)
	if err != nil {
		fmt.Println(err)
	}
	groups := r.([]common.TblGroup)
	for i := 0; i < len(groups); i++ {
		fmt.Println(groups[i].Id)
	}
}
func tblSelectVagueLines(t *testing.T) {
	r, err := findLike(common.TblGroup{}, "uuid", "1%")
	if err != nil {
		fmt.Println(err)
	}
	groups := r.([]common.TblGroup)
	for i := 0; i < len(groups); i++ {
		fmt.Println(groups[i].Id)
	}
}
func TestAll(t *testing.T) {
	t.Run("insert", tblGroupInsert)
	t.Run("delete", tblGroupDelete)
	t.Run("update", tblGroupUpdate)
	t.Run("search specify a line", tblSelectOneLine)
	t.Run("search specify lines", tblSelectLines)
	t.Run("search vague lines", tblSelectVagueLines)
}

