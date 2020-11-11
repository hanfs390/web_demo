package main

import (
	"web_demo/readConf"
	"web_demo/route"
)

func main() {
	readConf.ReadAllFilesInfo()
	route.InitRoutes()
}