package v2

import (
	"github.com/labstack/echo"
	"ByzoroAC/controllers/api"
	"fmt"
	"net/http"
	"strconv"
	"ByzoroAC/common"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	"ByzoroAC/controllers/task"
	"time"
	"ByzoroAC/hal/config"
	"errors"
)

func wlanRedisRefresh(groupUrl string) (result string, err error) {
	wlanArr, err := api.HandleFindWlanByGroupUrl(groupUrl)
	if err != nil {
		fmt.Println("json:", err)
		return "", err
	}
	if len(wlanArr) == 0 {
		err := redis_driver.RedisDb.DeleteKey("wlan_"+groupUrl)
		if err != nil {

		}
		err = redis_driver.RedisDb.DeleteKey("wlan_"+groupUrl+"MD5")
		if err != nil {

		}
		return "", err
	}
	value, err := json.Marshal(&wlanArr)
	if err != nil {
		fmt.Println("json:", err)
		return "", err
	}
	err = redis_driver.RedisDb.Set("wlan_"+groupUrl, string(value))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"wlan_"+groupUrl, Value:string(value)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("wlan_"+groupUrl+"_MD5", getMD5(string(value)))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"wlan_"+groupUrl+"_MD5", Value:getMD5(string(value))}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return string(value), nil
}

func getAllWlan(c echo.Context) error {
	result, err := api.HandleFindWlansAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
func getWlanById(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "WLAN ID is error")
	}
	result, err := api.HandleFindWlanById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func createWlan(c echo.Context) error {
	wlan := common.TblWlan{}
	if err := c.Bind(&wlan); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(wlan)
	if wlan.GroupId <= 1 { /* default group has not wlans */
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}
	/* get the group info by group id */
	r, err := api.HandleFindGroupById(wlan.GroupId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	wlan.GroupId = r.Id
	wlan.GroupName = r.Name
	wlan.GroupUrl = r.Url
	/* verify the wlan count */
	count, _ := api.HandleFindWlansCountByGroupId(wlan.GroupId)
	if count >= 4 {
		return c.JSON(http.StatusBadRequest, "WLAN is too many")
	}
	/* check ssid */
	sameSsidWlans, err := api.HandleFindWlanByGroupIdAndSsid(wlan.GroupId, wlan.Ssid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for len(sameSsidWlans) > 0 {
		return c.JSON(http.StatusBadRequest, "Ssid has exist")
	}
	/* insert */
	err = api.HandleCreateWlan(wlan)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, e := wlanRedisRefresh(wlan.GroupUrl)
	if e != nil {
		return c.JSON(http.StatusBadRequest, e.Error())
	}

	aps, err := api.HandleFindApsByUrl(wlan.GroupUrl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	config.CfgModifySendToAP("wlan", aps)
	return c.JSON(http.StatusOK, "OK")
}

func updateWlan(c echo.Context) error {
	wlan := common.TblWlan{}
	value := make(map[string]interface{})
	if err := c.Bind(&wlan); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if wlan.Id <= 0 {
		return c.JSON(http.StatusBadRequest, "WLAN ID is error")
	}
	/* verify elements */
	if wlan.GroupId > 1 {
		r, err := api.HandleFindGroupById(wlan.GroupId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		value["group_id"] = r.Id
		value["group_name"] = r.Name
		value["group_url"] = r.Url
	}

	if len(wlan.Ssid) <= 255 {
		value["ssid"] = wlan.Ssid
	}
	sameSsidWlans, err := api.HandleFindWlanByGroupIdAndSsid(wlan.GroupId, wlan.Ssid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for i := 0; i < len(sameSsidWlans); i++ {
		if sameSsidWlans[i].Id != wlan.Id {
			return c.JSON(http.StatusBadRequest, errors.New("ssid has exist"))
		}
	}
	if (wlan.Disabled >= 0) && (wlan.Disabled <= 1) {
		value["disabled"] = wlan.Disabled
	}
	if len(wlan.Encryption) <= 32 {
		value["encryption"] = wlan.Encryption
	}
	if len(wlan.Hidden) <= 32 {
		value["hidden"] = wlan.Hidden
	}
	if (wlan.VlanSwitch >= 0) && (wlan.VlanSwitch <= 1) {
		value["vlan_switch"] = wlan.VlanSwitch
	}
	if wlan.VlanId >= 0 {
		value["vlan_id"] = wlan.VlanId
	}
	if len(wlan.Radio) <= 32 {
		value["radio"] = wlan.Radio
	}
	if wlan.PerUserRate >= 0 {
		value["per_user_rate"] = wlan.PerUserRate
	}
	if wlan.UpPerUserRate >= 0 {
		value["up_per_user_rate"] = wlan.UpPerUserRate
	}
	if wlan.DownPerUserRate >= 0 {
		value["down_per_user_rate"] = wlan.DownPerUserRate
	}
	if len(wlan.AuthServer) <= 32 {
		value["auth_server"] = wlan.AuthServer
	}
	if len(wlan.AuthPort) <= 32 {
		value["auth_port"] = wlan.AuthPort
	}
	if len(wlan.AuthSecret) <= 255 {
		value["auth_secret"] = wlan.AuthSecret
	}
	if wlan.MultServer >= 0 {
		value["mult_server"] = wlan.MultServer
	}
	if len(wlan.Servercfg) <= 1024 {
		value["servercfg"] = wlan.Servercfg
	}
	if len(wlan.FirstRoamcfg) <= 1024 {
		value["first_roamcfg"] = wlan.FirstRoamcfg
	}
	if len(wlan.Key) <= 64 {
		value["key"] = wlan.Key
	}
	/* update to mySQL */
	err = api.HandleUpdateWlan(wlan.Id, value)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	/* update redis */
	_, e := wlanRedisRefresh(wlan.GroupUrl)
	if e != nil {
		return c.JSON(http.StatusBadRequest, e.Error())
	}
	aps, err := api.HandleFindApsByUrl(wlan.GroupUrl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	config.CfgModifySendToAP("wlan", aps)

	return c.JSON(http.StatusOK, "OK")
}
func deleteWlan(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "WLAN ID is error")
	}
	wlan, err := api.HandleFindWlanById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = api.HandleDeleteWlan(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* update redis */
	_, e := wlanRedisRefresh(wlan.GroupUrl)
	if e != nil {
		return c.JSON(http.StatusBadRequest, e.Error())
	}
	aps, err := api.HandleFindApsByUrl(wlan.GroupUrl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	config.CfgModifySendToAP("wlan", aps)

	return c.JSON(http.StatusOK, "OK")
}

func routeWlan(e *echo.Echo) {
	e.GET("/wlan", getAllWlan) /* get the all WLANs info */
	e.GET("/wlan/id", getWlanById) /* get the WLAN by id */
	e.POST("/wlan", createWlan) /* create a new WLAN, reserve */
	e.PUT("/wlan", updateWlan) /* update a WLAN by id. configure, group move all it */
	e.DELETE("/wlan", deleteWlan) /* delete a WLAN by id */
}