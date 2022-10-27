package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

// Add tags to the struct elements for Redis' scan struct function.
// Not disimilar to the way you unmarshal a JSON response into a struct,
// This is essentially unmarshalling a Redis reply into a struct.
type Podcast struct {
	Title    string  `redis:"title"`
	Creator  string  `redis:"creator"`
	Category string  `redis:"category"`
	Fee      float64 `redis:"membership_fee"`
}

// checkError checks if there is an error and logs accordingly.
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Connect to the Redis instance. The default port is 6379.
	// Side note: The only issue with redis.Dial is that the connection object you get back is not
	// 						safe for concurrent use. So if you have multiple Go routines that need to access
	// 						the redis instance, you need to use connection pooling to get its own connection
	// 						whenever the process needs it. Once the process does what it needs with the connection,
	// 						the connection needs to be closed and returned to the pool for other Go routines to use.
	connect, err := redis.Dial("tcp", "localhost:6379")
	checkError(err)
	defer connect.Close()

	// Set the key.
	_, err = connect.Do(
		"HMSET",     // Sets multiple hash fields on a key.
		"podcast:1", // Cached instance of a type. "podcast:1" is the key, the rest is technically the value.
		"title",
		"The WAN Show",
		"creator",
		"Linus Tech Tips",
		"category",
		"technology",
		"membership_fee",
		9.99,
	)
	checkError(err)
	// Get the map from the Redis instance, get the podcast field, convert
	// it to a string, and then print it out to the console.
	title, err := redis.String(connect.Do("HGET", "podcast:1", "title"))
	checkError(err)
	// Get the membership fee from the map and convert it to a float64 to be used.
	fee, err := redis.Float64(connect.Do("HGET", "podcast:1", "membership_fee"))
	checkError(err)

	// Get all elements from the map.
	values, err := redis.StringMap(connect.Do("HGETALL", "podcast:1"))
	checkError(err)

	// Print individual fields from the map.
	fmt.Println("Podcast Title:", title)
	fmt.Printf("Podcast Membership Fee: $%v\n", fee)

	// Print all fields from the map via a for-each loop.
	for k, v := range values {
		fmt.Println("Key:", k)
		fmt.Println("Value:", v)
	}

	// This time, the reply is needeed to run the ScanStruct function.
	reply, err := redis.Values(connect.Do("HGETALL", "podcast:1"))
	checkError(err)

	// Create a variable to store the struct.
	var podcast Podcast

	// Use ScanStruct to create a slice with all the information from the map.
	err = redis.ScanStruct(reply, &podcast)
	checkError(err)

	// Print all the podcast details from the map via the Podcast struct.
	fmt.Printf("Podcast: %+v\n", podcast)
}
