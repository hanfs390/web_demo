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
	"errors"
)

func delAPsLoadBalanceConf(aps []string) error {
	value := make(map[string]interface{})
	value["lb_name"] = ""

	g := make(map[string]interface{})
	g["LBName"] = ""

	for i := 0; i < len(aps); i++ {
		/* mysql */
		err := api.HandleUpdateApByMac(aps[i], value)
		if err != nil {
			fmt.Println(err)
		}
		/* redis */
		err = redis_driver.RedisDb.BatchHashSet(aps[i], g)
		if err != nil {
			fmt.Println(err)
			t := task.Task{Time:time.Now().Unix(), Op:"HashSet", Key:aps[i], Value:g}
			task.TaskRedis = append(task.TaskRedis, t)
		}
	}
	return nil
}
func addAPsLoadBalanceConf(aps []string, lbName string) error {
	value := make(map[string]interface{})
	value["lb_name"] = lbName

	g := make(map[string]interface{})
	g["LBName"] = lbName

	for i := 0; i < len(aps); i++ {
		/* mysql */
		err := api.HandleUpdateApByMac(aps[i], value)
		if err != nil {
			fmt.Println(err)
		}
		/* redis */
		err = redis_driver.RedisDb.BatchHashSet(aps[i], g)
		if err != nil {
			fmt.Println(err)
			t := task.Task{Time:time.Now().Unix(), Op:"HashSet", Key:aps[i], Value:g}
			task.TaskRedis = append(task.TaskRedis, t)
		}
	}
	return nil
}
func isExitLbInAp(mac string) error {
	ap, err := api.HandleFindApByMac(mac)
	if err != nil {
		return err
	}
	if ap.LBName != "" {
		return errors.New("AP "+ mac + " load balance group name not NULL")
	}
	return nil
}
func getAllLbConfs(c echo.Context) error {
	result, err := api.HandleFindLbsAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
func getLbConfById(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "Load Balance ID is error")
	}
	result, err := api.HandleFindLbById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func createLbConf(c echo.Context) error {
	lb := common.TblLoadBalanceGroup{}
	if err := c.Bind(&lb); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if lb.GroupId <= 1 { /* default group has not led */
		return c.JSON(http.StatusBadRequest, "Group ID is error")
	}
	if lb.Name == "" {
		return c.JSON(http.StatusBadRequest, "load balance group name can not NULL")
	}
	/* get the group info by group id */
	r, err := api.HandleFindGroupById(lb.GroupId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lb.GroupId = r.Id
	lb.GroupName = r.Name
	lb.GroupUrl = r.Url

	var aps []string
	err = json.Unmarshal([]byte(lb.Member), &aps)
	if err != nil {
		fmt.Println(err)
	}
	/* check the AP has lb or not */
	for i := 0; i < len(aps); i++ {
		err := isExitLbInAp(aps[i])
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}
	/* insert */
	err = api.HandleCreateLb(lb)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for i := 0 ; i <len(aps); i++ {
		fmt.Println(aps[i])
	}
	addAPsLoadBalanceConf(aps, lb.Name)

	err = redis_driver.RedisDb.Set("lb_"+lb.Name, lb.Member)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"lb_"+lb.Name, Value:lb.Member}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("lb_"+lb.Name+"_MD5", getMD5(lb.Member))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"lb_"+lb.Name+"_MD5", Value:getMD5(lb.Member)}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return c.JSON(http.StatusOK, "OK")
}

func updateLbConf(c echo.Context) error {
	lb := common.TblLoadBalanceGroup{}
	value := make(map[string]interface{})
	if err := c.Bind(&lb); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if lb.Id <= 0 {
		return c.JSON(http.StatusBadRequest, "Load Balance ID is error")
	}
	/* verify elements */
	if lb.GroupId > 1 {
		r, err := api.HandleFindGroupById(lb.GroupId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		value["group_id"] = r.Id
		value["group_name"] = r.Name
		value["group_url"] = r.Url
	}

	if len(lb.Name) <= 255 {
		value["name"] = lb.Name
	}
	if len(lb.Member) <= 255 {
		value["member"] = lb.Member
	}
	old, err := api.HandleFindLbById(lb.Id)
	if err != nil {
		fmt.Println(err)
	}
	var oldAps []string
	err = json.Unmarshal([]byte(old.Member), &oldAps)
	if err != nil {
		fmt.Println(err)
	}
	var aps []string
	err = json.Unmarshal([]byte(lb.Member), &aps)
	if err != nil {
		fmt.Println(err)
	}
	/* update to mySQL */
	err = api.HandleUpdateLb(lb.Id, value)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	delAPsLoadBalanceConf(oldAps)
	addAPsLoadBalanceConf(aps, lb.Name)

	err = redis_driver.RedisDb.Set("lb_"+lb.Name, lb.Member)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Set", Key:"lb_"+lb.Name, Value:lb.Member}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.Set("lb_"+lb.Name+"_MD5", getMD5(lb.Member))
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"lb_"+lb.Name+"_MD5", Value:getMD5(lb.Member)}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	return c.JSON(http.StatusOK, "OK")
}
func deleteLbConf(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "Load Balance ID is error")
	}
	lb, err := api.HandleFindLbById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var oldAps []string
	err = json.Unmarshal([]byte(lb.Member), &oldAps)
	if err != nil {
		fmt.Println(err)
	}

	err = api.HandleDeleteLb(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	delAPsLoadBalanceConf(oldAps)

	/* del redis */
	err = redis_driver.RedisDb.DeleteKey("lb_"+lb.Name)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"lb_"+lb.Name}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	err = redis_driver.RedisDb.DeleteKey("lb_"+lb.Name+"_MD5")
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:"lb_"+lb.Name+"_MD5"}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return c.JSON(http.StatusOK, "OK")
}

func routeLoadBalance(e *echo.Echo) {
	e.GET("/lb", getAllLbConfs)
	e.GET("/lb/id", getLbConfById)
	e.POST("/lb", createLbConf)
	e.PUT("/lb", updateLbConf)
	e.DELETE("/lb", deleteLbConf)
}