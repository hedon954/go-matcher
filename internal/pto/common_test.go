package pto

import (
	"testing"

	"github.com/hedon954/go-matcher/fixtures"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestFromGameResultJson(t *testing.T) {
	t.Run("integer should work", func(t *testing.T) {
		res, err := FromGameResultJson[int64](&GameResult{
			Result: []byte("1"),
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), *res)
	})

	t.Run("string should work", func(t *testing.T) {
		res, err := FromGameResultJson[string](&GameResult{
			Result: []byte(`"hello"`),
		})
		assert.Nil(t, err)
		assert.Equal(t, "hello", *res)
	})

	t.Run("struct should work", func(t *testing.T) {
		type Info struct {
			Name string `json:"name"`
		}
		res, err := FromGameResultJson[Info](&GameResult{
			Result: []byte(`{"name": "hello"}`),
		})
		assert.Nil(t, err)
		assert.Equal(t, "hello", res.Name)
	})

	t.Run("invalid json should fail", func(t *testing.T) {
		res, err := FromGameResultJson[int64](&GameResult{
			Result: []byte("hello"),
		})
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestMustFromGameResultJson(t *testing.T) {
	t.Run("integer should work", func(t *testing.T) {
		res := MustFromGameResultJson[int64](&GameResult{
			Result: []byte("1"),
		})
		assert.Equal(t, int64(1), *res)
	})

	t.Run("string should work", func(t *testing.T) {
		res := MustFromGameResultJson[string](&GameResult{
			Result: []byte(`"hello"`),
		})
		assert.Equal(t, "hello", *res)
	})

	t.Run("invalid json should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromGameResultJson[int64](&GameResult{
				Result: []byte("hello"),
			})
		})
	})
}

func TestFromGameResultPb(t *testing.T) {
	req := fixtures.Request{
		Name: "hedon",
	}
	bs, err := proto.Marshal(&req)
	assert.Nil(t, err)

	t.Run("proto struct should work", func(t *testing.T) {
		data, err := FromGameResultPb[fixtures.Request](&GameResult{
			Result: bs,
		})
		assert.Nil(t, err)
		assert.Equal(t, "hedon", data.Name)
	})

	t.Run("invalid proto should fail", func(t *testing.T) {
		data, err := FromGameResultPb[fixtures.Request](&GameResult{
			Result: []byte("hello"),
		})
		assert.NotNil(t, err)
		assert.Nil(t, data)
	})

	t.Run("basic type should fail", func(t *testing.T) {
		data, err := FromGameResultPb[int64](&GameResult{
			Result: bs,
		})
		assert.NotNil(t, err)
		assert.Nil(t, data)
	})
}

//nolint:dupl
func TestMustFromGameResultPb(t *testing.T) {
	req := fixtures.Request{
		Name: "hedon",
	}
	bs, err := proto.Marshal(&req)
	assert.Nil(t, err)

	t.Run("proto struct should work", func(t *testing.T) {
		data := MustFromGameResultPb[fixtures.Request](&GameResult{
			Result: bs,
		})
		assert.Equal(t, "hedon", data.Name)
	})

	t.Run("invalid proto should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromGameResultPb[fixtures.Request](&GameResult{
				Result: []byte("hello"),
			})
		})
	})

	t.Run("invalid type should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromGameResultPb[int64](&GameResult{
				Result: bs,
			})
		})
	})
}
