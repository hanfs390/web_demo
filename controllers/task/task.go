package task

import (
	"time"
	"ByzoroAC/models/redis_driver"
	"fmt"
)

type Task struct {
	Op string
	Key string
	Value interface{}
	Time int64
}
var TaskRedis []Task

func ReadTask() {
	for {
		current := time.Now().Unix()

		for i := 0; i < len(TaskRedis); i++ {
			var err error
			if (current - TaskRedis[i].Time) < 10 {
				continue
			}
			switch TaskRedis[i].Op {
			case "Del":
				err = redis_driver.RedisDb.DeleteKey(TaskRedis[i].Key)
			case "Set":
				err = redis_driver.RedisDb.Set(TaskRedis[i].Key, TaskRedis[i].Value.(string))
			case "HashSet":
				err = redis_driver.RedisDb.BatchHashSet(TaskRedis[i].Key, TaskRedis[i].Value)
			}
			if err != nil {
				fmt.Println(TaskRedis[i].Key,"err:", err)
			} else {
				TaskRedis = append(TaskRedis[:i], TaskRedis[i+1:]...)
				i--
			}
		}
		time.Sleep(time.Second)
	}
}