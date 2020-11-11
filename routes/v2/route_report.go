package v2

import (
	"github.com/labstack/echo"
	"net/http"
	"fmt"
)
type grafanaConf struct {
	Url string
}
func getGrafanaURL(c echo.Context) error {
	fmt.Println("grafana")
	return c.JSON(http.StatusOK, grafanaConf{Url:"http://192.168.100.7:3000/d/hXnpxsHGz/globalinfo?orgId=1"})
}
func routeReport(e *echo.Echo) {
	e.GET("/grafanaURL", getGrafanaURL)
}
