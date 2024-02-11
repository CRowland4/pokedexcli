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
	PokemonBaseExperience map[string]float64
}

type cacheEntry struct{
	createdAt time.Time
	LocationName string
	LocationPokemon []string
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

func (c *Cache) AddLocation(id int, areaName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	newAreaEntry := cacheEntry{
		createdAt: time.Now(),
		LocationName: areaName,
		LocationPokemon: []string{},
	}

	c.Info[id] = newAreaEntry
	return
}

func (c *Cache) AddPokemonToLocation(locationID int, pokemonName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	location := c.Info[locationID]
	location.LocationPokemon = append(location.LocationPokemon, pokemonName)
	c.Info[locationID] = location
	return
}

func (c *Cache) AddPokemon(name string, base_experience float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.PokemonBaseExperience[name] = base_experience
	return
}

func (c *Cache) GetLocation(id int) (entry cacheEntry, isFound bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, isFound = c.Info[id]
	if isFound {
		return entry, true
	}

	return entry, false
}

func (c *Cache) GetPokemonBaseExperience(name string) (experience float64, isFound bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	experience, isFound = c.PokemonBaseExperience[name]
	return experience, isFound
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
		PokemonBaseExperience: make(map[string]float64),
	}
	go pokeCache.reapLoop(interval)
	return pokeCache
}