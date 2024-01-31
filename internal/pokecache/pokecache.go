package pokecache

import (
	"time"
	"sync"
)

var (
	pokeCache map[int]cacheEntry
	mu sync.Mutex
)
const cacheEntryLifeSpan = time.Duration(1 * time.Minute)

type cacheEntry struct{
	createdAt time.Time
	areaName string  // TODO have the cache entry be the full JSON response, not just the name
}

func InitializePokeCache() {
	pokeCache = make(map[int]cacheEntry)
	go periodicallyRemoveOldPokeCacheEntries()
	return
}

func IsCached(id int) bool {
	_, ok := Get(id)
	return ok
}

func Add(id int, areaName string) {
	mu.Lock()
	defer mu.Unlock()
	
	newAreaEntry := cacheEntry{createdAt: time.Now(), areaName: areaName}
	pokeCache[id] = newAreaEntry
	return
}

func Get(id int) (locationName string, isFound bool) {
	entry, isFound := pokeCache[id]
	if isFound {
		return entry.areaName, true
	}

	return "", false
}

/*=========================================================================================================================*/

func periodicallyRemoveOldPokeCacheEntries() {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		currentTime := <- ticker.C
		clearOldPokeCacheEntries(currentTime)
	}
}

func clearOldPokeCacheEntries(timeValue time.Time) {
	for id, entry := range pokeCache {
		entryAge := timeValue.Sub(entry.createdAt)
		if entryAge > cacheEntryLifeSpan {
			delete(pokeCache, id)
		}
	}

	return
}
