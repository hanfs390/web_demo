package config

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)


func getGroupUrlAndDevCfg(mac string)(url,cfg string){
	v,err := redis_driver.RedisDb.BatchHashGet(mac,"GroupUrl","DevCfg")
	if err != nil || v[0] == nil{
		aclog.Error("Get %s cfg from TblAllCfg error :%s",mac,err.Error())
		return "",""
	}
	url = string(v[0].([]uint8))
	cfg = string(v[1].([]uint8))
	return url,cfg
}

func getCfgFromRedis(key string)string{
	data,err := redis_driver.RedisDb.Get(key)
	if err != nil{
		aclog.Error("Get led config error:%s",err.Error())
		return ""
	}
	return data
}

func updateAllCfg(url,devCfg string)string{
	var update common.UpdateCfg
	update.DevCfg = devCfg
	key := "wlan_" + url
	data, err := redis_driver.RedisDb.Get(key)
	if err == nil{
		update.Wlan = data
	}

	key = "policy_" + url
	data, err = redis_driver.RedisDb.Get(key)
	if err == nil{
		update.Policy = data
	}

	key = "led_" + url
	data, err = redis_driver.RedisDb.Get(key)
	if err == nil{
		update.LED = data
	}

	key = "lb_" + url
	data, err = redis_driver.RedisDb.Get(key)
	if err == nil{
		update.LB = data
	}

	key = "bw_" + url
	data, err = redis_driver.RedisDb.Get(key)
	if err == nil{
		update.BW = data
	}
	 update.DevCfg = devCfg
	 cfg, err := json.Marshal(&update)
	 if err != nil{
	 	aclog.Error("Change config to json error",err.Error())
	 	return ""
	 }
	return string(cfg)
}

func CfgModifySendToAP(cfg string, ap[]string)error{
	if len := len(ap); len == 0{
		aclog.Error("AP list is empty")
		return errors.New("Invalid parameter")
	}
	if _,ok := common.APMethod[cfg];ok == false{
		aclog.Error("Config type error")
		return errors.New("Type error")
	}
	for _,mac := range ap {
		var conf common.ResponseInfo
		var devCfg,url string
		var send []string
		conf.Version = common.ResponseVersion
		conf.Method = common.APMethod[cfg]
		url,devCfg = getGroupUrlAndDevCfg(string(mac))
		if url == "" || devCfg == ""{
			aclog.Error("Get group url or dev config error")
			continue
		}
		if cfg == "all"{
			conf.Contents = updateAllCfg(url,devCfg)
		}else if cfg == "dev"{
			conf.Contents = devCfg
		}else {
			key := cfg + "_" + url
			conf.Contents = getCfgFromRedis(key)
			}
		str, err := json.Marshal(&conf)
		if err != nil {
			aclog.Error("Json marshal error:%s", err.Error())
			return err
		}
		send = append(send,string(mac))
		mqtt_driver.MqttHandle.SendMsgToAP(send, string(str))
	}
	return  nil
}

func checkAndGetCfg(cfg, url, md5 string) (value string, err error){
	var res string
	md5key := cfg + "_" + url + "_MD5"
	v, err := redis_driver.RedisDb.Get(md5key)
	if err != nil{
		return "",err
	}
	if v != md5 {
		cfgKey := cfg + "_" + url
		tmp, err := redis_driver.RedisDb.Get(cfgKey)
		if err != nil{
			return "",err
		}
		res = tmp
	}
	return res, nil
}

func sendUpdateCfgToAp(mac string, cfg *common.UpdateCfg)error{
	var res common.ResponseInfo
	var ap [] string
	res.Version = common.ResponseVersion
	res.Method = common.CfgPollingMethod
	v, err := json.Marshal(cfg)
	if err != nil{
		return err
	}
	res.Contents = string(v)
	ap = append(ap, mac)
	data, err  := json.Marshal(res)
	if err != nil{
		return  err
	}
	mqtt_driver.MqttHandle.SendMsgToAP(ap, string(data))
	return  nil
}

func checkAllCfg(cfg *common.CfgPolling){
	var update common.UpdateCfg
	var url string
	v,err := redis_driver.RedisDb.BatchHashGet(cfg.Mac,"GroupUrl","DevCfgMD5","DevCfg")
	if err != nil || v[0] == nil{
		aclog.Error("Get %s cfg from TblAllCfg error :%s",cfg.Mac,err.Error())
		return
	}
	if cfg.DevCfgMd5 != string(v[1].([]uint8)){
		update.DevCfg = string(v[2].([]uint8))
	}
	url = string(v[0].([]uint8))
	conf ,_:= checkAndGetCfg("wlan",url,cfg.WlanMd5)
	update.Wlan = conf

	conf, _ = checkAndGetCfg("policy",url,cfg.PolicyMd5)
	update.Policy = conf

	conf, _ = checkAndGetCfg("led", url, cfg.LEDMd5)
	update.LED = conf

	conf, _ = checkAndGetCfg("bw", url, cfg.BWMd5)
	update.BW = conf

	conf, _ = checkAndGetCfg("lb", url, cfg.LBMd5)
	update.LB = conf

	err = sendUpdateCfgToAp(cfg.Mac,&update)
	if err != nil{
		aclog.Error("send update config to AP[%s]error",cfg.Mac)
	}
}

func HandleCfgPolling(client mqtt.Client, message mqtt.Message){
	var cfg common.CfgPolling
	err := json.Unmarshal([]byte(message.Payload()),&cfg)
	if err != nil{
		aclog.Error("Parse cfg polling msg error:%s",err.Error())
		return
	}
	go checkAllCfg(&cfg)
}

func initAPMethodSet(){
	common.APMethod = make(map[string]string)
	common.APMethod["wlan"] = common.WlanCfgMethod
	common.APMethod["policy"] = common.PolicyCfgMethod
	common.APMethod["dev"] = common.DevCfgMethod
	common.APMethod["led"] = common.LEDCfgMethod
	common.APMethod["bw"] = common.BWCfgMethod
	common.APMethod["lb"] = common.LBCfgMethod
	common.APMethod["all"] = common.UpdateAllMethod
}

func CfgCheckProcessStart() error{
	err := mqtt_driver.MqttHandle.SubscribeTopic(common.CfgCheckTopic, 0, HandleCfgPolling)
	if err != nil{
		aclog.Error("Subcribe topic :%s error will exit process.",common.CfgCheckTopic)
		panic(err)
	}
	return  nil
}
