package reg

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/influxdb_driver"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/mysql_driver"
	"ByzoroAC/models/redis_driver"
	md52 "crypto/md5"
	"encoding/hex"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var APIsBlocked int = 1

func updateDevTabInfo(msg *common.RegInfo)map[string]interface{}{
	tab := make(map[string]interface{})
	tab["sn"] = msg.Sn
	tab["model"] = msg.Model
	tab["hw_ver"] = msg.HwVer
	tab["hw_unique"] = msg.HwUnique
	tab["dev_name"] = msg.DevName
	tab["firmware_ver"] = msg.Version
	return tab
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

func createDevCfg(msg *common.RegInfo) string{
	switch msg.DevName {
	case "WitMAX-AP530-G" :
		str, err := initGp530Cfg()
		if err != nil{
			aclog.Error("Create gp530 cfg error")
			return ""
		}
		return str
	case "WitMAX-AP630-G":
		str, err := initGp630Cfg()
		if err != nil{
			aclog.Error("Create gp630 cfg error")
			return ""
		}
		return  str
	case "WitMAX-AP830-G":
		str, err := initGp830Cfg()
		if err != nil{
			aclog.Error("Create gp830 cfg error")
			return  ""
		}
		return str
	default:
		aclog.Error("Unrecognized AP type")
		return ""
	}
}

func createDevTabInfo(msg *common.RegInfo, new *common.TblDevInfo){
	new.Mac = msg.Mac
	new.Sn = msg.Sn
	new.Model = msg.Model
	new.HwUnique = msg.HwUnique
	new.DevName = msg.DevName
	new.GroupId = 1
	new.GroupUrl ="/root"
	new.GroupName = "default"
	new.DevCfg = createDevCfg(msg)
	new.FirstLoginTime = time.Now().Format("2006-01-02 15:04:05")
	new.SignUpTime = time.Now().Format("2006-01-02 15:04:05")
	new.FirmwareVer = msg.Version
}

func getCfgMd5(s string) string{
	md5 := md52.New()
	md5.Write([]byte(s))
	str := hex.EncodeToString(md5.Sum(nil))
	return  str
}

func copyDevInfoToTablAllCfg(info *common.TblDevInfo ,cfg*common.TblAllCfg){
	cfg.GroupName = info.GroupName
	cfg.DevCfg = info.DevCfg
	cfg.GroupUrl = info.GroupUrl
	cfg.Mac = info.Mac
	cfg.DevScript = info.DevScript
	cfg.TargetFirmwareVer = info.TargetFirmwareVer
	cfg.GroupId = info.GroupId
	cfg.DevCfgMD5 = getCfgMd5(info.DevCfg)
}

func updateSqlAndGetCfgInfo(msg *common.RegInfo,cfg*common.TblAllCfg )(isBlocked int, err error){
	var serach common.TblDevInfo
	serach.Mac = msg.Mac
	r, err := mysql_driver.Handle.FindFirst(serach)
	if (err == nil) && (r.(common.TblDevInfo) != common.TblDevInfo {}){
		result := r.(common.TblDevInfo)
		if result.Blocked == APIsBlocked{
			aclog.Info("The AP: %s is Blocked",result.Mac)
			return APIsBlocked,nil
		}
		copyDevInfoToTablAllCfg(&result,cfg)
		tab := updateDevTabInfo(msg)
		if result.FirstLoginTime == ""{
			tab["FirstLoginTime"] = time.Now().Format("2006-01-02 15:04:05")
		}
		err := mysql_driver.Handle.Update(serach,tab)
		if err != nil{
			aclog.Error("Update TblDevinfo [%s] error:%s.",serach.Mac,err.Error())
			return  0, err
		}
	}else{
		createDevTabInfo(msg, &serach)
		err := mysql_driver.Handle.Insert(serach)
		copyDevInfoToTablAllCfg(&serach,cfg)
		if err != nil{
			aclog.Error("Insert [%s] error: %s.",serach.Mac,err.Error())
			return  0, err
		}
	}
	return 0, nil
}

func updateTblAllCfgDB(info *common.TblAllCfg) error{
	val := make(map[string]interface{})
	val["DevCfg"] = info.DevCfg
	val["DevCfgMD5"] = info.DevCfgMD5
	val["DevScript"] = info.DevScript
	val["TargetFirmwareVer"] = info.TargetFirmwareVer
	val["GroupName"] = info.GroupName
	val["GroupUrl"] = info.GroupUrl
	val["GroupId"] = info.GroupId
	err := redis_driver.RedisDb.BatchHashSet(info.Mac, val)
	if err != nil{
		return  err
	}
	return nil
}

func updateTblOnlineDev(key string) error{
	val := make(map[string]interface{})
	val["timestamp"] = time.Now().Unix()
	hkey := "online_" + key
	err := redis_driver.RedisDb.BatchHashSet(hkey,val)
	if err != nil{
		return  err
	}
	err = redis_driver.RedisDb.ExpireHashKey(hkey,common.KeyExpireTime)
	if err != nil{
		return err
	}
	return nil
}


func regSuccessSendToAP(mac string)error{
	var res common.ResponseInfo
	var ap []string
	res.Version = common.ResponseVersion
	res.Method = common.RegisterMethod
	res.Contents = "Success"
	data, err := json.Marshal(&res)
	if err != nil{
		return err
	}
	ap = append(ap, mac)
	mqtt_driver.MqttHandle.SendMsgToAP(ap, string(data))
	return nil
}

func updateRegisterState(mac string) error{
	var search common.TblDevStat
	search.Mac = mac
	r, err := mysql_driver.Handle.FindFirst(search)
	if (err == nil) && (r.(common.TblDevStat) != common.TblDevStat {}){
		tbl := make(map[string]interface{})
		tbl["State"] = "online"
		tbl["LastLogInTime"] = time.Now().Format("2006-01-02 15:04:05")
		err := mysql_driver.Handle.Update(search,tbl)
		if err != nil{
			aclog.Error("Update AP [%s] register state error:%s",search.Mac,err.Error())
			return err
		}
	}else{
		search.State = "online"
		search.LastLogInTime = time.Now().Format("2006-01-02 15:04:05")
		err := mysql_driver.Handle.Insert(search)
		if err != nil{
			aclog.Error("Insert TblDevState[%s]error:%s",search.Mac,err.Error())
			return err
		}
	}
	return nil
}

func updateTblOnOffLog(mac string) error {
	var point []influxdb_driver.InfluxPoints
	var influx influxdb_driver.InfluxPoints
	tag := make(map[string]string)
	filed := make(map[string]interface{})
	influx.TableName = "TblOnOffLog"
	tag["Mac"] = mac
	influx.Tag = tag
	filed["State"] = "online"
	filed["LoginTime"] = time.Now().Format("2006-01-02 15:04:05")
	influx.Field = filed
	point = append(point, influx)
	err := influxdb_driver.Service.Insert(point)
	if err != nil {
		return err
	}
	return nil
}

func saveRegMsgToDB(msg *common.RegInfo){
	var allCfg common.TblAllCfg
	 blocked, err := updateSqlAndGetCfgInfo(msg,&allCfg)
	if err != nil || blocked == APIsBlocked{
		aclog.Error("Update sql error or AP is blocked.")
		return
	}
	err = updateTblAllCfgDB(&allCfg)
	if err != nil{
		aclog.Error("Update TblAllCfg error :%s.",err.Error())
		return
	}
	err  = updateTblOnlineDev(msg.Mac)
	if err != nil{
		aclog.Error("Update TblOnlineDev error:%s",err.Error())
		return
	}
	err = updateRegisterState(msg.Mac)
	if err != nil{
		aclog.Error("Update register state error:%s",err.Error())
	}
	err = updateTblOnOffLog(msg.Mac)
	if err != nil{
		aclog.Error("Update TblOnOffLog error:%s",err.Error())
	}
	err = regSuccessSendToAP(msg.Mac)
	if err != nil{
		aclog.Error("Register success send to ap error:%s.",err.Error())
	}
}

func HandleRegiserMsg(client mqtt.Client, message mqtt.Message){
	var reginfo common.RegInfo
	err := json.Unmarshal([]byte(message.Payload()),&reginfo)
	if err != nil{
		aclog.Error("Parse ap register message error")
		return
	}
	go saveRegMsgToDB(&reginfo)
}

func RegisterProcessStart() error {
	 err  := mqtt_driver.MqttHandle.SubscribeTopic(common.RegTopic,0,HandleRegiserMsg)
	 if err != nil{
	 	aclog.Error("Subcribe topic %s error will exit process",common.RegTopic)
	 	panic(err)
	 }
	 return  nil
}