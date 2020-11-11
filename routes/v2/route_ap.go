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
	"ByzoroAC/common/model"
	"errors"
	"encoding/json"
	"ByzoroAC/hal/config"
)
type webApList struct {
	Id uint64
	Mac string
	Ip string
	Model string
	Version string
	Alias string
	Tag string
	Detail string
	State string
	LBName string
	GroupName string
	Block int
}
func initGp530Cfg() (r string,err error){
	var cfg common.ModelGP530
	cfg.Country = "CN"
	cfg.Channel2G = 0
	cfg.Channel5G = 0
	cfg.HtMode2G = "auto"
	cfg.HwMode2G = "11ng"
	cfg.HtMode5G = "auto"
	cfg.HwMode5G = "11ac"
	cfg.TxPower2G = 0
	cfg.TxPower5G = 0
	v, err := json.Marshal(&cfg)
	if err != nil{
		return  "",err
	}
	return string(v),nil
}

func initGp630Cfg()(r string,err error){
	var cfg common.ModelGP630
	cfg.Country = "CN"
	cfg.Channel2G = 0
	cfg.Channel5G1 = 0
	cfg.Channel5G2 = 0
	cfg.HtMode2G = "auto"
	cfg.HwMode2G = "11ng"
	cfg.HtMode5G1 = "auto"
	cfg.HwMode5G1 = "11ac"
	cfg.HtMode5G2 = "auto"
	cfg.HwMode5G2 = "11ac"
	cfg.TxPower2G = 0
	cfg.TxPower5G1 = 0
	cfg.TxPower5G2 = 0
	v, err := json.Marshal(&cfg)
	if err != nil{
		return  "",err
	}
	return string(v),nil
}

func initGp830Cfg()(r string,err error){
	var cfg common.ModelGP830
	cfg.Country = "CN"
	cfg.Channel2G = 0
	cfg.Channel5G1 = 0
	cfg.Channel5G2 = 0
	cfg.HtMode2G = "auto"
	cfg.HwMode2G = "11ng"
	cfg.HtMode5G1 = "auto"
	cfg.HwMode5G1 = "11ac"
	cfg.HtMode5G2 = "auto"
	cfg.HwMode5G2 = "11ac"
	cfg.TxPower2G = 0
	cfg.TxPower5G1 = 0
	cfg.TxPower5G2 = 0
	v, err := json.Marshal(&cfg)
	if err != nil{
		return  "",err
	}
	return string(v),nil
}
func generateDefaultDevCfg(mod string) (result string) {
	if mod == "GP530" {
		result, _ = initGp530Cfg()
	} else if mod == "GP830" {
		result, _ = initGp830Cfg()
	} else if mod == "GP630" {
		result, _ = initGp630Cfg()
	}
	return result
}
func getModelInfo(c echo.Context) error {
	mod := c.QueryParam("model")
	fmt.Println("get the mode", mod)
	info := model.GetBoardInfoByModel(mod)
	if info == nil {
		return c.JSON(http.StatusBadRequest, errors.New("this model is not exist"))
	}

	return c.JSON(http.StatusOK, *info)
}

func getAllAp(c echo.Context) error {
	var list []webApList
	result, err := api.HandleFindApsAll()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for i := 0; i < len(result); i++ {

		temp := webApList{Id:result[i].Id, Model:result[i].Model, Alias:result[i].Alias, Tag:result[i].Tag, Detail:result[i].Detail,
							LBName:result[i].LBName, GroupName:result[i].GroupName, Block:result[i].Blocked}
		temp.Mac = model.TransformMacStrToMacForm(result[i].Mac)
		state, err := api.HandleFindApStateByMac(result[i].Mac)
		if err == nil {
			temp.Version = state.FirmwareVer
			temp.Ip = state.Ip
			temp.State = state.State
		}

		list = append(list, temp)
	}
	return c.JSON(http.StatusOK, list)
}
func getApById(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "AP ID is error")
	}
	result, err := api.HandleFindApById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	result.Mac = model.TransformMacStrToMacForm(result.Mac)
	return c.JSON(http.StatusOK, result)
}

