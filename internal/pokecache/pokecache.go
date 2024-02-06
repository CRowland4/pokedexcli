package pokecache

import (
	"time"
	"sync"
)
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

const cacheEntryLifeSpan = time.Duration(1 * time.Minute)

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type Cache struct{
	info map[int]cacheEntry
	mu *sync.Mutex
}

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type Cacher interface{
	Add(key int, areaName string)
	Get(key int) (areaName string, isFound bool)
	reapLoop(interval time.Duration)
}

func (c *Cache) Add(id int, areaName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	newAreaEntry := cacheEntry{
		createdAt: time.Now(),
		val: []byte(areaName),
	}

	c.info[id] = newAreaEntry
	return
}

func (c *Cache) Get(id int) (locationName string, isFound bool) {
	entry, isFound := c.info[id]
	if isFound {
		return string(entry.val), true
	}

	return "", false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		currentTime := <- ticker.C
		for id, entry := range (*c).info {
			entryAge := currentTime.Sub(entry.createdAt)
			if entryAge > cacheEntryLifeSpan {
				delete((*c).info, id)
			}
		}
	}
}

/*==================================================================================================================================*/


func NewCache(interval time.Duration) (pokeCache Cache) {
	pokeCache = Cache{
		info: make(map[int]cacheEntry),
		mu: new(sync.Mutex),
	}
	go pokeCache.reapLoop(interval)
	return pokeCache
}