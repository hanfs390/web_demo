package model

type ModelSupport struct {
	Model string  /* for display */
	DevName string /* real name for device */

	/* hardware info */
	Cpu string
	FlashType string
	FlashSize int
	MemSize int

	Radio int /* 1 - 2G; 2 - 5G; 4 - 5G2; */
	PortNumber int /* the number of the LANs */
	Chan2G []int
	Chan5G []int
	Chan5G2 []int
	SoftwareFlag string
}
var ModelSupportList = []ModelSupport{
	ModelSupport{
		"PatrolFlow-Air-GP530",
		"WitFi-CAP528-G",
		"QCA9563",
		"NOR",
		32,
		128,
		3,
		4,
		[]int{1,2,3,4,5,6,7,8,9,10,11,12,13},
		[]int{36,40,44,48,52,56,60,64, 149,153,157,161,165},
		[]int{},
		"w",
	},
	ModelSupport{
		"PatrolFlow-Air-GP830",
		"WitMAX-AP830-G",
		"IPQ4019",
		"NAND",
		128,
		256,
		7,
		1,
		[]int{1,2,3,4,5,6,7,8,9,10,11,12,13},
		[]int{36,40,44,48,52,56,60,64},
		[]int{149,153,157,161,165},
		"n",
	},
	ModelSupport{
		"PatrolFlow-Air-GP630",
		"WitMAX-AP630-G",
		"IPQ4019",
		"NAND",
		128,
		256,
		7,
		4,
		[]int{1,2,3,4,5,6,7,8,9,10,11,12,13},
		[]int{36,40,44,48,52,56,60,64},
		[]int{149,153,157,161,165},
		"n",
	},
}
func CheckModelIsExist(model string) bool {
	for i :=0; i < len(ModelSupportList); i++ {
		if ModelSupportList[i].Model == model {
			return true
		}
	}
	return false
}
func GetBoardInfoByModel(model string) *ModelSupport {
	for i :=0; i < len(ModelSupportList); i++ {
		if ModelSupportList[i].Model == model {
			return &ModelSupportList[i]
		}
	}
	return nil
}
