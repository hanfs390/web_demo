
package main

import (
	"ByzoroAC/aclog"
	"ByzoroAC/conf"
	"ByzoroAC/controllers/task"
	"ByzoroAC/models/influxdb_driver"
	"ByzoroAC/models/mysql_driver"
	"ByzoroAC/models/redis_driver"
	"ByzoroAC/routes"
)

func main() {
	conf.ReadGlobalConf()
	aclog.Init()
	mysql_driver.ModuleInit() /* init mysql tables */
	//mqtt_driver.InitMqtt()
	redis_driver.InitRedis()
	influxdb_driver.Init()
/*	reg.RegisterProcessStart()
	config.CfgCheckProcessStart()
	devState.APStateReportStart()
	gpon.GponReportStart()
	heartbeat.HeartBeatProcessStart()
	upgrade.UpgradeProcessStart()
	user.UserReportStart()*/
	go task.ReadTask()
	routes.Init()
}
