# Sync

Thread-safe generic map implementation.

## Map[K, V]

`Map` is a generic concurrent-safe map protected by `sync.RWMutex`.

```go
import sdksync "github.com/dipdup-net/indexer-sdk/pkg/sync"

m := sdksync.NewMap[string, int]()
```

### Methods

```go
m.Set("key", 42)              // set value
val, ok := m.Get("key")       // get value
m.Delete("key")               // delete key
m.Clear()                     // remove all entries
length := m.Len()             // number of entries

keys := m.Keys()              // []K — all keys
values := m.Values()          // []V — all values

// Iterate over entries
err := m.Range(func(key string, value int) (error, bool) {
    fmt.Println(key, value)
    return nil, false // return (nil, true) to stop iteration
})
```

### Thread Safety

All operations acquire the appropriate lock:
- `Get`, `Len`, `Keys`, `Values`, `Range` — read lock (`RLock`)
- `Set`, `Delete`, `Clear` — write lock (`Lock`)
