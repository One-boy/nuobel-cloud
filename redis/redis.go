package redis

import (
	"github.com/garyburd/redigo/redis"
	. "nuobelcloud/tools"
	"sync"
)

/*

redis连接

*/

var MAX_POOL_SIZE = 200

type RedisProvider struct {
	poollist *redis.Pool
	lock     sync.RWMutex
}

func NewRedisClient(addr string, dbnum int) (*RedisProvider, error) {

	provider := new(RedisProvider)
	//新建一个pool
	provider.poollist = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", addr)
		if err != nil {
			Log(ERROR, "连接redis 失败，", err)
			return nil, err
		}
		_, err = c.Do("SELECT", dbnum)

		if err != nil {
			return nil, err
		}
		return c, nil
	}, MAX_POOL_SIZE)

	return provider, provider.poollist.Get().Err()
}

//根据key和超时时间设置值
func (c *RedisProvider) Set(key string, timeout int, val string) {
	rc := c.poollist.Get()
	//defer rc.Close()
	c.lock.Lock()
	defer c.lock.Unlock()
	//用setex命令重复设置相同的key，生存时间会被覆盖
	rc.Do("SETEX", key, timeout, val)
}

//获取key值
func (c *RedisProvider) Get(key string) (string, error) {
	rc := c.poollist.Get()
	//defer rc.Close()
	c.lock.RLock()
	defer c.lock.RUnlock()
	kvs, err := redis.String(rc.Do("GET", key))
	if err != nil {
		Log(ERROR, "redis get error:，", err, kvs)
		return "", err
	}

	return kvs, nil
}

//删除key值
func (c *RedisProvider) Del(key string) {
	rc := c.poollist.Get()
	//defer rc.Close()
	c.lock.Lock()
	defer c.lock.Unlock()
	rc.Do("DEL", key)
}

//更新key
func (c *RedisProvider) Update(oldkey, newKey string, timeout int, newVal string) {
	rc := c.poollist.Get()
	//defer rc.Close()
	c.lock.Lock()
	defer c.lock.Unlock()
	rc.Do("RENAME", oldkey, newKey)
	rc.Do("SETEX", newKey, timeout, newVal)
}
