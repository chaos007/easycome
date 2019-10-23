package redis

import (
	rg "github.com/garyburd/redigo/redis"
)

type (
	connection struct {
		c rg.Conn
		p *rg.Pool
	}

	RedisCommands interface {
		// Redigo Function
		Do(cmd string, args ...interface{}) (interface{}, error)
		Send(cmd string, args ...interface{}) error
		Flush() error
		Receive() (interface{}, error)

		// Cluster

		// Connection
		Auth(password string) (bool, error)
		Echo(message string) (string, error)
		Ping() (string, error)
		Select(index int) (bool, error)
		Quit() (string, error)

		// Hashes
		HDel(hashKey string, fields []string) (int, error)
		HExists(hashKey, field string) (bool, error)
		HGet(hashKey, field string) (string, error)
		HGetFloat64(hashKey, field string) (float64, error)
		HGetAll(hashKey string) (map[string]string, error)
		HGetAllInterface(hashKey string) (map[string]interface{}, error)
		HGetAllFloat64(hashKey string) (map[string]float64, error)
		HIncrBy(hashKey, field string, increment int) (int, error)
		HIncrByFloat(hashKey, field string, increment float64) (float64, error)
		HKeys(hashKey string) ([]string, error)
		HLen(hashKey string) (int, error)
		HMGet(hashKey string, fields []string) ([]string, error)
		HMGetFloat64(hashKey string, fields []string) ([]float64, error)
		HMSet(hashKey string, fieldValue map[string]string) (bool, error)
		HMSetInterface(hashKey string, fieldValue map[string]interface{}) (bool, error)
		HMSetFloat64(hashKey string, fieldValue map[string]float64) (bool, error)
		// HScan() ()
		HSet(hashKey, field, value string) (int, error)
		HSetFloat64(hashKey, field string, value float64) (int, error)
		HSetNX(hashKey, field, value string) (int, error)
		HStrLen(hashKey, field string) (int, error)
		HVals(hashKey string) ([]string, error)

		// HyperLogLog

		// Keys
		Del(keys []string) (int, error)
		Dump(key string) (string, error)
		Exists(key string) (bool, error)
		Expire(key string, seconds int) (bool, error)
		Expireat(key string, timestamp int64) (bool, error)
		Keys(pattern string) ([]string, error)
		// Migrate(host, port, key, destDB, timeout) (int)
		Move(key, db string) (bool, error)
		// Object(subcommand string, args []string) (error)
		Persist(key string) (bool, error)
		PExpire(key string, millisec int64) (bool, error)
		PExpireat(key string, millisecTimestamp int64) (bool, error)
		PTTL(key string) (int, error)
		RandomKey() (string, error)
		Rename(key, newKey string) (bool, error)
		RenameNX(key, newKey string) (bool, error)
		Restore(key string, ttl int, serializedValue string) (bool, error)
		RestoreWithReplace(key string, ttl int, serializedValue string, replace string) (bool, error)
		// Scan(scanValue string) (map[string][]string, error)
		Sort(args ...interface{}) ([]string, error)
		TTL(key string) (int, error)
		Type(key string) (string, error)
		Wait(numSlaves, timeout int) (int, error)

		// Lists
		RPUSH(key, data string) (int, error)
		// Pub/Sub

		// Scripting

		// Server
		FlushAll() (bool, error)
		FlushDB() (bool, error)
		Time() (map[string]int64, error)

		// Sets

		// Sorted Sets
		ZAdd(key string, keyValue map[float64]string) (int, error)
		ZCard(key string) (int, error)
		ZCount(key, min, max string) (int, error)
		ZIncrBy(key string, increment float64, member string) (float64, error)
		// ZInterStore()
		ZLexCount(key, min, max string) (int, error)
		ZRange(key string, start, stop int) ([]string, error)
		ZRangeWithScores(key string, start, stop int) (map[float64]string, error)
		// ZRangeByLex()
		ZRangeByScore(key, min, max string) ([]string, error)
		ZRangeByScoreWithScores(key, min, max string) (map[float64]string, error)
		ZRank(key, member string) (int, error)
		ZRem(key string, members []string) (int, error)
		// ZRemRangeByLex()
		ZRemRangeByRank(key string, start, stop int) (int, error)
		ZRemRangeByScore(key, min, max string) (int, error)
		ZRevRange(key string, start, stop int) ([]string, error)
		// ZRevRangeByLex()
		ZRevRangeByScore(key, min, max string) ([]string, error)
		ZRevRangeByScoreWithScores(key, min, max string) (map[float64]string, error)
		ZRevRank(key, member string) (int, error)
		// ZScan()
		ZScore(key, member string) (int, error)
		// ZUnionStore()

		// Strings
		Append(key, value string) (int, error)
		BitCount(key string) (int, error)
		BitCountRange(key string, start, end int) (int, error)
		BitOP(operation, destKey string, keys []interface{}) (int, error)
		BitPos(key string, start int) (int, error)
		BitPosRange(key string, start, end int) (int, error)
		Decr(key string) (int, error)
		DecrBy(key string, decrement int) (int, error)
		Get(key string) (string, error)
		GetFloat64(key string) (float64, error)
		GetBit(key string, offset int) (int, error)
		GetRange(key string, start, end int) (string, error)
		GetSet(key, value string) (string, error)
		Incr(key string) (int, error)
		IncrBy(key string, increment int) (int, error)
		IncrByFloat(key string, increment float64) (float64, error)
		MGet(keys []string) ([]string, error)
		MGetFloat64(keys []string) ([]float64, error)
		MSet(keyValue map[string]string) (bool, error)
		MSetFloat64(keyValue map[string]float64) (bool, error)
		MSetNX(keyValue map[string]string) (int, error)
		PSetEX(key, value string, millisec int64) (bool, error)
		Set(key, value string) (bool, error)
		SetFloat64(key string, value float64) (bool, error)
		SetBit(key, value string, offset int) (int, error)
		SetEX(key, value string, seconds int) (bool, error)
		SetNX(key, value string) (int, error)
		SetRange(key, value string, offset int) (int, error)
		StrLen(key string) (int, error)

		// Transactions
		Discard() (bool, error)
		Exec() ([]interface{}, error)
		Multi() (bool, error)
		Unwatch() (bool, error)
		Watch(keys []string) (bool, error)
	}

	ConnectionCommands interface {
		Close() bool
	}

	PoolCommands interface {
		GetConnection() (PoolConnection, error)
		ActiveCount() int
		Release()
		PoolClose()
	}
)
