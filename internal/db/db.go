package db

import (
	"github.com/pocketbase/pocketbase"
	"github.com/uforg/ufogateway/internal/cache"
)

type DB struct {
	app           *pocketbase.PocketBase
	cacheInstance *cache.CacheInstance
}

func NewDB(
	app *pocketbase.PocketBase,
	cacheInstance *cache.CacheInstance,
) *DB {
	return &DB{
		app:           app,
		cacheInstance: cacheInstance,
	}
}