func createAp(c echo.Context) error {
	ap := common.TblDevInfo{}
	/* parse the params */
	if err := c.Bind(&ap); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if !model.CheckModelIsExist(ap.Model) {
		return c.JSON(http.StatusBadRequest, errors.New("Error Model: " + ap.Model))
	}
	macStr, err := model.CheckMacForm(ap.Mac)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ap.Mac = macStr
	/* full the default group info */
	ap.GroupId = 1;
	ap.GroupName = "default"
	ap.GroupUrl = "/root"

	ap.SignUpType = "import"
	ap.DevCfg = generateDefaultDevCfg(ap.Model)
	err = api.HandleCreateAp(ap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "OK")
}

func updateAp(c echo.Context) error {
	ap := common.TblDevInfo{}
	value := make(map[string]interface{})

	if err := c.Bind(&ap); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if ap.Id <= 0 {
		fmt.Println("wrong id")
		return c.JSON(http.StatusBadRequest, "AP ID is error")
	}
	/* verify elements */
	if ap.GroupId > 0 {
		r, err := api.HandleFindGroupById(ap.GroupId)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		value["group_id"] = r.Id
		value["group_name"] = r.Name
		value["group_url"] = r.Url
	}
	if len(ap.Alias) <= 255 {
		value["alias"] = ap.Alias
	}
	if len(ap.Detail) <= 255 {
		value["detail"] = ap.Detail
	}
	if ap.Blocked >= 0 && ap.Blocked <= 1 {
		value["blocked"] = ap.Blocked
	}
	if len(ap.Tag) <= 255 {
		value["tag"] = ap.Tag
	}
	if len(ap.DevCfg) <= 1024 {
		value["dev_cfg"] = ap.DevCfg
	}
	/* update to mySQL */
	err := api.HandleUpdateAp(ap.Id, value)
	if err != nil {
		fmt.Println("create", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	/* update redis */
	g := make(map[string]interface{})

	g["DevCfg"] = ap.DevCfg
	g["DevCfgMD5"] = getMD5(ap.DevCfg)
	g["DevScript"] = ap.DevScript
	g["TargetFirmwareVer"] = ap.TargetFirmwareVer
	g["LBName"] = ap.LBName
	g["GroupId"] = value["group_id"]
	g["GroupName"] = value["group_name"]
	g["GroupUrl"] = value["group_url"]

	mac, err := model.CheckMacForm(ap.Mac)
	if err != nil {
		fmt.Println(err)
	}
	err = redis_driver.RedisDb.BatchHashSet(mac, g)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"HashSet", Key:ap.Mac, Value:g}
		task.TaskRedis = append(task.TaskRedis, t)
	}
	var aps []string
	aps = append(aps, mac)
	config.CfgModifySendToAP("dev", aps)
	return c.JSON(http.StatusOK, "OK")
}
func deleteAp(c echo.Context) error {
	id_str := c.QueryParam("id")
	id, err := strconv.ParseUint(id_str, 10, 64) /* base is 2 8 10 16 */
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, "AP ID is error")
	}
	ap, err := api.HandleFindApById(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	api.HandleDeleteAp(id)

	err = redis_driver.RedisDb.DeleteKey(ap.Mac)
	if err != nil {
		fmt.Println(err)
		t := task.Task{Time:time.Now().Unix(), Op:"Del", Key:ap.Mac}
		task.TaskRedis = append(task.TaskRedis, t)
	}

	return c.JSON(http.StatusOK, "OK")
}

func routeAp(e *echo.Echo) {
	e.GET("/ap", getAllAp) /* get the all Aps info */
	e.GET("/ap/id", getApById) /* get the AP by id */
	e.GET("/ap/model", getModelInfo) /* get the AP by id */
	e.POST("/ap", createAp) /* create a new group, reserve */
	e.PUT("/ap", updateAp) /* update a AP by id. configure, group move all it */
	e.DELETE("/ap", deleteAp) /* delete a AP by id */
}