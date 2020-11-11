package common

import "time"

const(
	RegTopic  = "v1/ap/telemetry/register"
	CfgCheckTopic = "v1/ap/cfg/polling"
	DevStateReportTopic = "v1/ap/telemetry/state"
	GponStateTopic = "v1/ap/telemetry/gpon"
	HeartBeatTopic = "v1/ap/heart/polling"
	UpgradeTopic = "v1/ap/upgrade/polling"
	UserReportTopic = "v1/ap/telemetry/user"
 	RegisterMethod  = "registerStatus"
	CfgPollingMethod = "configCheck"
	WlanCfgMethod  = "updateWlan"
	PolicyCfgMethod = "updatePolicy"
	LEDCfgMethod  = "updateLED"
 	BWCfgMethod = "updateBWList"
	LBCfgMethod = "updateLB"
	UpgradeMethod  = "upgradeFirmware"
	DevCfgMethod  = "updateDevCfg"
	UpdateAllMethod = "updateAllCfg"
	ResponseVersion = "v1"
	KeyExpireTime = 30 * time.Second
)

var APMethod map[string]string

type RegInfo struct {
	Mac string
	Sn string
	Model string
	HwVer string
	HwUnique string
	DevName string
	Version string

}

type TblAllCfg struct {
	Mac  string
	DevCfg string
	DevCfgMD5 string
	GroupUrl string
	DevScript string
	TargetFirmwareVer string
	LBName string
	GroupName  string
	GroupId uint64
	Gpon string
	UpgradeUrl string
}

type CfgPolling struct {
	Mac string
	DevCfgMd5 string
	WlanMd5 string
	PolicyMd5 string
	LEDMd5 string
	BWMd5 string
	LBMd5 string
}

type UpdateCfg struct {
	DevCfg string
	Wlan string
	Policy string
	LED string
	BW string
	LB string
}

type ResponseInfo struct{
	Version string
	Method string
	Contents string
}

type DevStatus struct {
	Mac string
	DevBootTime string
	FirmwareVer string
	Ip string
	UserConnected uint
	WiredUserConnectedCount uint
	WirelessUserConnectedCount uint
	UserAuthned uint

	RxWan0 uint64
	TxWan0 uint64
	OfflineTimesWAN0 int

	RxEth0 uint64
	TxEth0 uint64
	TxSpeedEth0 int
	RxSpeedEth0 int
	RxAth0 uint64
	TxAth0 uint64
	RxAth1 uint64
	TxAth1 uint64
	RxAth2 uint64
	TxAth2 uint64
	ChAth0 int
	ChAth1 int
	ChAth2 int
	TxPwrAth0 int
	TxPwrAth1 int
	TxPwrAth2 int
	LinkQualityAth0 int
	SignalLevelAth0 int
	NoiseLevelAth0 int
	UserConnectedAth0 int
	UserAuthnedAth0 int
	UserVagSpeedAth0 int
	LinkQualityAth1 int
	SignalLevelAth1 int
	NoiseLevelAth1 int
	UserConnectedAth1 int
	UserAuthnedAth1 int
	UserVagSpeedAth1 int
	LinkQualityAth2 int
	SignalLevelAth2 int
	NoiseLevelAth2 int
	UserConnectedAth2 int
	UserAuthnedAth2 int
	UserVagSpeedAth2 int
	/* usage rate: 0 ~ 99 */
	UsedRateCpu int
	UsedRateMemory int
	UsedRateFlash int
}

type GponStatus struct {
	Mac string
	GponState string
	OffReason string
	GponSn string
	GponPwd string
	LosStatus string
	TxPower string
	RxPower string
	Temperature string
	SupplyVoltage string
	TxbiasCurrent string
	OnuState string
	PhyStatus string
	TrafficStatus string
	Manufacturer string
	ManufacturerOui string
	OperatorId string
	ModelName string
	CustomerHwversion string
	CustomerSwversion string
}

type HeartBeat struct {
	Mac string
}

type UpgradeInfo struct {
	Mac string
	TargetFirmwareVer string
}

type WirelessUserInfo struct {
	Mac string
	Ip string
	Wlan string
	Radio string
	Channel int
	DevModel string
	DevType string
	OSType string
	Signal int
	TxRate int
	TxBytes int
	RxRate int
	RxBytes int
	HostName string
	LoginTime string
}

type WiredUserInfo struct {
	Mac string
	LoginTime string
	Port int
}

type UserReport struct {
	APMac string
	WirelessUser [] WirelessUserInfo
	WiredUser [] WiredUserInfo
}