/**
 * the struct for mysql tables
 */

package common
/*
 * string is varchar255
 */
type TblGroup struct {
	Id uint64
	Uuid string
	Url string
	Name string
	Detail string
	Tag string
	GroupCfg string
	Rsv1 string
	Rsv2 string
	Rsv3 int
	Others string
}

type TblDevInfo struct {
	Id uint64
	Mac string `gorm:"size:32"`
	Sn string `gorm:"size:32"`
	Model string `gorm:"size:32"` //产品型号，用于页面展示
	DevName string `gorm:"size:32"` // 对内产品型号用于区分不同的产品
	FirmwareVer string `gorm:"size:32"` //软件版本号
	HwVer string `gorm:"size:32"` //产品硬件版本
	HwUnique string `gorm:"size:32"` //硬件唯一标识，从硬件读取的不可更改唯一标识
	/* 为其他公司贴牌预留 */
	VendorModel string `gorm:"size:32"`
	VendorVer string `gorm:"size:32"`
	VendorSn string `gorm:"size:32"`
	VendorMac string `gorm:"size:32"`

	HwInfo string //硬件信息，JSON格式。与型号关联
	Alias string //用户自定义名称
	Detail string //用户自定义描述
	Blocked int //禁止注册标志（1为禁止注册）
	Tag string //用户字定义标签组，是一个关键词列表
	SignUpType string `gorm:"size:32"`//设备录入的方式，暂时分为：
										//import 人工导入
										//register 设备注册
	SignUpTime string `gorm:"size:64"` //gorm:"size:32"
	FirstLoginTime string `gorm:"size:64"` //首次登陆时间

	GroupId uint64
	GroupName string
	GroupUrl string

	TargetFirmwareVer string `gorm:"size:32"` //目标软件版本
	DevScript string //需要设备获取并执行的脚本路径
	LBName string //AP所在负载均衡组的名称，默认为空
	DevCfg string `gorm:"size:1024"` //设备配置，每个设备的配置有一个独立的配置文件
	Rsv1 string
	Rsv2 string
	Rsv3 int
	Others string `gorm:"size:1024"`
}

type TblDevStat struct {
	Id uint64
	Mac string `gorm:"size:32"`
	State string `gorm:"size:32"`

	LastLogInTime string `gorm:"size:64"`
	LastLogOutTime string `gorm:"size:64"`
	DevBootTime string `gorm:"size:64"`
	DevOfflineTimes string `gorm:"size:64"`

	FirmwareVer string `gorm:"size:32"`
	Ip string `gorm:"size:32"`
	UserConnected uint //连接用户数量
	WiredUserConnectedCount uint
	WirelessUserConnectedCount uint
	UserAuthned uint //reserve

	GroupId uint64
	GroupUrl string
	GroupName string

	RxWan0 uint64
	TxWan0 uint64
	OfflineTimesWAN0 int
	/* reserve */
	RxWan1 uint64
	TxWan1 uint64
	OfflineTimesWAN1 int
	/* out interface */
	RxEth0 uint64
	TxEth0 uint64
	TxSpeedEth0 int
	RxSpeedEth0 int
	/*reserve*/
	RxEth1 uint64
	TxEth1 uint64
	RxEth2 uint64
	TxEth2 uint64
	RxEth3 uint64
	TxEth3 uint64

	RxAth0 uint64
	TxAth0 uint64
	RxAth1 uint64
	TxAth1 uint64
	RxAth2 uint64
	TxAth2 uint64
	RxAth3 uint64
	TxAth3 uint64

	ChAth0 int
	ChAth1 int
	ChAth2 int
	ChAth3 int

	TxPwrAth0 int
	TxPwrAth1 int
	TxPwrAth2 int
	TxPwrAth3 int

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

	LinkQualityAth3 int
	SignalLevelAth3 int
	NoiseLevelAth3 int
	UserConnectedAth3 int
	UserAuthnedAth3 int
	UserVagSpeedAth3 int
	/* usage rate: 0 ~ 99 */
	UsedRateCpu int
	UsedRateMemory int
	UsedRateFlash int

	Rsv0 int
	Rsv1 int
	Rsv2 int
	Rsv3 int
	Rsv4 string
	Rsv5 string
	Others string `gorm:"size:1024"`
}

type TblWlan struct {
	Id uint64

	GroupId uint64
	GroupName string
	GroupUrl string

	Ssid string
	Disabled int //该WLAN是否启用。1为关闭，0为启用
	Encryption string `gorm:"size:32"`
	Hidden string `gorm:"size:32"`
	VlanSwitch int
	VlanId uint
	Radio string `gorm:"size:64"`

	PerUserRate int
	UpPerUserRate int
	DownPerUserRate int

	Key string
	AuthSecret string
	AuthServer string `gorm:"size:32"`
	AuthPort string `gorm:"size:32"`
	MultServer int
	Servercfg string `gorm:"size:1024"`
	FirstRoamcfg string `gorm:"size:1024"`
}

type TblLED struct {
	Id uint64
	GroupId uint64
	GroupName string
	GroupUrl string
	Timer string
}

type TblLoadBalanceGroup struct {
	Id uint64
	GroupId uint64
	GroupName string
	GroupUrl string
	Name string
	Member string
}

type TblPolicy struct {
	Id uint64
	GroupId uint64
	GroupName string
	GroupUrl string

	Dual_max_user int
	Single_max_user int
	Rssi_threshold int
	Access_policy int
	Reject_max int
	Rssi_max int
	L2_isolation int
	Band_steering int
	Thredhold_5g int
	Thredhold_5g_rssi int
	Roaming_detect int
	Roaming_assoc_rssi int
}

type TblBWList struct {
	Id uint64
	GroupId uint64
	GroupName string
	GroupUrl string
	Type string
	BlackList string `gorm:"size:10240"`
	WhiteList string `gorm:"size:10240"`
}
