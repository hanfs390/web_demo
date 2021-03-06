package route

import "github.com/labstack/echo"

func InitRoutes() {
	e := echo.New()
	//routePicture(e)
	routeTxt(e)
	//routeVideo(e)
	e.Static("/", "static")
	e.Logger.Fatal(e.Start(":1323"))
}