package redis

import (
	"errors"
	"strconv"

	rg "github.com/garyburd/redigo/redis"
)

func (con *connection) Do(command string, args ...interface{}) (interface{}, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Do(command, args...)
	} else {
		return con.c.Do(command, args...)
	}
}

func (con *connection) Send(command string, args ...interface{}) error {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Send(command, args...)
	} else {
		return con.c.Send(command, args...)
	}
}

func (con *connection) Flush() error {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Flush()
	} else {
		return con.c.Flush()
	}
}

func (con *connection) Receive() (interface{}, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Receive()
	} else {
		return con.c.Receive()
	}
}

/*
	Connection
*/
func (con *connection) Auth(password string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Auth(password)
	} else {
		res, err := rg.String(con.c.Do("AUTH", password))
		return getBool(res), err
	}
}

func (con *connection) Echo(message string) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Echo(message)
	} else {
		return rg.String(con.c.Do("ECHO", message))
	}
}

func (con *connection) Ping() (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Ping()
	} else {
		return rg.String(con.c.Do("PING"))
	}
}

func (con *connection) Select(index int) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Select(index)
	} else {
		res, err := con.c.Do("SELECT", index)
		return rg.Bool(res, err)
	}
}

func (con *connection) Quit() (string, error) {
	return rg.String(con.c.Do("QUIT"))
}

/*
	Hashes
*/
func (con *connection) HDel(hashKey string, fields []string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HDel(hashKey, fields)
	} else {
		req := make([]interface{}, len(fields)+1)
		req[0] = hashKey
		for idx, val := range fields {
			req[idx+1] = val
		}

		return rg.Int(con.c.Do("HDEL", req...))
	}
}

func (con *connection) HExists(hashKey, field string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HExists(hashKey, field)
	} else {
		return rg.Bool(con.c.Do("HEXISTS", hashKey, field))
	}
}

func (con *connection) HGet(hashKey, field string) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HGet(hashKey, field)
	} else {
		return rg.String(con.c.Do("HGET", hashKey, field))
	}
}

func (con *connection) HGetFloat64(hashKey, field string) (float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HGetFloat64(hashKey, field)
	} else {
		return rg.Float64(con.c.Do("HGET", hashKey, field))
	}
}

func (con *connection) HGetAll(hashKey string) (map[string]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HGetAll(hashKey)
	} else {
		return rg.StringMap(con.c.Do("HGETALL", hashKey))
	}
}

func (con *connection) HGetAllInterface(hashKey string) (map[string]interface{}, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HGetAllInterface(hashKey)
	} else {
		result, err := rg.Values(con.c.Do("HGETALL", hashKey))

		if err != nil {
			return nil, err
		}

		res := make(map[string]interface{})

		for key, value := range result {
			res[string(key)] = value
		}

		return res, nil
	}
}

func (con *connection) HGetAllFloat64(hashKey string) (map[string]float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HGetAllFloat64(hashKey)
	} else {
		result, err := rg.StringMap(con.c.Do("HGETALL", hashKey))

		if err != nil {
			return map[string]float64{}, err
		}

		res := make(map[string]float64)

		for key, value := range result {
			res[key] = strToFloat64(value)
		}

		return res, nil
	}
}

func (con *connection) HIncrBy(hashKey, field string, increment int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HIncrBy(hashKey, field, increment)
	} else {
		return rg.Int(con.c.Do("HINCRBY", hashKey, field, increment))
	}
}

func (con *connection) HIncrByFloat(hashKey, field string, increment float64) (float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HIncrByFloat(hashKey, field, increment)
	} else {
		return rg.Float64(con.c.Do("HINCRBYFLOAT", hashKey, field, increment))
	}
}

func (con *connection) HKeys(hashKey string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HKeys(hashKey)
	} else {
		return rg.Strings(con.c.Do("HKEYS", hashKey))
	}
}

func (con *connection) HLen(hashKey string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HLen(hashKey)
	} else {
		return rg.Int(con.c.Do("HLEN", hashKey))
	}
}

