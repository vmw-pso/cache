package cache

import (
	"testing"
	"time"
)

type TestStruct struct {
	Name string
}

func TestCache(t *testing.T) {
	c := NewCache[TestStruct](15 * time.Second)

	item, err := c.ReadById(1)
	if err == nil || item.Name != "" {
		t.Errorf("found value that does not exist. wanted: \"\", got: '%s'\n", item.Name)
	}

	result, err := c.ReadAll()
	if err != nil || len(result) > 0 {
		t.Errorf("found values that should not exist. wanted length: 0, got: %d\n", len(result))
	}

	c.Add(TestStruct{"John Doe"}, 1, 60*60)

	item, err = c.ReadById(1)
	if err != nil || item.Name != "John Doe" {
		t.Errorf("no item found when expected. wanted: `John Doe`, got: %s\n", item.Name)
	}

	c.Add(TestStruct{"Jane Doe"}, 2, 60*60)

	result, err = c.ReadAll()
	if err != nil || len(result) != 2 {
		t.Errorf("expected length 2, got %d\n", len(result))
	}

	c.DeleteById(2)
	result, err = c.ReadAll()
	if err != nil || len(result) != 1 {
		t.Errorf("expected length 2, got %d\n", len(result))
	}
}
