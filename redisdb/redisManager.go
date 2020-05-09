package redisdb

import (
	"sync"

	"gopkg.in/redis.v5"
)

// redisConnURLMap connection URL map
type redisConnURLMap struct {
	mutex         *sync.RWMutex
	connectionMap map[string]*redis.Client
}

var (
	connectionList *redisConnURLMap
	cOnce          sync.Once
)

func init() {
	connectionList = newConectionList()
}

func newConectionList() *redisConnURLMap {

	cOnce.Do(func() {
		connectionList = &redisConnURLMap{
			mutex:         new(sync.RWMutex),
			connectionMap: make(map[string]*redis.Client),
		}
	})
	return connectionList
}

func connectRedis(address, password string, db int) *redis.Client {
	return redis.NewClient(loadRedisOptions(address, password, db))
}

func loadRedisOptions(address, password string, db int) *redis.Options {

	return &redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	}
}

func (c *redisConnURLMap) setConnection(restaurantID string, redisConnection *redis.Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.connectionMap[restaurantID] = redisConnection
}

func (c *redisConnURLMap) getConnection(restaurantID string) (*redis.Client, error) {
	c.mutex.RLock()
	redisConnection, cached := c.connectionMap[restaurantID] // check if connection is cached
	c.mutex.RUnlock()

	if !cached {

		addr := "" //hostName + ":" + port
		password := ""
		db := 0

		redisClient := connectRedis(addr, password, db)
		c.setConnection(restaurantID, redisClient)
		return redisClient, nil
	}
	return redisConnection, nil
}

// Get gets the value for given key
func Get(restaurantID string, key string) (string, error) {
	redisClient, connectionError := connectionList.getConnection(restaurantID)

	if connectionError != nil {
		return "", connectionError
	}
	r := redisClient.Get(key)
	d, e := r.Result()
	return d, e
}
