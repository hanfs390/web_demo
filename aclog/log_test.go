package aclog

import (
	"testing"
	"ByzoroAC/conf"
	"fmt"
)

func TestMain(m *testing.M) {
	conf.ReadGlobalConf()
	Init()
	m.Run()
}

func TestPrint(t *testing.T) {
	a := "Hello World"
	Debug("hanfushun %d %s", 1, a)
	fmt.Println(a)
	Info("hanfushun %d %s", 2, a)
	Warning("hanfushun %d %s", 3, a)
	Error("hanfushun %d %s", 4, a)
}