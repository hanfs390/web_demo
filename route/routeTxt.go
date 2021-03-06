package route

import (
	"github.com/labstack/echo"
	"web_demo/global"
	"web_demo/readConf"
	"net/http"
	"fmt"
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

/*func getTxtByName(c echo.Context) error {
	name := c.QueryParam("name") //name is md5 string
	if name == "" {
		return c.JSON(http.StatusBadRequest, "name is error")
	}
	fmt.Println("name is ", name)
	for i := 0; i < len(readConf.TxtData); i++ {
		for j := 0; j < len(readConf.TxtData[i].TxtList); j++ {
			if name == readConf.TxtData[i].TxtList[j].FileName {
				fmt.Println("we found it")
				var imgHtml string
				for n := 1; n <= readConf.TxtData[i].TxtList[j].Number; n++ {
					number := strconv.Itoa(n)
					imgHtml = imgHtml + `<img src="`+readConf.TxtData[i].TxtList[j].Url+`/`+ number +`.jpg">`
				}
				return c.HTML(http.StatusOK, `<html><head><title>`+name+`</title></head><body><div>`+imgHtml+`</div></body></html>`)
			}
		}
	}
	return c.JSON(http.StatusOK, "no found")
}*/

func routeTxt(e *echo.Echo) {
	e.GET("/txt/head", getTxtHead) /* get dir list*/
	e.GET("/txt/allTxt", getAllTxtTxt)
	e.GET("/txt/dirTxt", getTxtByDir)
	e.GET("/txt/nameTxt", getTxtByName)
	e.Static(global.RouteTxtDir, global.TxtDir)
}
