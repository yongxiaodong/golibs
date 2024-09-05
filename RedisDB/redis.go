package RedisDB

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	mu                      sync.Mutex
	redisStringDataCacheTag sync.Map
	blockNum                int64
)

type ConnRedisParams struct {
	Addr     string
	Password string
	DB       int
}

func validateOpts(opts *redis.Options) {
	if opts.MinIdleConns == 0 || opts.MinIdleConns > 2000 {
		opts.MinIdleConns = 50
	}
	if opts.PoolSize == 0 || opts.PoolSize > 2000 {
		opts.PoolSize = 100
	}
	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = time.Minute * 5
	}
	if opts.PoolTimeout == 0 {
		opts.PoolTimeout = time.Minute * 10
	}

}

func NewRedisClient(opts *redis.Options) *redis.Client {
	// 校验opts
	validateOpts(opts)

	redisClient := redis.NewClient(opts)
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return redisClient

}

// 计算查询阻塞数量
func blockNumCounter(i int64) {
	atomic.AddInt64(&blockNum, i)
}

// 测试阻塞情况
func blockNumCounterDebug() {
	if blockNum < 5 {
		log.Println("阻塞:", blockNum)
	} else {
		remainder := blockNum % 1000
		if remainder == 0 {
			log.Println("阻塞:", blockNum)
		}
	}
}

func FromCacheOrDB(r *redis.Client, key string, f interface{}, cacheTime time.Duration) (interface{}, error) {
	ctx := context.Background()
	// 优先使用redis缓存
	getFromCache := func() (interface{}, error) {
		value, err := r.Get(ctx, key).Result()
		return value, err
	}
	value, err := getFromCache()
	if err == nil && value != "" {
		return value, nil
	} else if err != redis.Nil {
		return "", err
	}
	// 从DB数据获取
	if wgValue, ok := redisStringDataCacheTag.Load(key); ok {
		// 等待其它请求从DB获取数据写入缓存
		blockNumCounter(1)
		// 单元测试显示阻塞数据
		//blockNumCounterDebug()
		wgValue.(*sync.WaitGroup).Wait()
		value, err := getFromCache()
		blockNumCounter(-1)
		return value, err
	}
	// 运行函数从DB获取数据
	wg := &sync.WaitGroup{}
	wg.Add(1)
	redisStringDataCacheTag.Store(key, wg)
	result := f.(func() interface{})()
	// 写入缓存
	err = r.Set(ctx, key, result, cacheTime).Err()
	if err != nil {
		redisStringDataCacheTag.Delete(key)
		wg.Done()
		return "", err
	}
	redisStringDataCacheTag.Delete(key)
	wg.Done()
	return result, nil
}
