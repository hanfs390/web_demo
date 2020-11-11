package mqtt_driver

import (
	"ByzoroAC/aclog"
	"ByzoroAC/conf"
	"testing"
)

func TestMain(m *testing.M){
	aclog.Init()
	conf.ReadGlobalConf()
	InitMqtt()
	var mac []string
	var data string ="Mqtt Test"
	mac = append(mac, "112233445566")
	MqttHandle.SendMsgToAP(mac,data)
}