func (con *connection) HMGet(hashKey string, fields []string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HMGet(hashKey, fields)
	} else {
		req := make([]interface{}, len(fields)+1)
		req[0] = hashKey
		for idx, val := range fields {
			req[idx+1] = val
		}
		return rg.Strings(con.c.Do("HMGET", req...))
	}
}

func (con *connection) HMGetFloat64(hashKey string, fields []string) ([]float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HMGetFloat64(hashKey, fields)
	} else {
		req := make([]interface{}, len(fields)+1)
		req[0] = hashKey
		for idx, val := range fields {
			req[idx+1] = val
		}

		result, err := rg.Strings(con.c.Do("HMGET", req...))

		if err != nil {
			return nil, err
		}

		res := make([]float64, len(result))

		for idx, value := range result {
			res[idx] = strToFloat64(value)
		}

		return res, nil
	}
}

func (con *connection) HMSet(hashKey string, fieldValue map[string]string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HMSet(hashKey, fieldValue)
	} else {
		req := make([]interface{}, len(fieldValue)*2+1)
		req[0] = hashKey
		idx := 1
		for name, value := range fieldValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		res, err := rg.String(con.c.Do("HMSET", req...))
		return getBool(res), err
	}
}

func (con *connection) HMSetInterface(hashKey string, fieldValue map[string]interface{}) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HMSetInterface(hashKey, fieldValue)
	} else {
		req := make([]interface{}, len(fieldValue)*2+1)
		req[0] = hashKey
		idx := 1
		for name, value := range fieldValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		res, err := rg.String(con.c.Do("HMSET", req...))
		return getBool(res), err
	}
}

func (con *connection) HMSetFloat64(hashKey string, fieldValue map[string]float64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HMSetFloat64(hashKey, fieldValue)
	} else {
		req := make([]interface{}, len(fieldValue)*2+1)
		req[0] = hashKey
		idx := 1
		for name, value := range fieldValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		res, err := rg.String(con.c.Do("HMSET", req...))
		return getBool(res), err
	}
}

// func HScan() ()

func (con *connection) HSet(hashKey, field, value string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HSet(hashKey, field, value)
	} else {
		return rg.Int(con.c.Do("HSET", hashKey, field, value))
	}
}

func (con *connection) HSetFloat64(hashKey, field string, value float64) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HSetFloat64(hashKey, field, value)
	} else {
		return rg.Int(con.c.Do("HSET", hashKey, field, value))
	}
}

func (con *connection) HSetNX(hashKey, field, value string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HSetNX(hashKey, field, value)
	} else {
		return rg.Int(con.c.Do("HSETNX", hashKey, field, value))
	}
}

func (con *connection) HStrLen(hashKey, field string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HStrLen(hashKey, field)
	} else {
		return rg.Int(con.c.Do("HSTRLEN", hashKey, field))
	}
}

func (con *connection) HVals(hashKey string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.HVals(hashKey)
	} else {
		return rg.Strings(con.c.Do("HVALS", hashKey))
	}
}

/*
	Keys
*/
func (con *connection) Del(keys []string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Del(keys)
	} else {
		req := make([]interface{}, len(keys))
		for idx, val := range keys {
			req[idx] = val
		}
		return rg.Int(con.c.Do("DEL", req...))
	}
}

func (con *connection) Dump(key string) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Dump(key)
	} else {
		return rg.String(con.c.Do("DUMP", key))
	}
}

func (con *connection) Exists(key string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Exists(key)
	} else {
		return rg.Bool(con.c.Do("EXISTS", key))
	}
}

func (con *connection) Expire(key string, seconds int) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Expire(key, seconds)
	} else {
		return rg.Bool(con.c.Do("EXPIRE", key, seconds))
	}
}

func (con *connection) Expireat(key string, timestamp int64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Expireat(key, timestamp)
	} else {
		return rg.Bool(con.c.Do("EXPIREAT", key, timestamp))
	}
}

