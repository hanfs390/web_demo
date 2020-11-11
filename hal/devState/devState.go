package devState

import (
	"ByzoroAC/aclog"
	"ByzoroAC/common"
	"ByzoroAC/models/influxdb_driver"
	"ByzoroAC/models/mqtt_driver"
	"ByzoroAC/models/mysql_driver"
	"ByzoroAC/models/redis_driver"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"reflect"
)

func updateDevState(info *common.DevStatus)map[string]interface{}{
	dev := make(map[string]interface{})
	elem := reflect.ValueOf(info).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField();i++{
		dev[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return  dev
}

func copyInfoToDevStateTbl(info *common.DevStatus, tbl *common.TblDevStat){
	tbl.Mac = info.Mac
	tbl.DevBootTime = info.DevBootTime
	tbl.FirmwareVer = info.FirmwareVer
	tbl.Ip = info.Ip
	tbl.UserConnected = info.UserConnected
	tbl.WirelessUserConnectedCount = info.WirelessUserConnectedCount
	tbl.WiredUserConnectedCount = info.WiredUserConnectedCount
	tbl.UserAuthned = info.UserAuthned
	tbl.RxWan0 = info.RxWan0
	tbl.TxWan0 = info.TxWan0
	tbl.OfflineTimesWAN0 = info.OfflineTimesWAN0
	tbl.RxEth0 = info.RxEth0
	tbl.TxEth0 = info.TxEth0
	tbl.TxSpeedEth0 = info.TxSpeedEth0
	tbl.RxSpeedEth0 = info.RxSpeedEth0
	tbl.RxAth0 = info.RxAth0
	tbl.TxAth0 = info.TxAth0
	tbl.RxAth1 = info.RxAth1
	tbl.TxAth1 = info.TxAth1
	tbl.RxAth2 = info.RxAth2
	tbl.TxAth2 = info.TxAth2
	tbl.ChAth0 = info.ChAth0
	tbl.ChAth1 = info.ChAth1
	tbl.ChAth2 = info.ChAth2
	tbl.TxPwrAth0 = info.TxPwrAth0
	tbl.TxPwrAth1 = info.TxPwrAth1
	tbl.TxPwrAth2 = info.TxPwrAth2

	tbl.LinkQualityAth0 = info.LinkQualityAth0
	tbl.SignalLevelAth0 = info.SignalLevelAth0
	tbl.NoiseLevelAth0 = info.NoiseLevelAth0
	tbl.UserConnectedAth0 = info.UserConnectedAth0
	tbl.UserAuthnedAth0 = info.UserAuthnedAth0
	tbl.UserVagSpeedAth0 = info.UserVagSpeedAth0

	tbl.LinkQualityAth1 = info.LinkQualityAth1
	tbl.SignalLevelAth1 = info.SignalLevelAth1
	tbl.NoiseLevelAth1 = info.NoiseLevelAth1
	tbl.UserConnectedAth1 = info.UserConnectedAth1
	tbl.UserAuthnedAth1 = info.UserAuthnedAth1
	tbl.UserVagSpeedAth1 = info.UserVagSpeedAth1

	tbl.LinkQualityAth2 = info.LinkQualityAth2
	tbl.SignalLevelAth2 = info.SignalLevelAth2
	tbl.NoiseLevelAth2 = info.NoiseLevelAth2
	tbl.UserConnectedAth2 = info.UserConnectedAth2
	tbl.UserAuthnedAth2 = info.UserAuthnedAth2
	tbl.UserVagSpeedAth2 = info.UserVagSpeedAth2

	tbl.UsedRateCpu = info.UsedRateCpu
	tbl.UsedRateMemory = info.UsedRateMemory
	tbl.UsedRateFlash = info.UsedRateFlash
}


func updateTblDevStat(info *common.DevStatus)error{
	var search common.TblDevStat
	search.Mac = info.Mac
	r, err := mysql_driver.Handle.FindFirst(search)
	if (err == nil) && (r.(common.TblDevStat) != common.TblDevStat {}){
		update := updateDevState(info)
		err := mysql_driver.Handle.Update(search,update)
		if err != nil{
			aclog.Error("Update TblDevState[%s]error:%s",search.Mac,err.Error())
			return err
		}
	}else {
		copyInfoToDevStateTbl(info, &search)
		err := mysql_driver.Handle.Insert(search)
		if err != nil{
			aclog.Error("Insert TblDevState[%s]error:%s",search.Mac,err.Error())
			return err
		}
	}
	return nil
}

func updateTblDevLog(info * common.DevStatus) error{
	var point []influxdb_driver.InfluxPoints
	var influx influxdb_driver.InfluxPoints
	tag := make(map[string]string)
	filed := make(map[string]interface{})
	tag["Mac"] = info.Mac
	influx.TableName = "TblDevLog"
	influx.Tag = tag
	filed["WirelessUserCount"] = info.WirelessUserConnectedCount
	filed["WiredUserCount"] = info.WiredUserConnectedCount
	filed["RxEth0"] = uint(info.RxEth0)
	filed["TxEth0"] = uint(info.TxEth0)
	filed["RxAth"] = uint(info.RxAth0 + info.RxAth1 + info.RxAth2)
	filed["TxAth"] = uint(info.TxAth0 + info.TxAth1 + info.TxAth2)
	filed["UsedRateCpu"] = info.UsedRateCpu
	filed["UsedRateMemory"] = info.UsedRateMemory
	filed["UsedRateFlash"] = info.UsedRateFlash
	influx.Field = filed
	point = append(point, influx)
	err := influxdb_driver.Service.Insert(point)
	if err != nil{
		return  err
	}
	return nil
}

func updateDevDB(info *common.DevStatus){
	v, err := redis_driver.RedisDb.IsKeyExit(info.Mac)
	if err != nil || v == 0{
		aclog.Warning("AP [%s] unregistered",info.Mac)
		return
	}
	err = updateTblDevStat(info)
	if err != nil{
		aclog.Error("Update table dev stat error:%s",err)
		return
	}
	err = updateTblDevLog(info)
	if err != nil{
		aclog.Error("Update table dev log error:%s",err.Error())
		return
	}
}

func HandleDevStateMsg(client mqtt.Client, message mqtt.Message){
	var info common.DevStatus
	err := json.Unmarshal([]byte(message.Payload()),&info)
	if err != nil{
		aclog.Error("Parse dev state msg error :%s",err.Error())
		return
	}
	go updateDevDB(&info)
}

func APStateReportStart() error{
	err := mqtt_driver.MqttHandle.SubscribeTopic(common.DevStateReportTopic,0,HandleDevStateMsg)
	if err != nil{
		aclog.Error("Subscribe topic:%s error will exit process.",common.DevStateReportTopic)
		panic(err)
	}
	return  nil
}
