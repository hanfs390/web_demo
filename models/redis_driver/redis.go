package redis_driver

import (
	"ByzoroAC/aclog"
	"ByzoroAC/conf"
	"bytes"
	"github.com/gomodule/redigo/redis"
	"time"
)

type redisDb struct {
	proto string
	addr string
	pool *redis.Pool
}

var RedisDb  *redisDb

func( handle *redisDb) initRedisPool(){
	handle.pool = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial(
				handle.proto,
				handle.addr,
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialDatabase(0),
				)
		},
		TestOnBorrow:    nil,
		MaxIdle:         100,
		MaxActive:       1000,
		IdleTimeout:     0,
		Wait:            true,
		MaxConnLifetime: 0,
	}
}

func appendArg(dst []interface{}, arg interface{}) []interface{} {
	switch arg := arg.(type) {
	case []string:
		for _, s := range arg {
			dst = append(dst, s)
		}
		return dst
	case []interface{}:
		dst = append(dst, arg...)
		return dst
	case map[string]interface{}:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	default:
		return append(dst, arg)
	}
}

func appendArgs(dst, src []interface{}) []interface{} {
	if len(src) == 1 {
		return appendArg(dst, src[0])
	}

	dst = append(dst, src...)
	return dst
}
//BatchHashSet values in following formats:
//   -BatchHashSet("myhash", "key1", "value1", "key2", "value2")
//   - BatchHashSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - BatchHashSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (handle *redisDb)BatchHashSet(key string, valus ...interface{}) error {
	con := handle.pool.Get()
	if err := con.Err();err != nil{
		return err
	}
	defer con.Close()
	args := make([]interface{},1,1+len(valus))
	args[0] = key
	args = appendArgs(args, valus)
	_, err := redis.String(con.Do("HMSET",args...))
	if err != nil{
		return err
	}
	return nil
}

func(handle *redisDb)BatchHashGet(key string, fileds ...string) (val []interface{},error error){
	var res [] interface{}
	con := handle.pool.Get()
	if err := con.Err(); err != nil{
		return res,err
	}
	defer con.Close()
	args := make([]interface{},1+len(fileds))
	args[0] = key
	for i, filed := range fileds{
		args[1+i] = filed
	}
	res, err := redis.Values(con.Do("HMGET", args...))
	if err != nil{
		return res, err
	}
	return  res,err
}

func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		aclog.Warning(
			"specified duration is %s, but minimal supported value is %s - truncating to 1s",
			dur,time.Second,
			)
		return 1
	}
	return int64(dur / time.Second)
}

func(handle *redisDb)ExpireHashKey(key string, expiration time.Duration)error{
	con := handle.pool.Get()
	if err := con.Err(); err != nil{
		return err
	}
	defer con.Close()
	_, err := con.Do("Expire",key,formatSec(expiration))
	if err != nil{
		return  err
	}
	return nil
}

func(handle *redisDb)DeleteKey(key string)error{
	con := handle.pool.Get()
	if err := con.Err();err != nil{
		return  err
	}
	defer con.Close()
	_, err := con.Do("DEL",key)
	if err != nil{
		return  err
	}
	return  nil
}
// return value 1 exit
// 0 not exit
func(handle *redisDb)IsKeyExit(key string)(res int, error error){
	var v int
	con := handle.pool.Get()
	if err := con.Err(); err != nil{
		return v, err
	}
	defer con.Close()
	v, err := redis.Int(con.Do("EXISTS",key))
	if err != nil{
		return v,err
	}
	return v,nil
}

func(hadle *redisDb)Set(key string, value string) error{
	con := hadle.pool.Get()
	if err := con.Err(); err != nil{
		return  err
	}
	defer con.Close()
	_, err := con.Do("SET", key, value)
	if err != nil{
		return  err
	}
	return nil
}

func (handle *redisDb)Get(key string) (res string, err error){
	con := handle.pool.Get()
	if err := con.Err(); err != nil{
		return  "", err
	}
	defer con.Close()
	v, err := redis.String(con.Do("GET",key))
	if err != nil{
		return "",err
	}
	return v,nil
}

func (handle *redisDb)GetProtoAndAddr()(proto, addr string){
	return handle.proto,handle.addr
}

func InitRedis(){
	var host, port string
	var addr bytes.Buffer
	if conf.GlobalConf["RedisHost"] != ""{
		host = conf.GlobalConf["RedisHost"]
	}else{
		host = "127.0.0.1"
	}
	if conf.GlobalConf["RedisPort"] != ""{
		port = conf.GlobalConf["RedisPort"]
	}else {
		 port = "6379"
	}
	addr.WriteString(host)
	addr.WriteString(":")
	addr.WriteString(port)
	RedisDb = new(redisDb)
	RedisDb.proto = "tcp"
	RedisDb.addr = addr.String()
	RedisDb.initRedisPool()
}


