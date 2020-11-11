package influxdb_driver

import (
	"fmt"
	"net/url"
	"time"
	"github.com/influxdata/influxdb1-client/v2"

	"errors"
	"ByzoroAC/conf"
)


const (
	databaseName = "acDB"
)

type serviceProxy struct {
	influx   client.Client
}
type InfluxPoints struct {
	TableName string
	Tag map[string]string
	Field map[string]interface{}
}
type influxConf struct {
	host string
	port string
	user string
	password string
}
var Service *serviceProxy

/**
 * function: query data from Table
 *
 * Tables Name: two_hours.TblDevLog; two_days.TblDevLog_1d; one_month.TblDevLog_1m;
 *              one_month.TblUser; half_year.TblOnOffLog
 */
func (s *serviceProxy) Select(tableName string, key string, value string) ([]client.Result, error) {
	q := client.NewQuery(fmt.Sprintf("SELECT * from %s WHERE %s='%s'", tableName, key, value), databaseName, "s")
	response, err := Service.influx.Query(q)
	if err == nil && response.Error() == nil {
		return response.Results, nil
	}
	return nil, err
}

func (s *serviceProxy) Insert(points []InfluxPoints) error {
	var rp string
	if len(points) <= 0 {
		return errors.New("empty data")
	}
	if points[0].TableName == "TblDevLog" {
		rp = "two_hours"
	} else if points[0].TableName == "TblUser" {
		rp = "one_month"
	} else if points[0].TableName == "TblOnOffLog" {
		rp = "half_year"
	} else {
		return errors.New("error table name")
	}
	now := time.Now()
	fmt.Println(now.Unix())
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Precision:"s",
		Database:  databaseName,
		RetentionPolicy:rp,
	})

	for i := 0; i < len(points); i++ {
		point, _ := client.NewPoint(
			points[i].TableName,
			points[i].Tag,
			points[i].Field,
			now)
		bp.AddPoint(point)
	}

	err := Service.influx.Write(bp)
	return err
}

func (s *serviceProxy) Close() {
	s.influx.Close()
}
/**
 * initInflux - create a connect to influxDB
 */
func initInflux() error {
	confInflux := influxConf{}
	if conf.GlobalConf["InfluxDBHost"] != "" {
		confInflux.host = conf.GlobalConf["InfluxDBHost"]
	} else {
		confInflux.host = "127.0.0.1"
	}
	if conf.GlobalConf["InfluxDBPort"] != "" {
		confInflux.port = conf.GlobalConf["InfluxDBPort"]
	} else {
		confInflux.port = "8086"
	}
	if conf.GlobalConf["InfluxDBUserName"] != "" {
		confInflux.user = conf.GlobalConf["InfluxDBUserName"]
	} else {
		confInflux.user = "root"
	}
	if conf.GlobalConf["InfluxDBPassword"] != "" {
		confInflux.password = conf.GlobalConf["InfluxDBPassword"]
	} else {
		confInflux.password = "123456"
	}
	urlInflux := fmt.Sprintf("influx://%s:%s@%s:%s", confInflux.user, confInflux.password, confInflux.host, confInflux.port)
	opt, err := url.Parse(urlInflux)
	if err != nil {
		return err
	}
	pwd, _ := opt.User.Password()
	Service.influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s", opt.Host),
		Username: opt.User.Username(),
		Password: pwd,
		Timeout:  time.Second * 5,
	})
	if err != nil {
		return err
	}
	/* create the database */
	createDbSQL := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", databaseName), "", "s")
	if _, err := Service.influx.Query(createDbSQL); err != nil {
		return err
	}
	fmt.Println("influxDB services started")
	return nil
}
/**
 * initRetentionPolicy - create the Retention Policies
 * RPs: two_hours, two_days, one_month, one_year
 */
func initRetentionPolicies() error {
	rp_2h := client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY two_hours ON %s DURATION 2h REPLICATION 1", databaseName), databaseName, "")
	if _, err := Service.influx.Query(rp_2h); err != nil {
		return err
	}
	rp_2d:= client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY two_days ON %s DURATION 2d REPLICATION 1", databaseName), databaseName, "s")
	if _, err := Service.influx.Query(rp_2d); err != nil {
		return err
	}
	rp_30d:= client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY one_month ON %s DURATION 30d REPLICATION 1", databaseName), databaseName, "s")
	if _, err := Service.influx.Query(rp_30d); err != nil {
		return err
	}
	rp_180d:= client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY half_year ON %s DURATION 180d REPLICATION 1", databaseName), databaseName, "s")
	if _, err := Service.influx.Query(rp_180d); err != nil {
		return err
	}
	return nil
}
/**
 * initContinuousQueries - create the Continuous Queries
 * CQs: cq_1h, cq_1d
 */
func initContinuousQueries() error {
	cq_1h := client.NewQuery(fmt.Sprintf("CREATE CONTINUOUS QUERY cq_1h ON %s BEGIN" +
												" SELECT mean(WirelessUserCount) AS WirelessUserCount,mean(WiredUserCount) AS WiredUserCount," +
														"mean(UsedRateCpu) AS UsedRateCpu,mean(UsedRateMemory) AS UsedRateMemory," +
														"mean(UsedRateFlash) AS UsedRateFlash," +
														"last(RxEth0) AS RxEth0,last(TxEth0) AS TxEth0," +
														"last(RxAth) AS RxAth,last(TxAth) AS TxAth " +
												" INTO two_days.TblDevLog_1d " +
												" FROM two_hours.TblDevLog " +
												" GROUP BY time(1h) END", databaseName), databaseName, "s")
	if _, err := Service.influx.Query(cq_1h); err != nil {
		fmt.Println(err)
		return err
	}
	cq_1d := client.NewQuery(fmt.Sprintf("CREATE CONTINUOUS QUERY cq_1d ON %s BEGIN" +
												" SELECT mean(WirelessUserCount) AS WirelessUserCount,mean(WiredUserCount) AS WiredUserCount," +
													"mean(UsedRateCpu) AS UsedRateCpu,mean(UsedRateMemory) AS UsedRateMemory," +
													"mean(UsedRateFlash) AS UsedRateFlash," +
													"last(RxEth0) AS RxEth0,last(TxEth0) AS TxEth0," +
													"last(RxAth) AS RxAth,last(TxAth) AS TxAth " +
												" INTO one_month.TblDevLog_1m " +
												" FROM two_days.TblDevLog_1d " +
												" GROUP BY time(1d) END", databaseName), databaseName, "s")
	if _, err := Service.influx.Query(cq_1d); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func Init() error {
	Service = new(serviceProxy)
	err := initInflux()
	if err != nil {
		return err
	}

	err = initRetentionPolicies()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = initContinuousQueries()
	if err != nil {
		return err
	}
	return nil
}