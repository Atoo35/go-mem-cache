# Go Mem Cache

This repository is a Go module for creating a memory cache with a specified capacity and key eviction policy.

## Overview

`db.go` provides a set of functions to perform operations on a memory cache. It uses the `container/heap` and `sync` packages from the Go standard library.

## Types

- `DB`: The main type representing the memory cache.
- `KeyEvictionPolicy`: A type used for setting the eviction policy.
- `Option`: A type for optional configuration functions.

## Functions

- `New()`: This function creates a new memory cache.
- `WithKeyEvictionPolicy()`: This function sets the key eviction policy for the memory cache.
- `WithCapacity()`: This function sets the capacity for the memory cache.
- `GetCapacity()`: This function gets the capacity of the memory cache.
- `GetKeyEvictionPolicy()`: This function gets the key eviction policy of the memory cache.
- `GetKeys()`: This function gets all the keys in the memory cache.
- `Get()`: This function gets an item from the memory cache.
- `SetWithExpiry()`: This function sets an item in the memory cache with a specified expiry.
- `Set()`: This function sets an item in the memory cache.
- `Size()`: This function gets the size of the memory cache.

## Usage

To use this module, import it in your Go code:

```go
import "github.com/Atoo35/go-mem-cache"
```

Then, you can call its functions like this:

```go
    // Configurable options.
    db := gomemcache.New()

    db.Set("key", "value", 1)
    value, err := db.Get("key")
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
