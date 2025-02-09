package typeconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToSlice(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		m := map[int]int{1: 1, 2: 2, 3: 3}
		s := MapToSlice(m)
		assert.ElementsMatch(t, []int{1, 2, 3}, s)
	})

	t.Run("string", func(t *testing.T) {
		m := map[string]string{"1": "1", "2": "2", "3": "3"}
		s := MapToSlice(m)
		assert.ElementsMatch(t, []string{"1", "2", "3"}, s)
	})

	t.Run("struct", func(t *testing.T) {
		m := map[string]struct {
			ID   int
			Name string
		}{"1": {ID: 1, Name: "1"}, "2": {ID: 2, Name: "2"}, "3": {ID: 3, Name: "3"}}
		s := MapToSlice(m)
		assert.ElementsMatch(t, []string{"1", "2", "3"}, s)
	})
}

func TestSliceToMap(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		s := []int{1, 2, 3}
		m := SliceToMap(s)
		assert.Equal(t, map[int]bool{1: true, 2: true, 3: true}, m)
	})

	t.Run("string", func(t *testing.T) {
		s := []string{"1", "2", "3"}
		m := SliceToMap(s)
		assert.Equal(t, map[string]bool{"1": true, "2": true, "3": true}, m)
	})

	t.Run("struct", func(t *testing.T) {
		type item struct {
			ID   int
			Name string
		}
		s := []item{
			{ID: 1, Name: "1"},
			{ID: 2, Name: "2"},
			{ID: 3, Name: "3"},
		}
		m := SliceToMap(s)
		assert.Equal(t, map[item]bool{
			{ID: 1, Name: "1"}: true,
			{ID: 2, Name: "2"}: true,
			{ID: 3, Name: "3"}: true,
		}, m)
	})
}
