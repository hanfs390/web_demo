package routes

import (
	"ByzoroAC/routes/v2"
	"github.com/labstack/echo"
)
var version = "v2"
func Init() {
	e := echo.New()
	e.Static("/", "views")
	if version == "v2" {
		v2.RouteInit(e)
	}
	e.Logger.Fatal(e.Start(":1323"))
}