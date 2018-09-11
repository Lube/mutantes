package testUtils

import (
	"github.com/go-redis/redis"
	"github.com/lube/mutantes/app"
)

// GetDB Connects to a Redis DB for testing purposes
func GetDB() *redis.Client {
	err := app.LoadConfig("./config", "../config")
	if err != nil {
		panic(err)
	}

	db := redis.NewClient(&redis.Options{
		Addr:     app.Config.DSN,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = db.Ping().Result()
	if err != nil {
		panic(err)
	}

	return db
}

// ResetDB deletes all the elements in the set.
// This method is mainly used in tests.
func ResetDB() *redis.Client {
	err := app.LoadConfig("./config", "../config")
	if err != nil {
		panic(err)
	}

	db := redis.NewClient(&redis.Options{
		Addr:     app.Config.DSN,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = db.Ping().Result()
	if err != nil {
		panic(err)
	}

	err = db.Del("mutantes").Err()
	if err != nil {
		panic(err)
	}
	err = db.Del("humanos").Err()
	if err != nil {
		panic(err)
	}
	return db
}
