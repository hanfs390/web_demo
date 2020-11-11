package gpon

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func updateGponInfo(info *common.GponStatus) {
	v, err := redis_driver.RedisDb.IsKeyExit(info.Mac)
	if err != nil || v == 0{
		aclog.Warning("AP [%s] unregistered",info.Mac)
		return
	}

	data, err := json.Marshal(&info)
	if err != nil{
		aclog.Error("Gpon info to json error:%s",err.Error())
		return
	}
	err = redis_driver.RedisDb.BatchHashSet(info.Mac,"Gpon",string(data))
	if err != nil{
		aclog.Error("Update gpon info error:%s",err.Error())
		return
	}
}

func HandleGponMsg(client mqtt.Client, message mqtt.Message){
	var gpon common.GponStatus
	err := json.Unmarshal([]byte(message.Payload()),&gpon)
	if err != nil{
		aclog.Error("Parse gpon message error")
		return
	}
	go updateGponInfo(&gpon)
}

func GponReportStart() error {
	err := mqtt_driver.MqttHandle.SubscribeTopic(common.GponStateTopic, 0, HandleGponMsg)
	if err != nil{
		aclog.Error("Subcribe topic %s error will exit process",common.GponStateTopic)
		panic(err)
	}
	return  nil
}