func (con *connection) Keys(pattern string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Keys(pattern)
	} else {
		return rg.Strings(con.c.Do("KEYS", pattern))
	}
}

// Migrate()

func (con *connection) Move(key, db string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Move(key, db)
	} else {
		res, err := rg.String(con.c.Do("MOVE", key, db))
		return getBool(res), err
	}
}

// Object()

func (con *connection) Persist(key string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Persist(key)
	} else {
		return rg.Bool(con.c.Do("PERSIST", key))
	}
}

func (con *connection) PExpire(key string, millisec int64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.PExpire(key, millisec)
	} else {
		return rg.Bool(con.c.Do("PEXPIRE", key, millisec))
	}
}

func (con *connection) PExpireat(key string, millisecTimestamp int64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.PExpireat(key, millisecTimestamp)
	} else {
		return rg.Bool(con.c.Do("PEXPIREAT", key, millisecTimestamp))
	}
}

func (con *connection) PTTL(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.PTTL(key)
	} else {
		return rg.Int(con.c.Do("PTTL", key))
	}
}

func (con *connection) RandomKey() (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.RandomKey()
	} else {
		return rg.String(con.c.Do("RANDOMKEY"))
	}
}

func (con *connection) Rename(key, newKey string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Rename(key, newKey)
	} else {
		res, err := rg.String(con.c.Do("Rename", key, newKey))
		return getBool(res), err
	}
}

func (con *connection) RenameNX(key, newKey string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.RenameNX(key, newKey)
	} else {
		res, err := rg.String(con.c.Do("RenameNX", key, newKey))
		return getBool(res), err
	}
}

func (con *connection) Restore(key string, ttl int, serializedValue string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Restore(key, ttl, serializedValue)
	} else {
		res, err := rg.String(con.c.Do("RESTORE", key, ttl, serializedValue))
		return getBool(res), err
	}
}

func (con *connection) RestoreWithReplace(key string, ttl int, serializedValue, replace string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.RestoreWithReplace(key, ttl, serializedValue, replace)
	} else {
		res, err := rg.String(con.c.Do("RESTORE", key, ttl, serializedValue, replace))
		return getBool(res), err
	}
}

// Scan()

func (con *connection) Sort(args ...interface{}) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Sort(args...)
	} else {
		return rg.Strings(con.c.Do("SORT", args...))
	}
}

func (con *connection) TTL(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.TTL(key)
	} else {
		return rg.Int(con.c.Do("TTL", key))
	}
}

func (con *connection) Type(key string) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Type(key)
	} else {
		return rg.String(con.c.Do("TYPE", key))
	}
}

func (con *connection) Wait(numSlaves, ttl int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Wait(numSlaves, ttl)
	} else {
		return rg.Int(con.c.Do("Wait", numSlaves, ttl))
	}
}

/*
	Server
*/
func (con *connection) FlushAll() (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.FlushAll()
	} else {
		res, err := rg.String(con.Do("FLUSHALL"))
		return getBool(res), err
	}
}

func (con *connection) FlushDB() (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.FlushDB()
	} else {
		res, err := rg.String(con.Do("FLUSHDB"))
		return getBool(res), err
	}
}

func (con *connection) Time() (map[string]int64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Time()
	} else {
		res, err := rg.Strings(con.Do("TIME"))

		if err != nil {
			return nil, err
		}

		result := make(map[string]int64)
		secs, _ := strconv.ParseInt(res[0], 10, 64)
		result["seconds"] = secs
		microSecs, _ := strconv.ParseInt(res[1], 10, 64)
		result["microseconds"] = microSecs

		return result, err
	}
}

/*
	Sorted Sets
*/
func (con *connection) ZAdd(key string, keyValue map[float64]string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZAdd(key, keyValue)
	} else {
		req := make([]interface{}, len(keyValue)*2+1)
		req[0] = key
		idx := 1
		for name, value := range keyValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		res, err := rg.Int(con.c.Do("ZADD", req...))
		return res, err
	}
}

