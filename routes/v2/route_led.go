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
)

func getAllLedConfs(c echo.Context) error {
	result, err := api.HandleFindLedsAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
func getLedConfById(c echo.Context) error {
	fmt.Println(" get LED by id")
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "LED ID is error")
	}
	result, err := api.HandleFindLedById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func createLedConf(c echo.Context) error {
	led := common.TblLED{}
	if err := c.Bind(&led); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if led.GroupId <= 1 { /* default group has not led */
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}
	/* get the group info by group id */
	r, err := api.HandleFindGroupById(led.GroupId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	led.GroupId = r.Id
	led.GroupName = r.Name
	led.GroupUrl = r.Url

	/* checkout the led exist or not in this group */
	temp, err := api.HandleFindPolicyByGroupUrl(led.GroupUrl)
	if len(temp) > 0 {
		return c.JSON(http.StatusBadRequest, "Led has exist in this group"+led.GroupName)
	}
	/* insert */
	err = api.HandleCreateLed(led)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = redis_driver.RedisDb.Set("led_"+led.GroupUrl, led.Timer)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"led_"+led.GroupUrl, Value:led.Timer}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("led_"+led.GroupUrl+"_MD5", getMD5(led.Timer))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"led_"+led.GroupUrl+"_MD5", Value:getMD5(led.Timer)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}

func updateLedConf(c echo.Context) error {
	led := common.TblLED{}
	value := make(map[string]interface{})
	if err := c.Bind(&led); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if led.Id <= 0 {
		return c.JSON(http.StatusBadRequest, "LED ID is error")
	}
	/* verify elements */
	if led.GroupId > 1 {
		r, err := api.HandleFindGroupById(led.GroupId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		value["group_id"] = r.Id
		value["group_name"] = r.Name
		value["group_url"] = r.Url
	}
	/* checkout the led exist or not in this group */
	temp, err := api.HandleFindPolicyByGroupUrl(value["group_url"].(string))
	if len(temp) > 0 {
		return c.JSON(http.StatusBadRequest, "Led has exist in this group"+value["group_name"].(string))
	}

	if len(led.Timer) <= 255 {
		value["timer"] = led.Timer
	}

	/* update to mySQL */
	err = api.HandleUpdateLed(led.Id, value)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = redis_driver.RedisDb.Set("led_"+led.GroupUrl, led.Timer)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"led_"+led.GroupUrl, Value:led.Timer}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("led_"+led.GroupUrl+"_MD5", getMD5(led.Timer))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"led_"+led.GroupUrl+"_MD5", Value:getMD5(led.Timer)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}
func deleteLedConf(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "LED ID is error")
	}
	led, err := api.HandleFindLedById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = api.HandleDeleteLed(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* del redis */
	err = redis_driver.RedisDb.DeleteKey("led_"+led.GroupUrl)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"led_"+led.GroupUrl}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("led_"+led.GroupUrl+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"led_"+led.GroupUrl+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return c.JSON(http.StatusOK, "OK")
}

func routeLed(e *echo.Echo) {
	e.GET("/led", getAllLedConfs) /* get the all WLANs info */
	e.GET("/led/id", getLedConfById) /* get the WLAN by id */
	e.POST("/led", createLedConf) /* create a new WLAN, reserve */
	e.PUT("/led", updateLedConf) /* update a WLAN by id. configure, group move all it */
	e.DELETE("/led", deleteLedConf) /* delete a WLAN by id */
}