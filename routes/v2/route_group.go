package v2

import (
	"net/http"
	"github.com/labstack/echo"
	"ByzoroAC/controllers/api"
	"fmt"
	"strconv"
	"ByzoroAC/common"
	"ByzoroAC/models/redis_driver"
	"ByzoroAC/controllers/task"
	"time"
	"errors"
)

type webGroupList struct {
	Id uint64
	Name string
	APCount string
	WlanCount int
	Tag string
	Detail string
}

func delGroupCfgRedis(group_url string) error {
	err := redis_driver.RedisDb.DeleteKey("wlan_"+group_url)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"wlan_"+group_url}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("wlan_"+group_url+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"wlan_"+group_url+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("led_"+group_url)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"led_"+group_url}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("led_"+group_url+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"led_"+group_url+"MD5_"}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("policy_"+group_url)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"policy_"+group_url}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("policy_"+group_url+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"policy_"+group_url+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("bwlist_"+group_url)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"bwlist_"+group_url}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("bwlist_"+group_url+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"bwlist_"+group_url+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return err
}
func delGroupCfgMysql(groupUrl string) error {
	/* del group policy */
	policy, err := api.HandleFindPolicyByGroupUrl(groupUrl)
	if err != nil {
		fmt.Println("find error")
	}
	for i := 0; i < len(policy); i++ {
		fmt.Println(policy[i].Id)
		err := api.HandleDeletePolicy(policy[i].Id)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("del policy OK")
	/* del group led */
	leds, err := api.HandleFindLedByGroupUrl(groupUrl)
	if err != nil {
		fmt.Println("find error")
	}
	for i := 0; i < len(leds); i++ {
		err := api.HandleDeleteLed(leds[i].Id)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("del led OK")
	/* del group bw list */
	bws, err := api.HandleFindBWlistByGroupUrl(groupUrl)
	if err != nil {
		fmt.Println("find error")
	}
	for i := 0; i < len(bws); i++ {
		err := api.HandleDeleteBWlist(bws[i].Id)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("del bw OK")
	/* del group wlan */
	wlans, err := api.HandleFindWlanByGroupUrl(groupUrl)
	if err != nil {
		fmt.Println("find error")
	}
	for i := 0; i < len(wlans); i++ {
		err := api.HandleDeleteWlan(wlans[i].Id)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("del wlan OK", wlans[i].Id)
	}
	return nil
}
func moveApsToDefaultGroup(groupUrl string) error {
	value := make(map[string]interface{})

	value["group_id"] = 1;
	value["group_name"] = "default"
	value["group_url"] = "/root"

	api.HandleUpdateApByGroupUrl(groupUrl, value)
	return nil
}
func getAllGroups(c echo.Context) error {
	fmt.Println("get groups")
	result, err := api.HandleFindGroupsAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var info []webGroupList
	for i := 0; i < len(result); i++ {
		if result[i].Id <= 1 {
			continue
		}
		allAPCount, onAPCount, err := api.HandleFindApCountByGroupId(result[i].Id)
		fmt.Println(allAPCount, onAPCount)
		if err != nil {
			allAPCount = -1
			onAPCount = -1
		}

		wlanCount, err := api.HandleFindWlansCountByGroupId(result[i].Id)
		if err != nil {
			wlanCount = -1
		}
		temp := webGroupList{Id:result[i].Id, Name:result[i].Name, Tag:result[i].Tag, Detail:result[i].Detail,
					APCount:fmt.Sprintf("%d/%d", onAPCount, allAPCount), WlanCount:wlanCount}

		info = append(info, temp)
	}
	return c.JSON(http.StatusOK, info)
}
func getGroupById(c echo.Context) error {
	fmt.Println(" get group by id")
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 1 {
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}

	result, err := api.HandleFindGroupById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func createGroup(c echo.Context) error {
	group := common.TblGroup{}
	if err := c.Bind(&group); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(group)
	if group.Name == "" {
		fmt.Println("Group Name is empty")
		return c.JSON(http.StatusBadRequest, errors.New("Group Name is empty"))
	}
	/* generate the uuid */
	uu, e := generateUuid()
	if e != nil {
		fmt.Println("create", e)
		return c.JSON(http.StatusBadRequest, e.Error())
	}
	group.Uuid = uu
	/* generate the url */
	group.Url = "/root/" + group.Uuid
	err := api.HandleCreateGroup(group)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "OK")
}
func updateGroup(c echo.Context) error {
	group := common.TblGroup{}
	value := make(map[string]interface{})
	if err := c.Bind(&group); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if group.Id <= 1 {
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}
	fmt.Println(group)
	if len(group.Detail) <= 255 {
		value["detail"] = group.Detail
	}
	if len(group.Tag) <= 255 {
		value["tag"] = group.Tag
	}
	if len(group.GroupCfg) <= 255 {
		value["group_cfg"] = group.GroupCfg
	}
	err := api.HandleUpdateGroup(group.Id, value)
	if err != nil {
		fmt.Println("create", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "OK")
}
func deleteGroup(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 1 {
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}
	result, err := api.HandleFindGroupById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = api.HandleDeleteGroup(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* delete group cfg */
	delGroupCfgRedis(result.Url)
	delGroupCfgMysql(result.Url)
	/* delete load balance */

	/* move ap to default group */
	moveApsToDefaultGroup(result.Url)

	return c.JSON(http.StatusOK, "OK")
}
func routeGroup(e *echo.Echo) {
	e.GET("/group", getAllGroups) /* get the all groups info */
	e.GET("/group/id", getGroupById) /* get the group by id */
	e.POST("/group", createGroup) /* create a new group */
	e.PUT("/group", updateGroup) /* update a group by id */
	e.DELETE("/group", deleteGroup) /* delete a group by id */
}