func (con *connection) ZCard(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZCard(key)
	} else {
		return rg.Int(con.c.Do("ZCARD", key))
	}
}

func (con *connection) ZCount(key, min, max string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZCount(key, min, max)
	} else {
		return rg.Int(con.c.Do("ZCOUNT", key, min, max))
	}
}

func (con *connection) ZIncrBy(key string, increment float64, member string) (float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZIncrBy(key, increment, member)
	} else {
		return rg.Float64(con.c.Do("ZINCRBY", key, increment, member))
	}
}

// ZInterStore()

func (con *connection) ZLexCount(key, min, max string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZLexCount(key, min, max)
	} else {
		return rg.Int(con.c.Do("ZLEXCOUNT", key, min, max))
	}
}

func (con *connection) ZRange(key string, start, stop int) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRange(key, start, stop)
	} else {
		return rg.Strings(con.c.Do("ZRANGE", key, start, stop))
	}
}

func (con *connection) ZRangeWithScores(key string, start, stop int) (map[float64]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRangeWithScores(key, start, stop)
	} else {
		result := make(map[float64]string)
		res, err := rg.StringMap(con.c.Do("ZRANGE", key, start, stop, "WITHSCORES"))

		if err != nil {
			return nil, err
		}

		for key, value := range res {
			val, err := strconv.ParseFloat(value, 64)

			if err != nil {
				return nil, errors.New("redigowrap: Cannot convert string to float64")
			}

			result[val] = key
		}

		return result, nil
	}
}

// ZRangeByLex()
func (con *connection) ZRangeByScore(key, min, max string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRangeByScore(key, min, max)
	} else {
		return rg.Strings(con.Do("ZRANGEBYSCORE", key, min, max))
	}
}

func (con *connection) ZRangeByScoreWithScores(key, min, max string) (map[float64]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRangeByScoreWithScores(key, min, max)
	} else {
		result := make(map[float64]string)
		res, err := rg.StringMap(con.c.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))

		if err != nil {
			return nil, err
		}

		for key, value := range res {
			val, err := strconv.ParseFloat(value, 64)

			if err != nil {
				return nil, errors.New("redigowrap: Cannot convert string to float64")
			}

			result[val] = key
		}

		return result, nil
	}
}

func (con *connection) ZRank(key, member string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRank(key, member)
	} else {
		return rg.Int(con.Do("ZRANK", key, member))
	}
}

func (con *connection) ZRem(key string, members []string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRem(key, members)
	} else {
		req := make([]interface{}, len(members)+1)
		req[0] = key
		for idx, value := range members {
			req[idx+1] = value
		}

		return rg.Int(con.Do("ZREM", req...))
	}
}

// ZRemRangeByLex()
func (con *connection) ZRemRangeByRank(key string, start, stop int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRemRangeByRank(key, start, stop)
	} else {
		return rg.Int(con.Do("ZREMRANGEBYRANK", key, start, stop))
	}
}

func (con *connection) ZRemRangeByScore(key, min, max string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRemRangeByScore(key, min, max)
	} else {
		return rg.Int(con.Do("ZREMRANGEBYSCORE", key, min, max))
	}
}

func (con *connection) ZRevRange(key string, start, stop int) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRevRange(key, start, stop)
	} else {
		return rg.Strings(con.c.Do("ZREVRANGE", key, start, stop))
	}
}

// ZRevRangeByLex()
func (con *connection) ZRevRangeByScore(key, min, max string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRevRangeByScore(key, min, max)
	} else {
		return rg.Strings(con.Do("ZREVRANGEBYSCORE", key, min, max))
	}
}

