package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

var RDB *redis.Client

func GetRedisKeyFromInputs(method string, params []string) string {
	// method:eth_blockNumber params:[]string{} -> method_eth_blockNumber_params_[]
	return fmt.Sprintf("method_%s_params_%v", method, params)
}

func RedisSetKey(key string, value map[string]interface{}, ttl time.Duration) error {
	// ttl for defining how long the key,value will be kept in the database, 0 means forever
	// ttl for defining how long the key,value will be kept in the database, 0 means forever
	bytes, _ := json.Marshal(value) // map value converted to bytes
	err := RDB.Set(ctx, key, bytes, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func RedisGetValue(key string) (map[string]interface{}, error) {
	// Redis key,value pairs, find the corresponding value of a key
	val, err := RDB.Get(ctx, key).Result()
	if err != nil {
		// Not found? return err
		return nil, err
	}
	// Un cast bytes data to map
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(val), &data)

	return data, nil
}
