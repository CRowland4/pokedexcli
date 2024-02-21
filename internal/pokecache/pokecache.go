package pokecache

import (
	"time"
	"sync"
	"slices"
)
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

const cacheEntryLifeSpan = time.Duration(1 * time.Minute)

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type Cache struct{
	Info map[int]locationEntry
	mu *sync.Mutex
	Pokemon map[string]PokemonData
}

type locationEntry struct{
	createdAt time.Time
	LocationName string
	LocationPokemon []string
}

type PokemonData struct{
	IsCaught bool
	BaseExperience int
	Height int
	Weight int
	HP int
	Attack int
	Defense int
	SpecialAttack int
	SpecialDefense int
	Speed int
	Types []string
}

/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

func (c *Cache) AddLocation(id int, areaName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	newAreaEntry := locationEntry{
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
	if !slices.Contains(c.Info[locationID].LocationPokemon, pokemonName) {
		location := c.Info[locationID]
		location.LocationPokemon = append(location.LocationPokemon, pokemonName)
		c.Info[locationID] = location
	}
	return
}

func (c *Cache) AddPokemon(name string, data PokemonData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Pokemon[name] = data
	return
}

func (c *Cache) GetLocation(id int) (entry locationEntry, isFound bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, isFound = c.Info[id]
	if isFound {
		return entry, true
	}

	return entry, false
}

func (c *Cache) GetCaughtPokemon() (caughtPokemon []string) {
	for name, data := range c.Pokemon {
		if data.IsCaught {
			caughtPokemon = append(caughtPokemon, name)
		}
	}

	return caughtPokemon
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
		Info: make(map[int]locationEntry),
		mu: new(sync.Mutex),
		Pokemon: make(map[string]PokemonData),
	}
	go pokeCache.reapLoop(interval)
	return pokeCache
}