package route

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"web_demo/global"
	"web_demo/readConf"
)
type pictureHead struct {
	Dir string
	Label []string
}

func getPictureHead(c echo.Context) error {
	var head []pictureHead
	for i := 0; i < len(readConf.PictureData); i++ {
		headNode := pictureHead{}
		headNode.Dir = readConf.PictureData[i].Dir
		for j := 0; j < len(readConf.PictureData[i].Label); j++ {
			headNode.Label = append(headNode.Label, readConf.PictureData[i].Label[j])
		}
		head = append(head, headNode)
	}
	return c.JSON(http.StatusOK, head)
}
func getAllPicturePicture(c echo.Context) error {
	var list []readConf.PictureNode
	for i := 0; i < len(readConf.PictureData); i++ {
		for j := 0; j < len(readConf.PictureData[i].PictureList); j++ {
			list = append(list, readConf.PictureData[i].PictureList[j])
		}
	}
	return c.JSON(http.StatusOK, list)
}
func getPictureByDir(c echo.Context) error {
	dir := c.QueryParam("dir")
	if dir == "" {
		return c.JSON(http.StatusBadRequest, "dir is error")
	}
	fmt.Println("dir is ", dir)
	for i := 0; i < len(readConf.PictureData); i++ {
		if dir == readConf.PictureData[i].Dir {
			return c.JSON(http.StatusOK, readConf.PictureData[i].PictureList)
		}
	}
	return c.JSON(http.StatusBadRequest, dir + " is no found")
}
func getPictureByLabel(c echo.Context) error {
	label := c.QueryParam("label")
	if label == "" {
		return c.JSON(http.StatusBadRequest, "label is error")
	}
	var list []readConf.PictureNode
	for i := 0; i < len(readConf.PictureData); i++ {
		for j := 0; j < len(readConf.PictureData[i].PictureList); j++ {
			if label == readConf.PictureData[i].PictureList[j].Label {
				list = append(list, readConf.PictureData[i].PictureList[j])
			}
		}
	}
	return c.JSON(http.StatusOK, list)
}
func getPictureByName(c echo.Context) error {
	name := c.QueryParam("name") //name is md5 string
	if name == "" {
		return c.JSON(http.StatusBadRequest, "name is error")
	}
	fmt.Println("name is ", name)
	for i := 0; i < len(readConf.PictureData); i++ {
		for j := 0; j < len(readConf.PictureData[i].PictureList); j++ {
			if name == readConf.PictureData[i].PictureList[j].FileName {
				fmt.Println("we found it")
				var imgHtml string
				for n := 1; n <= readConf.PictureData[i].PictureList[j].Number; n++ {
					number := strconv.Itoa(n)
					imgHtml = imgHtml + `<img src="`+readConf.PictureData[i].PictureList[j].Url+`/`+ number +`.jpg">`
				}
				return c.HTML(http.StatusOK, `<html><head><title>`+name+`</title></head><body><div>`+imgHtml+`</div></body></html>`)
			}
		}
	}
	return c.JSON(http.StatusOK, "no found")
}
func routePicture(e *echo.Echo) {
	e.GET("/picture/head", getPictureHead)
	e.GET("/picture/allPicture", getAllPicturePicture)
	e.GET("/picture/dirPicture", getPictureByDir)
	e.GET("/picture/labelPicture", getPictureByLabel)
	e.GET("/picture/namePicture", getPictureByName)
	e.Static(global.RoutePictureDir, global.PictureDir)
}
