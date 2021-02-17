package main

import (
	"web_demo/readConf"
	"web_demo/route"
	"fmt"
)

func main() {
	fmt.Println("My Web start")
	readConf.ReadAllFilesInfo()
	fmt.Println("Read files finished")
	route.InitRoutes()
}