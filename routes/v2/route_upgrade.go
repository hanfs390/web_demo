package v2

import (
	"github.com/labstack/echo"
	"net/http"
	"ByzoroAC/controllers/api"
	"fmt"
)
type upgradeInfo struct {
	Mac string
	Model string
	Version string
}
func uploadUpgradeFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	err = api.SaveUpgradeFile(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "OK")
}
func getNewestVersion(c echo.Context) error {
	model := c.QueryParam("model")
	newest := api.GetNewestUpgradeFilesName(model)
	if newest == "" {
		return c.JSON(http.StatusBadRequest, "No need version:" + model)
	}
	return c.JSON(http.StatusOK, newest)
}
func getAllVersion(c echo.Context) error {
	model := c.QueryParam("model")
	list := api.GetUpgradeFilesName(model)
	if list == nil {
		return c.JSON(http.StatusBadRequest, "No need versions:" + model)
	}
	return c.JSON(http.StatusOK, list)
}
func deleteVersion(c echo.Context) error {
	model := c.QueryParam("model")
	version := c.QueryParam("version")
	err := api.RemoveUpgradeFile(model, version)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "OK")
}
func upgradeThisAp(c echo.Context) error {
	info := upgradeInfo{}
	if err := c.Bind(&info); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err := api.UpgradeSingleAp(info.Mac, info.Model, info.Version)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "OK")
}
func routeUpgrade(e *echo.Echo) {
	e.POST("/upgrade", uploadUpgradeFile) /* upload the version file */
	e.GET("/upgrade/newest", getNewestVersion) /* get the model newest file name */
	e.GET("/upgrade", getAllVersion) /* get all version file name */
	e.DELETE("/upgrade", deleteVersion) /* del the file */
	e.POST("/apupgrade", upgradeThisAp)
	e.POST("/batchupgrade", nil)
}