func (con *connection) ZRevRangeByScoreWithScores(key, min, max string) (map[float64]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRevRangeByScoreWithScores(key, min, max)
	} else {
		result := make(map[float64]string)
		res, err := rg.StringMap(con.c.Do("ZREVRANGEBYSCORE", key, min, max, "WITHSCORES"))

		if err != nil {
			return nil, err
		}

		for key, value := range res {
			val, err := strconv.ParseFloat(value, 64)

			if err != nil {
				return nil, errors.New("redigowrap: Cannot convert string to float64")
			}

			result[val] = key
		}

		return result, nil
	}
}

func (con *connection) ZRevRank(key, member string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZRevRank(key, member)
	} else {
		return rg.Int(con.Do("ZREVRANK", key, member))
	}
}

// ZScan()
func (con *connection) ZScore(key, member string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.ZScore(key, member)
	} else {
		return rg.Int(con.Do("ZSCORE", key, member))
	}
}

// ZUnionStore()

/*
	Strings
*/
func (con *connection) Append(key, value string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Append(key, value)
	} else {
		return rg.Int(con.c.Do("APPEND", key, value))
	}
}

func (con *connection) BitCount(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.BitCount(key)
	} else {
		return rg.Int(con.c.Do("BITCOUNT", key))
	}
}

func (con *connection) BitCountRange(key string, start, end int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.BitCountRange(key, start, end)
	} else {
		return rg.Int(con.c.Do("BITCOUNTRANGE", key, start, end))
	}
}

func (con *connection) BitOP(key, destKey string, keys []interface{}) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.BitOP(key, destKey, keys)
	} else {
		return rg.Int(con.c.Do("BITOP", keys...))
	}
}

func (con *connection) BitPos(key string, start int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.BitPos(key, start)
	} else {
		return rg.Int(con.c.Do("BITPOS", key, start))
	}
}

func (con *connection) BitPosRange(key string, start, end int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.BitPosRange(key, start, end)
	} else {
		return rg.Int(con.c.Do("BITPOS", key, start, end))
	}
}

func (con *connection) Decr(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Decr(key)
	} else {
		return rg.Int(con.c.Do("DECR", key))
	}
}

func (con *connection) DecrBy(key string, decrement int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.DecrBy(key, decrement)
	} else {
		return rg.Int(con.c.Do("DECRBY", key, decrement))
	}
}

func (con *connection) Get(key string) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Get(key)
	} else {
		return rg.String(con.c.Do("GET", key))
	}
}

func (con *connection) GetFloat64(key string) (float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.GetFloat64(key)
	} else {
		return rg.Float64(con.c.Do("GET", key))
	}
}

func (con *connection) GetBit(key string, offset int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.GetBit(key, offset)
	} else {
		return rg.Int(con.c.Do("GETBIT", key, offset))
	}
}

func (con *connection) GetRange(key string, start, end int) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.GetRange(key, start, end)
	} else {
		return rg.String(con.c.Do("GETRANGE", key, start, end))
	}
}

func (con *connection) GetSet(key, value string) (string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.GetSet(key, value)
	} else {
		return rg.String(con.c.Do("GETSET", key, value))
	}
}

func (con *connection) Incr(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Incr(key)
	} else {
		return rg.Int(con.c.Do("INCR", key))
	}
}

func (con *connection) IncrBy(key string, increment int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.IncrBy(key, increment)
	} else {
		return rg.Int(con.c.Do("INCRBY", key, increment))
	}
}

func (con *connection) IncrByFloat(key string, increment float64) (float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.IncrByFloat(key, increment)
	} else {
		return rg.Float64(con.c.Do("INCRBYFLOAT", key, increment))
	}
}

func (con *connection) MGet(keys []string) ([]string, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.MGet(keys)
	} else {
		req := make([]interface{}, len(keys))
		for idx, val := range keys {
			req[idx] = val
		}
		return rg.Strings(con.c.Do("MGET", req...))
	}
}

