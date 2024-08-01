package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	type Item struct {
		Name string
	}

	// the type of T should be ineterface or pointer
	manager := New[int, *Item]()
	assert.Equal(t, 0, len(manager.items))

	// no item, get should return nil
	get := manager.Get(1)
	assert.Nil(t, get)

	// no item, delete should return nil
	d := manager.Delete(1)
	assert.Nil(t, d)

	// get after add, should return the item
	manager.Add(1, &Item{Name: "1"})
	item := manager.Get(1)
	assert.NotNil(t, item)
	assert.Equal(t, "1", item.Name)
	assert.Equal(t, 1, len(manager.items))
	assert.Equal(t, true, manager.Exists(1))

	// delete after add, should return the item
	d = manager.Delete(1)
	assert.NotNil(t, d)
	assert.Equal(t, "1", d.Name)
	assert.Equal(t, 0, len(manager.items))
	assert.Equal(t, false, manager.Exists(1))
}
