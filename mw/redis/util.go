package redis

import (
	"strconv"

	"github.com/go-redis/redis/v7"
)

// add key & value to redis
func add(c *redis.Client, key string, value int64) {
	tx := c.TxPipeline()
	tx.SAdd(key, value)
	tx.Expire(key, expireTime)
	_, err := tx.Exec()
	if err != nil {
		panic(err)
	}
}

// delete key & value from redis
func del(c *redis.Client, key string, value int64) {
	tx := c.TxPipeline()
	tx.SRem(key, value)
	_, err := tx.Exec()
	if err != nil {
		panic(err)
	}
}

// check if key exists
func check(c *redis.Client, key string) bool {
	res := c.Exists(key)
	if res.Val() == 0 {
		return false
	}
	return true
}

// check the relation of k and v if exists
func exist(c *redis.Client, k string, v int64) bool {
	if e, _ := c.SIsMember(k, v).Result(); e {
		c.Expire(k, expireTime)
		return true
	}
	return false

}

// count get the size of the set of key
func count(c *redis.Client, k string) (sum int64, err error) {
	if sum, err = c.SCard(k).Result(); err == nil {
		c.Expire(k, expireTime)
		return sum, err
	}
	return sum, err
}

func get(c *redis.Client, k string) (vt []int64) {
	v, _ := c.SMembers(k).Result()
	c.Expire(k, expireTime)
	for _, vs := range v {
		v_i64, _ := strconv.ParseInt(vs, 10, 64)
		vt = append(vt, v_i64)
	}
	return vt
}
