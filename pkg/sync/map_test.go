package sync

import "testing"

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
		t.Fatalf("error occured in range %+v", err)
	}

	for k, v := range checkData {
		if !v.checked {
			t.Fatalf("key %d was not applied in range", k)
		}
	}
}
