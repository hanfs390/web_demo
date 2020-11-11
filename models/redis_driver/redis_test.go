package redis_driver

import (
	"ByzoroAC/aclog"
	"ByzoroAC/conf"
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M){
	aclog.Init()
	conf.ReadGlobalConf()
	InitRedis()
	m.Run()
}

func TestRedisDb_BatchHashSet(t *testing.T){
	g := make(map[string]interface{})
	g["Mac"] = "112233445566"
	g["Name"] = "Test"
	g["age"] = 18
	key := "112233445566"
	err := RedisDb.BatchHashSet(key,g)
	if err != nil{
		fmt.Println(err)
		return
	}
}

func TestRedisDb_BatchHashGet(t *testing.T) {
	v, err := RedisDb.BatchHashGet("112233445566","age","Name")
	if err != nil{
		fmt.Println(err)
		return
	}
	for _, res := range v{
		fmt.Printf("%s\n",res)
	}
}

func TestRedisDb_ExpireHashKey(t *testing.T) {
	err := RedisDb.ExpireHashKey("112233445566",10*time.Second)
	if err != nil{
		fmt.Println(err)
	}
}

func TestRedisDb_DeleteHashKey(t *testing.T) {
	err := RedisDb.DeleteKey("112233445566")
	if err != nil{
		fmt.Println(err)
	}
}

func TestRedisDb_IsKeyExit(t *testing.T) {
	v, err := RedisDb.IsKeyExit("112233445566")
	if err != nil{
		fmt.Println(err)
		return
	}
	v1, err := RedisDb.IsKeyExit("11223344")
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(v,v1)
}

func TestAll(t *testing.T){
	t.Run("BatchHashSet",TestRedisDb_BatchHashSet)
	t.Run("BatchHashGet",TestRedisDb_BatchHashGet)
	t.Run("ExpireHashKey",TestRedisDb_ExpireHashKey)
	t.Run("DeleteHashKey",TestRedisDb_DeleteHashKey)
	t.Run("IsKeyExit",TestRedisDb_IsKeyExit)
}