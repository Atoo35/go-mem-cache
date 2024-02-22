package main

import (
	"fmt"
	"time"

	gomemcache "github.com/Atoo35/go-mem-cache"
	types "github.com/Atoo35/go-mem-cache/src/types"
)

func main() {
	db := gomemcache.New[int64](
		gomemcache.WithKeyEvictionPolicy[int64](
			types.KeyEvictionPolicy("LRU"),
		),
		gomemcache.WithCapacity[int64](5),
	)

	// db.SetWithExpiry("A", 5, 5, time.Now().Unix()+5)
	// db.SetWithExpiry("B", 4, 1, time.Now().Unix()+500)
	// db.SetWithExpiry("C", 3, 5, time.Now().Unix()+1000000)
	// db.SetWithExpiry("D", 2, 9, time.Now().Unix()+10000000)
	// db.SetWithExpiry("E", 1, 5, time.Now().Unix()+100000000)
	db.Set("A", 5, 5)
	fmt.Println(db.GetKeys())
	db.Set("B", 4, 1)
	fmt.Println(db.GetKeys())
	db.Set("C", 3, 5)
	fmt.Println(db.GetKeys())
	db.Set("D", 2, 9)
	fmt.Println(db.GetKeys())
	db.Set("E", 1, 5)
	fmt.Println(db.GetKeys())

	// Get the value associated with key "C"
	value, err := db.Get("C")
	if err == nil {
		fmt.Println("Value for key 'C':", value)
	} else {
		fmt.Println("Error:", err)
	}

	// Simulate time passing to check expiration and update values
	time.Sleep(time.Second * 3)

	// Set new values in the db
	db.Set("F", 10, 5)
	fmt.Println(db.GetKeys())
	db.Set("G", 9, 5)
	fmt.Println(db.GetKeys())
	db.Set("H", -1, 5)
	fmt.Println(db.GetKeys())
	db.Set("I", 1, 5)
	fmt.Println(db.GetKeys())

	// Update value for key "C"
	db.Set("C", 1, 5)
	fmt.Println(db.GetKeys())
	// Attempt to get the value associated with key "D" (should return an error)
	value, err = db.Get("D")
	if err == nil {
		fmt.Println("Value for key 'D':", value)
	} else {
		fmt.Println("Error:", err)
	}
}
