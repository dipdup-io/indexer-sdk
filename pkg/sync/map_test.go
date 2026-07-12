package sync

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap_Get(t *testing.T) {
	m := NewMap[int, string]()
	m.Set(10, "hello sdk sync map based on RWMutex")

	value, ok := m.Get(10)
	if !ok {
		t.Fatal("existing key was not found")
	}

	if value != "hello sdk sync map based on RWMutex" {
		t.Fatal("found value is incorrect")
	}
}

func TestMap_Delete(t *testing.T) {
	m := NewMap[int, string]()
	m.Set(10, "hello sdk sync map based on RWMutex")
	m.Set(11, "this value and key will be deleted")

	m.Delete(11)

	_, ok := m.Get(11)
	if ok {
		t.Fatal("non-existing key was found")
	}
}

func TestMap_Range(t *testing.T) {
	m := NewMap[int, string]()
	m.Set(10, "hello sdk sync map based on RWMutex")
	m.Set(11, "second value")

	checkData := map[int]*struct {
		checked bool
		value   string
	}{
		10: {value: "hello sdk sync map based on RWMutex"},
		11: {value: "second value"},
	}

	handler := func(k int, v string) (error, bool) {
		toCheck, ok := checkData[k]
		if !ok {
			t.Fatal("found non-existing key")
			return nil, true
		}

		if v != toCheck.value {
			t.Fatalf("found value is incorrect for key=%d with value=%s, looking for value=%s", k, v, toCheck.value)
			return nil, true
		}

		toCheck.checked = true
		return nil, false
	}

	if err := m.Range(handler); err != nil {
		t.Fatalf("error occured in Range %+v", err)
	}

	for k, v := range checkData {
		if !v.checked {
			t.Fatalf("key %d was not applied in Range", k)
		}
	}
}

func TestMap_ConcurrentRange(t *testing.T) {
	const mapSize = 1 << 10

	m := NewMap[int64, int64]()
	for n := int64(1); n <= mapSize; n++ {
		m.Set(n, n)
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()

	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if r.Int63n(mapSize) == 0 {
						m.Set(n, n*i*g)
					} else {
						m.Get(n)
					}
				}
			}
		}(g)
	}

	for n := 16; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		err := m.Range(func(k, v int64) (error, bool) {
			if v%k != 0 {
				t.Fatalf("while Setting multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return nil, false
		})

		if len(seen) != mapSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}

		if err != nil {
			t.Fatalf("error occured in Range %+v", err)
		}
	}
}

func TestMap_Clear(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"clear", "sync", "map"} {
		m.Set(i, v)
	}

	m.Clear()

	length := 0
	err := m.Range(func(key int, value string) (error, bool) {
		length++
		return nil, false
	})

	if err != nil {
		t.Fatalf("error occured in checking length of Range %+v", err)
	}

	if length != 0 {
		t.Fatalf("unexpected map size, got %v want %v", length, 0)
	}
}

func TestMap_Len(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"len", "sync", "map"} {
		m.Set(i, v)
	}

	length := m.Len()

	if length != 3 {
		t.Fatalf("unexpected map size, got %v want %v", length, 3)
	}
}

func TestMap_Values(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"len", "sync", "map"} {
		m.Set(i, v)
	}

	values := m.Values()
	require.Len(t, values, 3)
	require.Contains(t, values, "len")
	require.Contains(t, values, "sync")
	require.Contains(t, values, "map")
}

func TestMap_Keys(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"len", "sync", "map"} {
		m.Set(i, v)
	}

	keys := m.Keys()
	require.Len(t, keys, 3)
	require.Contains(t, keys, 0)
	require.Contains(t, keys, 1)
	require.Contains(t, keys, 2)
}

func TestMap_All(t *testing.T) {
	m := NewMap[int, string]()
	expected := map[int]string{0: "all", 1: "sync", 2: "map"}
	for k, v := range expected {
		m.Set(k, v)
	}

	got := make(map[int]string)
	for k, v := range m.All() {
		_, seen := got[k]
		require.False(t, seen, "All visited key %d twice", k)
		got[k] = v
	}

	require.Equal(t, expected, got)
}

func TestMap_All_Break(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"all", "sync", "map"} {
		m.Set(i, v)
	}

	count := 0
	for range m.All() {
		count++
		break
	}

	require.Equal(t, 1, count)
}

func TestMap_All_Empty(t *testing.T) {
	m := NewMap[int, string]()

	for k, v := range m.All() {
		t.Fatalf("All on empty map yielded %d=%s", k, v)
	}
}

func TestMap_AllKeys(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"all", "sync", "map"} {
		m.Set(i, v)
	}

	keys := make([]int, 0, 3)
	for k := range m.AllKeys() {
		keys = append(keys, k)
	}

	require.Len(t, keys, 3)
	require.Contains(t, keys, 0)
	require.Contains(t, keys, 1)
	require.Contains(t, keys, 2)
}

func TestMap_AllKeys_Break(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"all", "sync", "map"} {
		m.Set(i, v)
	}

	count := 0
	for range m.AllKeys() {
		count++
		break
	}

	require.Equal(t, 1, count)
}

func TestMap_AllValues(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"all", "sync", "map"} {
		m.Set(i, v)
	}

	values := make([]string, 0, 3)
	for v := range m.AllValues() {
		values = append(values, v)
	}

	require.Len(t, values, 3)
	require.Contains(t, values, "all")
	require.Contains(t, values, "sync")
	require.Contains(t, values, "map")
}

func TestMap_AllValues_Break(t *testing.T) {
	m := NewMap[int, string]()
	for i, v := range [3]string{"all", "sync", "map"} {
		m.Set(i, v)
	}

	count := 0
	for range m.AllValues() {
		count++
		break
	}

	require.Equal(t, 1, count)
}
