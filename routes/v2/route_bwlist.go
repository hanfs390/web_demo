package v2

import (
	"github.com/labstack/echo"
	"ByzoroAC/controllers/api"
	"fmt"
	"net/http"
	"strconv"
	"ByzoroAC/common"
	"ByzoroAC/models/redis_driver"
	"ByzoroAC/controllers/task"
	"time"
	"encoding/json"
)
/*
 * the bwlist msg that send to AP
 */
type bwlistMsg struct {
	Type int
	List string
}

func getAllBWlistConfs(c echo.Context) error {
	result, err := api.HandleFindBWlistsAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
func getBWlistConfById(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "BWlist ID is error")
	}
	result, err := api.HandleFindBWlistById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func createBWlistConf(c echo.Context) error {
	bwlist := common.TblBWList{}
	if err := c.Bind(&bwlist); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if bwlist.GroupId <= 1 { /* default group has not bwlist */
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}
	if bwlist.Type == "black" {

	} else if bwlist.Type == "white" {

	} else {
		return c.JSON(http.StatusBadRequest, "Black/White list type is error:" + bwlist.Type)
	}
	/* get the group info by group id */
	r, err := api.HandleFindGroupById(bwlist.GroupId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	bwlist.GroupId = r.Id
	bwlist.GroupName = r.Name
	bwlist.GroupUrl = r.Url

	/* checkout the bwlist exist or not in this group */
	temp, err := api.HandleFindPolicyByGroupUrl(bwlist.GroupUrl)
	if len(temp) > 0 {
		return c.JSON(http.StatusBadRequest, "bwlist has exist in this group"+bwlist.GroupName)
	}

	/* insert */
	err = api.HandleCreateBWlist(bwlist)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var data bwlistMsg
	if bwlist.Type == "black" {
		data.List = bwlist.BlackList
		data.Type = 0
	} else if bwlist.Type == "white" {
		data.List = bwlist.WhiteList
		data.Type = 1
	}
	value, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
	}
	err = redis_driver.RedisDb.Set("bwlist_"+bwlist.GroupUrl, string(value))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"bwlist_"+bwlist.GroupUrl, Value:string(value)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("bwlist_"+bwlist.GroupUrl+"_MD5", getMD5(string(value)))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"bwlist_"+bwlist.GroupUrl+"_MD5", Value:getMD5(string(value))}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}

func updateBWlistConf(c echo.Context) error {
	bwlist := common.TblBWList{}
	value := make(map[string]interface{})
	if err := c.Bind(&bwlist); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if bwlist.Id <= 0 {
		return c.JSON(http.StatusBadRequest, "BWlist ID is error")
	}
	/* verify elements */
	if bwlist.GroupId > 1 {
		r, err := api.HandleFindGroupById(bwlist.GroupId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		value["group_id"] = r.Id
		value["group_name"] = r.Name
		value["group_url"] = r.Url
	}
	/* checkout the bwlist exist or not in this group */
	temp, err := api.HandleFindPolicyByGroupUrl(value["group_url"].(string))
	if len(temp) > 0 {
		return c.JSON(http.StatusBadRequest, "bwlist has exist in this group"+value["group_name"].(string))
	}

	if bwlist.Type == "black" {
		value["type"] = "black"
	} else if bwlist.Type == "white" {
		value["type"] = "white"
	} else {
		return c.JSON(http.StatusBadRequest, "Black/White list type is error:" + bwlist.Type)
	}

	if len(bwlist.BlackList) <= 10240 {
		value["black_list"] = bwlist.BlackList
	}
	if len(bwlist.WhiteList) <= 10240 {
		value["white_list"] = bwlist.WhiteList
	}

	/* update to mySQL */
	err = api.HandleUpdateBWlist(bwlist.Id, value)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var data bwlistMsg
	if bwlist.Type == "black" {
		data.List = bwlist.BlackList
		data.Type = 0
	} else if bwlist.Type == "white" {
		data.List = bwlist.WhiteList
		data.Type = 1
	}
	j, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
	}
	err = redis_driver.RedisDb.Set("bwlist_"+bwlist.GroupUrl, string(j))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"bwlist_"+bwlist.GroupUrl, Value:string(j)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("bwlist_"+bwlist.GroupUrl+"_MD5", getMD5(string(j)))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"bwlist_"+bwlist.GroupUrl+"_MD5", Value:getMD5(string(j))}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}
func deleteBWlistConf(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "BWlist ID is error")
	}
	bwlist, err := api.HandleFindBWlistById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = api.HandleDeleteBWlist(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* del redis */
	err = redis_driver.RedisDb.DeleteKey("bwlist_"+bwlist.GroupUrl)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"bwlist_"+bwlist.GroupUrl}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("bwlist_"+bwlist.GroupUrl+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"bwlist_"+bwlist.GroupUrl+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return c.JSON(http.StatusOK, "OK")
}

func routeBWlist(e *echo.Echo) {
	e.GET("/bwlist", getAllBWlistConfs) /* get the all bwlist info */
	e.GET("/bwlist/id", getBWlistConfById) /* get the bwlist by id */
	e.POST("/bwlist", createBWlistConf) /* create a new bwlist, reserve */
	e.PUT("/bwlist", updateBWlistConf) /* update a bwlist by id. configure, group move all it */
	e.DELETE("/bwlist", deleteBWlistConf) /* delete a bwlist by id */
}
