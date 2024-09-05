package RedisDB

import (
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
	"time"
)

func mockQueryDB() interface{} {
	//fmt.Println("有请求开始查库")
	time.Sleep(time.Second * 3)
	//fmt.Println("查库 END")
	return "123"
}

func TestFromCacheOrDB(t *testing.T) {
	r := NewRedisClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "1317665590",
		DB:       0,
		PoolSize: 300,
	})
	var wg sync.WaitGroup
	stopChain := time.After(30 * time.Second)
	for {
		select {
		case <-stopChain:
			return
		case <-time.NewTicker(10).C:
			go func() {
				wg.Add(1)
				data, err := FromCacheOrDB(r, "test", mockQueryDB, time.Second*5)
				if err != nil {
					panic(err)
				}
				if data != "123" {
					panic("===")
				}
			}()
		}
	}
	wg.Wait()
}
