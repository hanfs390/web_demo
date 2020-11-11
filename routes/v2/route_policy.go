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

func getAllPolicyConfs(c echo.Context) error {
	result, err := api.HandleFindPolicysAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
func getPolicyConfById(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "Policy ID is error")
	}
	result, err := api.HandleFindPolicyById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func createPolicyConf(c echo.Context) error {
	policy := common.TblPolicy{}
	if err := c.Bind(&policy); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if policy.GroupId <= 1 { /* default group has not policy */
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}

	/* get the group info by group id */
	r, err := api.HandleFindGroupById(policy.GroupId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	policy.GroupId = r.Id
	policy.GroupName = r.Name
	policy.GroupUrl = r.Url

	/* checkout the policy exist or not in this group */
	temp, err := api.HandleFindPolicyByGroupUrl(policy.GroupUrl)
	if len(temp) > 0 {
		return c.JSON(http.StatusBadRequest, "Policy has exist in this group"+policy.GroupName)
	}
	/* insert */
	err = api.HandleCreatePolicy(policy)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	value, err := json.Marshal(&policy)
	if err != nil {
		fmt.Println(err)
	}
	err = redis_driver.RedisDb.Set("policy_"+policy.GroupUrl, string(value))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"policy_"+policy.GroupUrl, Value:string(value)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("policy_"+policy.GroupUrl+"_MD5", getMD5(string(value)))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"policy_"+policy.GroupUrl+"_MD5", Value:getMD5(string(value))}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}

func updatePolicyConf(c echo.Context) error {
	policy := common.TblPolicy{}
	value := make(map[string]interface{})
	if err := c.Bind(&policy); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if policy.Id <= 0 {
		return c.JSON(http.StatusBadRequest, "Policy ID is error")
	}
	/* verify elements */
	if policy.GroupId > 1 {
		r, err := api.HandleFindGroupById(policy.GroupId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		value["group_id"] = r.Id
		value["group_name"] = r.Name
		value["group_url"] = r.Url
	}
	/* checkout the policy exist or not in this group */
	temp, err := api.HandleFindPolicyByGroupUrl(value["group_url"].(string))
	if len(temp) > 0 {
		return c.JSON(http.StatusBadRequest, "Policy has exist in this group"+value["group_name"].(string))
	}
	if policy.Dual_max_user == 0 {
		value["dual_max_user"] = 0
	} else {
		value["dual_max_user"] = 1
	}
	if policy.Single_max_user == 0 {
		value["single_max_user"] = 0
	} else {
		value["single_max_user"] = 1
	}
	if policy.Rssi_threshold == 0 {
		value["rssi_threshold"] = 0
	} else {
		value["rssi_threshold"] = 1
	}
	if policy.Access_policy == 0 {
		value["access_policy"] = 0
	} else {
		value["access_policy"] = 1
	}
	if policy.Reject_max == 0 {
		value["reject_max"] = 0
	} else {
		value["reject_max"] = 1
	}
	if policy.Rssi_max == 0 {
		value["rssi_max"] = 0
	} else {
		value["rssi_max"] = 1
	}
	if policy.L2_isolation == 0 {
		value["l2_isolation"] = 0
	} else {
		value["l2_isolation"] = 1
	}
	if policy.Band_steering == 0 {
		value["band_steering"] = 0
	} else {
		value["band_steering"] = 1
	}
	if policy.Thredhold_5g == 0 {
		value["thredhold_5g"] = 0
	} else {
		value["thredhold_5g"] = 1
	}
	if policy.Thredhold_5g_rssi == 0 {
		value["thredhold_5g_rssi"] = 0
	} else {
		value["thredhold_5g_rssi"] = 1
	}
	if policy.Roaming_detect == 0 {
		value["roaming_detect"] = 0
	} else {
		value["roaming_detect"] = 1
	}
	if policy.Roaming_assoc_rssi == 0 {
		value["roaming_assoc_rssi"] = 0
	} else {
		value["roaming_assoc_rssi"] = 1
	}
	/* update to mySQL */
	err = api.HandleUpdatePolicy(policy.Id, value)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	data, err := json.Marshal(&policy)
	if err != nil {
		fmt.Println(err)
	}
	err = redis_driver.RedisDb.Set("policy_"+policy.GroupUrl, string(data))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"policy_"+policy.GroupUrl, Value:string(data)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("policy_"+policy.GroupUrl+"_MD5", getMD5(string(data)))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"policy_"+policy.GroupUrl+"_MD5", Value:getMD5(string(data))}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}
func deletePolicyConf(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "Policy ID is error")
	}
	policy, err := api.HandleFindPolicyById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = api.HandleDeletePolicy(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	/* del redis */
	err = redis_driver.RedisDb.DeleteKey("policy_"+policy.GroupUrl)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"policy_"+policy.GroupUrl}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("policy_"+policy.GroupUrl+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"policy_"+policy.GroupUrl+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return c.JSON(http.StatusOK, "OK")
}

func routePolicy(e *echo.Echo) {
	e.GET("/policy", getAllPolicyConfs) /* get the all policy info */
	e.GET("/policy/id", getPolicyConfById) /* get the policy by id */
	e.POST("/policy", createPolicyConf) /* create a new policy, reserve */
	e.PUT("/policy", updatePolicyConf) /* update a policy by id. configure, group move all it */
	e.DELETE("/policy", deletePolicyConf) /* delete a policy by id */
}