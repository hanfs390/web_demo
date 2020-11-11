package heartbeat

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/influxdb_driver"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/mysql_driver"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gomodule/redigo/redis"
	"time"
)

func saveTblOnOffLog(mac string) error {
	var point []influxdb_driver.InfluxPoints
	var influx influxdb_driver.InfluxPoints
	tag := make(map[string]string)
	filed := make(map[string]interface{})
	influx.TableName = "TblOnOffLog"
	tag["Mac"] = mac
	influx.Tag = tag
	filed["State"] = "offline"
	filed["Reason"] = "Time out"
	filed["LogoutTime"] = time.Now().Format("2006-01-02 15:04:05")
	influx.Field = filed
	point = append(point, influx)
	err := influxdb_driver.Service.Insert(point)
	if err != nil {
		return err
	}
	return nil
}


func updateTblDevStatDB(mac string)error{
	var search common.TblDevStat
	search.Mac = mac
	r, err := mysql_driver.Handle.FindFirst(search)
	if (err == nil) && (r.(common.TblDevStat) != common.TblDevStat {}){
		update := make(map[string]interface{})
		update["State"] = "offline"
		update["LastLogOutTime"] = time.Now().Format("2006-01-02 15:04:05")
		err := mysql_driver.Handle.Update(search,update)
		if err != nil{
			return  err
		}
	}else{
		aclog.Error("Can not find [%s] in TblDevStat")
	}
	return nil
}

//event key is online_
func phraseExpireMsgGetMac(data []byte) string{
	len := len(data)
	if len > 7 {
		return string(data[7:])
	}else {
		return ""
	}
}

func handleExipreEvent(data []byte) error {
	key := phraseExpireMsgGetMac(data)
	err := saveTblOnOffLog(key)
	if err != nil{
		aclog.Error("Save offline msg to TblOnOffLog error:%s",err.Error())
	}

	err = updateTblDevStatDB(key)
	if err != nil{
		aclog.Error("Update TblDevStat error:%s",err.Error())
		return err
	}
	return nil
}

func startListenExpireEvent(){
	proto,addr := redis_driver.RedisDb.GetProtoAndAddr()
	c, err := redis.Dial(proto,addr)
	if err != nil{
		aclog.Error("Heart beat process connect server error:%s",err.Error())
		panic(err)
	}
	defer c.Close()
	psc := redis.PubSubConn{c}
	psc.PSubscribe("__keyevent@0__:expired")
	for{
		switch v := psc.Receive().(type){
		case redis.Message:
			go handleExipreEvent(v.Data)
		}
	}
}

func checkAndUpdateOnlineDB(msg *common.HeartBeat){
	v, err := redis_driver.RedisDb.IsKeyExit(msg.Mac)
	if err != nil || v == 0{
		aclog.Warning("AP [%s] unregistered",msg.Mac)
		return
	}
	val := make(map[string]interface{})
	val["timestamp"] = time.Now().Unix()
	hkey := "online_" + msg.Mac
	err = redis_driver.RedisDb.BatchHashSet(hkey,val)
	if err != nil{
		return
	}
	err = redis_driver.RedisDb.ExpireHashKey(hkey,common.KeyExpireTime)
	if err != nil {
		return
	}
}

func HandleHeartBeatMsg(client mqtt.Client, message mqtt.Message){
	var info common.HeartBeat
	err := json.Unmarshal([]byte(message.Payload()),&info)
	if err != nil{
		aclog.Error("Parse dev state msg error :%s",err.Error())
		return
	}
	go checkAndUpdateOnlineDB(&info)
}

func HeartBeatProcessStart() error {
	err := mqtt_driver.MqttHandle.SubscribeTopic(common.HeartBeatTopic,0, HandleHeartBeatMsg)
	if err != nil{
		aclog.Error("Subscribe topic :%s error will exit process",common.HeartBeatTopic)
		panic(err)
	}
	go startListenExpireEvent()
	return  nil
}