func (con *connection) MGetFloat64(keys []string) ([]float64, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.MGetFloat64(keys)
	} else {
		req := make([]interface{}, len(keys))
		for idx, val := range keys {
			req[idx] = val
		}

		result, err := rg.Strings(con.c.Do("MGET", req...))

		if err != nil {
			return nil, err
		}

		res := make([]float64, len(result))

		for idx, value := range result {
			res[idx] = strToFloat64(value)
		}

		return res, nil
	}
}

func (con *connection) MSet(keyValue map[string]string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.MSet(keyValue)
	} else {
		req := make([]interface{}, len(keyValue)*2)
		idx := 0
		for name, value := range keyValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		res, err := rg.String(con.c.Do("MSET", req...))
		return getBool(res), err
	}
}

func (con *connection) MSetFloat64(keyValue map[string]float64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.MSetFloat64(keyValue)
	} else {
		req := make([]interface{}, len(keyValue)*2)
		idx := 0
		for name, value := range keyValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		res, err := rg.String(con.c.Do("MSET", req...))
		return getBool(res), err
	}
}

func (con *connection) MSetNX(keyValue map[string]string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.MSetNX(keyValue)
	} else {
		req := make([]interface{}, len(keyValue)*2)
		idx := 0
		for name, value := range keyValue {
			req[idx] = name
			idx++
			req[idx] = value
			idx++
		}

		return rg.Int(con.c.Do("MSETNX", req...))
	}
}

func (con *connection) PSetEX(key, value string, millisec int64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer con.Release()
		return c.PSetEX(key, value, millisec)
	} else {
		res, err := rg.String(con.c.Do("PSETEX", key, millisec, value))
		return getBool(res), err
	}
}

func (con *connection) Set(key, value string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Set(key, value)
	} else {
		res, err := rg.String(con.c.Do("SET", key, value))
		return getBool(res), err
	}
}

func (con *connection) SetFloat64(key string, value float64) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.SetFloat64(key, value)
	} else {
		res, err := rg.String(con.c.Do("SET", key, value))
		return getBool(res), err
	}
}

func (con *connection) SetBit(key, value string, offset int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.SetBit(key, value, offset)
	} else {
		return rg.Int(con.c.Do("SETBIT", key, value, offset))
	}
}

func (con *connection) SetEX(key, value string, seconds int) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.SetEX(key, value, seconds)
	} else {
		res, err := rg.String(con.c.Do("SETEX", key, seconds, value))
		return getBool(res), err
	}
}

func (con *connection) SetNX(key, value string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.SetNX(key, value)
	} else {
		return rg.Int(con.c.Do("SETNX", key, value))
	}
}

func (con *connection) SetRange(key, value string, offset int) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.SetRange(key, value, offset)
	} else {
		return rg.Int(con.c.Do("SETRANGE", key, offset, value))
	}
}

func (con *connection) StrLen(key string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.StrLen(key)
	} else {
		return rg.Int(con.c.Do("STRLEN", key))
	}
}

/*
   Transactions
*/
func (con *connection) Discard() (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Discard()
	} else {
		return rg.Bool(con.c.Do("DISCARD"))
	}
}

func (con *connection) Exec() ([]interface{}, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Exec()
	} else {
		return rg.Values(con.c.Do("EXEC"))
	}
}

func (con *connection) Multi() (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Multi()
	} else {
		return rg.Bool(con.c.Do("MULTI"))
	}
}

func (con *connection) Unwatch() (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Unwatch()
	} else {
		return rg.Bool(con.c.Do("UNWATCH"))
	}
}

func (con *connection) Watch(keys []string) (bool, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.Watch(keys)
	} else {
		req := make([]interface{}, len(keys))
		for idx, val := range keys {
			req[idx] = val
		}
		res, err := rg.String(con.c.Do("WATCH", req...))
		return getBool(res), err
	}
}

func (con *connection) RPUSH(key, data string) (int, error) {
	if con.p != nil {
		c, _ := con.GetConnection()
		defer c.Release()
		return c.RPUSH(key, data)
	} else {
		return rg.Int(con.c.Do("RPUSH", key, data))
	}
}
