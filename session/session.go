package session

import (
	"encoding/json"
	. "nuobelcloud/conf"
	"nuobelcloud/redis"
	. "nuobelcloud/tools"
)

/*

session操作

*/

type SessionInfo struct {
	Sid        string
	Uid        int
	UserName   string
	SetDate    int64
	ExpireTime int
	Dir1       string
}

var provider *redis.RedisProvider

func init() {
	var err error
	provider, err = redis.NewRedisClient(RedisAddr, 1)

	if err != nil {
		panic(err)
		return
	}
	Log(INFO, "redis连接成功.")
}

func (c *SessionInfo) Set() error {

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	provider.Set(c.Sid, c.ExpireTime, string(b))

	return nil
}

func (c *SessionInfo) Get() error {

	valStr, err := provider.Get(c.Sid)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(valStr), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *SessionInfo) Del() {
	provider.Del(c.Sid)
}

func (c *SessionInfo) Update(oldKey string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	provider.Update(oldKey, c.Sid, c.ExpireTime, string(b))
	return nil
}
