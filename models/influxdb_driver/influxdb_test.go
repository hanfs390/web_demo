package influxdb_driver

import (
	"testing"
	"fmt"
	"ByzoroAC/aclog"
	"ByzoroAC/conf"
)

func TestMain(m *testing.M) {
	aclog.Init()
	conf.ReadGlobalConf()
	Init()
	m.Run()
}
func insertTest(t *testing.T) {
	var points []InfluxPoints
	tag := make(map[string]string)
	tag["Mac"] = "12345678"
	filed := make(map[string]interface{})
	filed["WirelessUserCount"] = 12
	filed["WiredUserCount"] = 6
	point := InfluxPoints{TableName:"TblDevLog",Tag:tag,Field:filed}
	points = append(points, point)
	err := Service.Insert(points)
	if err != nil {
		fmt.Println(err)
	}
}
func queryTest(t *testing.T) {
	result, err := Service.Select("two_hours.TblDevLog", "Mac", "12345678")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i].Series); j++ {
			fmt.Println(result[i].Series[j].Name) /* table name */
			fmt.Println(result[i].Series[j].Tags) /* map[] */
			fmt.Println(result[i].Series[j].Columns) /* tags and fields */
			fmt.Println(result[i].Series[j].Values) /* array value */
			fmt.Println(result[i].Series[j].Partial) /* unknown */
		}
	}
}
func TestAll(t *testing.T) {
	//t.Run("Insert", insertTest)
	t.Run("Query", queryTest)
}
