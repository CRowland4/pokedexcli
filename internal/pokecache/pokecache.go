package pokecache

import (
	"time"
	"sync"
)
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

const cacheEntryLifeSpan = time.Duration(1 * time.Minute)

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type Cache struct{
	Info map[int]cacheEntry
	mu *sync.Mutex
}

type cacheEntry struct{
	createdAt time.Time
	LocationName string
	Pokemon []string
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type Cacher interface{
	Add(key int, areaName string)
	Get(key int) (areaName string, isFound bool)
	reapLoop(interval time.Duration)
}

func (c *Cache) Add(id int, areaName string, pokemon []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	newAreaEntry := cacheEntry{
		createdAt: time.Now(),
		LocationName: areaName,
		Pokemon: pokemon,
	}

	c.Info[id] = newAreaEntry
	return
}

func (c *Cache) Get(id int) (entry cacheEntry, isFound bool) {
	entry, isFound = c.Info[id]
	if isFound {
		return entry, true
	}

	return entry, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		currentTime := <- ticker.C
		for id, entry := range (*c).Info {
			entryAge := currentTime.Sub(entry.createdAt)
			if entryAge > cacheEntryLifeSpan {
				delete((*c).Info, id)
			}
		}
	}
}

/*==================================================================================================================================*/


func NewCache(interval time.Duration) (pokeCache Cache) {
	pokeCache = Cache{
		Info: make(map[int]cacheEntry),
		mu: new(sync.Mutex),
	}
	go pokeCache.reapLoop(interval)
	return pokeCache
}