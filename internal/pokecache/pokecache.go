package pokecache

var PokeCache map[int]string  // Maps location IDs to location names

func InitializePokeCache() {
	PokeCache = make(map[int]string)
	return
}