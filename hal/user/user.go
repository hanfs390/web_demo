package user

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/influxdb_driver"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func createWirelessTblUserInfo(ap string,info *common.WirelessUserInfo) influxdb_driver.InfluxPoints{
	var influx influxdb_driver.InfluxPoints
	tag := make(map[string]string)
	filed := make(map[string]interface{})
	tag["ApMac"] = ap
	tag["Mac"] = info.Mac
	filed["Logintime"] = info.LoginTime
	filed["Logofftime"] = time.Now().Format("2006-01-02 15:04:05")
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, info.LoginTime, loc)
	timestamp := tmp.Unix()
	filed["Uptime"] = time.Now().Unix() - timestamp
	filed["Ip"] = info.Ip
	filed["Wlan"] = info.Wlan
	filed["Radio"] = info.Radio
	filed["Channel"] = info.Channel
	filed["DevModel"] = info.DevModel
	filed["DevType"] = info.DevType
	filed["OSType"] = info.OSType
	filed["Signal"] = info.Signal
	filed["TxRate"] = info.TxRate
	filed["TxBytes"] = info.TxBytes
	filed["RxRate"] = info.RxRate
	filed["RxBytes"] = info.RxBytes
	filed["Hostname"] = info.HostName
	influx.TableName = "TblUser"
	influx.Tag = tag
	influx.Field = filed
	return influx
}

func createWiredTblUserInfo(ap string,info *common.WiredUserInfo) influxdb_driver.InfluxPoints{
	var influx influxdb_driver.InfluxPoints
	tag := make(map[string]string)
	filed := make(map[string]interface{})
	tag["ApMac"] = ap
	tag["Mac"] = info.Mac
	filed["Port"] = info.Port
	filed["LoginTime"] = info.LoginTime
	filed["Logofftime"] = time.Now().Format("2006-01-02 15:04:05")
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, info.LoginTime, loc)
	timestamp := tmp.Unix()
	filed["Uptime"] = time.Now().Unix() - timestamp
	influx.TableName = "TblUser"
	influx.Tag = tag
	influx.Field = filed
	return influx
}

func checkUserSaveToTblUser(info string, last *common.UserReport) error{
	var msg common.UserReport
	var point []influxdb_driver.InfluxPoints
	oldWireless:= make(map[string]common.WirelessUserInfo)
	oldWired := make(map[string]common.WiredUserInfo)
	newWireless := make(map[string]common.WirelessUserInfo)
	newWired := make(map[string]common.WiredUserInfo)
	err := json.Unmarshal([]byte(info),&msg)
	if err != nil{
		return err
	}
	for _, v := range msg.WirelessUser{
		oldWireless[v.Mac] = v
	}
	for _, v := range msg.WiredUser{
		oldWired[v.Mac] = v
	}
	for _,v := range last.WirelessUser{
		newWireless[v.Mac] = v
	}
	for _,v := range last.WiredUser{
		newWired[v.Mac] = v
	}
	for k, v := range oldWireless{
		if _,ok := newWireless[k];ok == false{
			influx := createWirelessTblUserInfo(last.APMac,&v)
			point = append(point,influx)
		}
	}
	for k, v := range oldWired{
		if _, ok := newWired[k]; ok == false{
			influx := createWiredTblUserInfo(last.APMac,&v)
			point = append(point,influx)
		}
	}
	if len(point) > 0 {
		err = influxdb_driver.Service.Insert(point)
		if err != nil {
			return err
		}
	}
	return  nil
}

func updateUserReportMsg(msg *common.UserReport){
	v, err := redis_driver.RedisDb.IsKeyExit(msg.APMac)
	if err != nil || v == 0{
		aclog.Warning("AP [%s] unregistered",msg.APMac)
		return
	}
	key := "user_" + msg.APMac
	user, err := redis_driver.RedisDb.Get(key)
	if err != nil{
		data, err := json.Marshal(msg)
		if err != nil{
			aclog.Error("Change user report info to json error:%s",err.Error())
			return
		}
		 err = redis_driver.RedisDb.Set(key, string(data))
		 if err != nil{
			 aclog.Error("Save Report user Info error:%s",err.Error())
			 return
		 }
	}else {
			err := checkUserSaveToTblUser(user,msg)
			if err != nil{
				aclog.Error("Update TblUser error:%s",err.Error())
				return
			}
	}
	data, err := json.Marshal(msg)
	if err != nil{
		aclog.Error("Change report user info to json error:%s",err.Error())
		return
	}
	err = redis_driver.RedisDb.Set(key,string(data))
	if err != nil{
		aclog.Error("Update user info to redis error:%s",err.Error())
		return
	}
}

func HandleUserReport(client mqtt.Client, message mqtt.Message){
	var user common.UserReport
	err := json.Unmarshal([]byte(message.Payload()),&user)
	if err != nil{
		aclog.Error("Parse user report message error")
		return
	}
	go updateUserReportMsg(&user)
}

func UserReportStart() error{
	err := mqtt_driver.MqttHandle.SubscribeTopic(common.UserReportTopic, 0, HandleUserReport)
	if err != nil{
		aclog.Error("Subcribe topic %s error will exit process",common.UserReportTopic)
		panic(err)
	}
	return nil
}
