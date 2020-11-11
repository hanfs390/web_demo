package upgrade

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SendUpgradeMsgToAP(ap [] string) error {
	for mac := range ap{
		var res  common.ResponseInfo
		var ap []string
		r, err := redis_driver.RedisDb.BatchHashGet(string(mac),"UpgradeUrl")
		if err == nil{
			upgradeUrl := string(r[0].([]uint8))
			res.Method = common.UpgradeMethod
			res.Version = common.ResponseVersion
			res.Contents = upgradeUrl
			ap = append(ap,string(mac))
			data, err := json.Marshal(&res)
			if err == nil{
				mqtt_driver.MqttHandle.SendMsgToAP(ap,string(data))
			}
		}
	}
	return nil
}

func checkFimwareVersion(info *common.UpgradeInfo){
	var res  common.ResponseInfo
	var ap []string
	v, err := redis_driver.RedisDb.IsKeyExit(info.Mac)
	if err != nil || v == 0{
		aclog.Warning("AP [%s] unregistered",info.Mac)
		return
	}
	r, err := redis_driver.RedisDb.BatchHashGet(info.Mac,"TargetFirmwareVer","UpgradeUrl")
	if err != nil {
		aclog.Error("Get info from TblAllCfg error:%s",err.Error())
		return
	}
	firmwareVer :=  string(r[0].([]uint8))
	upgradeUrl := string(r[1].([]uint8))
	if firmwareVer != info.TargetFirmwareVer{
		res.Method = common.UpgradeMethod
		res.Version = common.ResponseVersion
		res.Contents = upgradeUrl
		ap = append(ap,info.Mac)
		data, err := json.Marshal(&res)
		if err != nil{
			aclog.Error("Change upgrade info to json error :%s",err.Error())
			return
		}
		mqtt_driver.MqttHandle.SendMsgToAP(ap,string(data))
	}
}

func HandleUpgradeProcess(client mqtt.Client, message mqtt.Message){
	var info  common.UpgradeInfo
	err := json.Unmarshal([]byte(message.Payload()),&info)
	if err != nil{
		aclog.Error("Parse upgrade msg error :%s",err.Error())
		return
	}
	go checkFimwareVersion(&info)
}

func UpgradeProcessStart() error {
	err := mqtt_driver.MqttHandle.SubscribeTopic(common.UpgradeTopic,0, HandleUpgradeProcess)
	if err != nil{
		aclog.Error("Subscribe topic :%s error will exit process",common.UpgradeTopic)
		panic(err)
	}
	return  nil
}
