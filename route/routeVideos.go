package route

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"web_demo/global"
	"web_demo/readConf"
)

type videoHead struct {
	Dir string
	Label []string
}

func getVideoHead(c echo.Context) error {
	var head []videoHead
	for i := 0; i < len(readConf.VideoData); i++ {
		headNode := videoHead{}
		headNode.Dir = readConf.VideoData[i].Dir
		for j := 0; j < len(readConf.VideoData[i].Label); j++ {
			headNode.Label = append(headNode.Label, readConf.VideoData[i].Label[j])
		}
		head = append(head, headNode)
	}
	return c.JSON(http.StatusOK, head)
}
func getAllVideo(c echo.Context) error {
	var list []readConf.VideoNode
	for i := 0; i < len(readConf.VideoData); i++ {
		for j := 0; j < len(readConf.VideoData[i].VideoList); j++ {
			list = append(list, readConf.VideoData[i].VideoList[j])
		}
	}
	return c.JSON(http.StatusOK, list)
}
func getVideoByDir(c echo.Context) error {
	dir := c.QueryParam("dir")
	if dir == "" {
		return c.JSON(http.StatusBadRequest, "dir is error")
	}
	fmt.Println("dir is ", dir)
	for i := 0; i < len(readConf.VideoData); i++ {
		if dir == readConf.VideoData[i].Dir {
			return c.JSON(http.StatusOK, readConf.VideoData[i].VideoList)
		}
	}
	return c.JSON(http.StatusBadRequest, dir + " is no found")
}
func getVideoByLabel(c echo.Context) error {
	label := c.QueryParam("label")
	if label == "" {
		return c.JSON(http.StatusBadRequest, "label is error")
	}
	var list []readConf.VideoNode
	for i := 0; i < len(readConf.VideoData); i++ {
		for j := 0; j < len(readConf.VideoData[i].VideoList); j++ {
			if label == readConf.VideoData[i].VideoList[j].Label {
				list = append(list, readConf.VideoData[i].VideoList[j])
			}
		}
	}
	return c.JSON(http.StatusOK, list)
}
func routeVideo(e *echo.Echo) {
	e.GET("/video/head", getVideoHead)
	e.GET("/video/allVideo", getAllVideo)
	e.GET("/video/dirVideo", getVideoByDir)
	e.GET("/video/labelVideo", getVideoByLabel)
	e.Static(global.RouteVideoDir, global.VideoDir)
}