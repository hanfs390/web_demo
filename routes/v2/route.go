package v2

import (
	"github.com/labstack/echo"
)

func routeStatic(e *echo.Echo) {
	e.Static("/static", "static")
}
func RouteInit(e *echo.Echo) {
	routeGroup(e)
	routeAp(e)
	routeWlan(e)
	routeLed(e)
	routePolicy(e)
	routeBWlist(e)
	routeLoadBalance(e)
	routeUpgrade(e)
	routeReport(e)
	routeStatic(e)
}
