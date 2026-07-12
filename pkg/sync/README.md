# Sync

Thread-safe generic map implementation and pagination helpers.

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

## Paginate

`Paginate` returns an iterator (`iter.Seq2[T, error]`) over items fetched page by page with limit/offset pagination. It calls `fetch` with increasing offset until a page shorter than `limit` is returned.

```go
func Paginate[T any](ctx context.Context, limit int,
    fetch func(ctx context.Context, limit, offset int) ([]T, error),
) iter.Seq2[T, error]
```

### Usage

```go
import sdksync "github.com/dipdup-net/indexer-sdk/pkg/sync"

for item, err := range sdksync.Paginate(ctx, 100, storage.List) {
    if err != nil {
        return err
    }
    process(item)
}
```

### Behavior

- Items are yielded one by one; the next page is fetched only after the current one is fully consumed.
- On `fetch` failure the error is yielded as the last pair (with zero `T`) and the sequence ends — check the error on every iteration.
- If the total count is a multiple of `limit`, the final `fetch` call returns an empty page.
- If `limit` is not positive, the sequence is empty and `fetch` is never called.
- Breaking out of the loop stops the iterator without fetching further pages.
