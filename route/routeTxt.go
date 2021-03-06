package route

import (
	"github.com/labstack/echo"
	"web_demo/readConf"
	"net/http"
	"fmt"
	"strconv"
)

type txtHead struct {
	Dir string
}

func getTxtHead(c echo.Context) error {
	var head []txtHead
	for i := 0; i < len(readConf.TxtData); i++ {
		headNode := txtHead{}
		headNode.Dir = readConf.TxtData[i].Dir
		head = append(head, headNode)
	}
	return c.JSON(http.StatusOK, head)
}
func getAllTxtTxt(c echo.Context) error {
	var list []readConf.TxtNode
	for i := 0; i < len(readConf.TxtData); i++ {
		for j := 0; j < len(readConf.TxtData[i].TxtList); j++ {
			list = append(list, readConf.TxtData[i].TxtList[j])
		}
	}
	return c.JSON(http.StatusOK, list)
}
func getTxtByDir(c echo.Context) error {
	dir := c.QueryParam("dir")
	if dir == "" {
		return c.JSON(http.StatusBadRequest, "dir is error")
	}
	fmt.Println("dir is ", dir)
	for i := 0; i < len(readConf.TxtData); i++ {
		if dir == readConf.TxtData[i].Dir {
			return c.JSON(http.StatusOK, readConf.TxtData[i].TxtList)
		}
	}
	return c.JSON(http.StatusBadRequest, dir + " is no found")
}

func getTxtByName(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "name is error")
	}
	fmt.Println("name is ", name)
	for i := 0; i < len(readConf.TxtData); i++ {
		for j := 0; j < len(readConf.TxtData[i].TxtList); j++ {
			if name == readConf.TxtData[i].TxtList[j].FileName {
				err := readConf.CreateTxtChapters(&readConf.TxtData[i].TxtList[j])
				if err != nil {
					return c.JSON(http.StatusOK, "no found")
				}
				html := readConf.BuildTxtDirHTML(&readConf.TxtData[i].TxtList[j])
				return c.HTML(http.StatusOK, html)
			}
		}
	}
	return c.JSON(http.StatusOK, "no found")
}
func getChapter(c echo.Context) error {
	name := c.QueryParam("name")
	number := c.QueryParam("chapters") //name is md5 string
	if name == "" || number == "" {
		return c.JSON(http.StatusBadRequest, "name or chapters is error")
	}
	num, err :=  strconv.Atoi(number)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "number is error")
	}
	data := readConf.GetWordinTxt(name, num)
	return c.HTML(http.StatusOK, `<html><head><title>`+name+`</title></head><body><div>`+ data +`</div></body></html>`)
}
func routeTxt(e *echo.Echo) {
	e.GET("/txt/head", getTxtHead) /* get dir list*/
	e.GET("/txt/allTxt", getAllTxtTxt)
	e.GET("/txt/dirTxt", getTxtByDir)
	e.GET("/txt/nameTxt", getTxtByName)
	e.GET("/txt/chapters", getChapter)
	//e.Static(global.RouteTxtDir, global.TxtDir)

}
