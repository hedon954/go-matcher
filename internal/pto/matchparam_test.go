package pto

import (
	"testing"

	"github.com/hedon954/go-matcher/fixtures"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestMustFromAttrJson(t *testing.T) {
	t.Run("integer should work", func(t *testing.T) {
		res := MustFromAttrJson[int64](&UploadPlayerAttr{
			Extra: []byte("1"),
		})
		assert.Equal(t, int64(1), *res)
	})

	t.Run("string should work", func(t *testing.T) {
		res := MustFromAttrJson[string](&UploadPlayerAttr{
			Extra: []byte(`"hello"`),
		})
		assert.Equal(t, "hello", *res)
	})

	t.Run("bool should work", func(t *testing.T) {
		res := MustFromAttrJson[bool](&UploadPlayerAttr{
			Extra: []byte("true"),
		})
		assert.Equal(t, true, *res)
	})

	t.Run("float should work", func(t *testing.T) {
		res := MustFromAttrJson[float64](&UploadPlayerAttr{
			Extra: []byte("3.14"),
		})
		assert.Equal(t, 3.14, *res)
	})

	t.Run("struct should work", func(t *testing.T) {
		type Info struct {
			Name string `json:"name"`
		}
		res := MustFromAttrJson[Info](&UploadPlayerAttr{
			Extra: []byte(`{"name": "hello"}`),
		})
		assert.Equal(t, "hello", res.Name)
	})

	t.Run("map should work", func(t *testing.T) {
		res := MustFromAttrJson[map[string]string](&UploadPlayerAttr{
			Extra: []byte(`{"name": "hello"}`),
		})
		assert.Equal(t, "hello", (*res)["name"])
	})

	t.Run("slice should work", func(t *testing.T) {
		res := MustFromAttrJson[[]int64](&UploadPlayerAttr{
			Extra: []byte(`[1,2,3,4]`),
		})
		assert.Equal(t, int64(1), (*res)[0])
	})

	t.Run("invalid json should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromAttrJson[int64](&UploadPlayerAttr{
				Extra: []byte("hello"),
			})
		})
	})
}

func TestFromAttrJson(t *testing.T) {
	t.Run("struct should work", func(t *testing.T) {
		type Info struct {
			Name string `json:"name"`
		}
		res, err := FromAttrJson[Info](&UploadPlayerAttr{
			Extra: []byte(`{"name": "hello"}`),
		})
		assert.Nil(t, err)
		assert.Equal(t, "hello", res.Name)
	})

	t.Run("invalid string should fail", func(t *testing.T) {
		res, err := FromAttrJson[map[string]string](&UploadPlayerAttr{
			Extra: []byte(`{"name": hello"}`),
		})
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestFromAttrPb(t *testing.T) {
	t.Run("proto struct should work", func(t *testing.T) {
		req := fixtures.Request{
			Name: "hedon",
		}
		bs, err := proto.Marshal(&req)
		assert.Nil(t, err)
		data, err := FromAttrPb[fixtures.Request](&UploadPlayerAttr{
			Extra: bs,
		})
		assert.Nil(t, err)
		assert.Equal(t, "hedon", data.Name)
	})

	t.Run("invalid proto should fail", func(t *testing.T) {
		data, err := FromAttrPb[fixtures.Request](&UploadPlayerAttr{
			Extra: []byte("hello"),
		})
		assert.NotNil(t, err)
		assert.Nil(t, data)
	})

	t.Run("basic type should failed", func(t *testing.T) {
		data, err := FromAttrPb[int64](&UploadPlayerAttr{
			Extra: []byte("1"),
		})
		assert.NotNil(t, err)
		assert.Nil(t, data)
	})
}

//nolint:dupl
func TestMustFromAttrPb(t *testing.T) {
	req := fixtures.Request{
		Name: "hedon",
	}
	bs, err := proto.Marshal(&req)
	assert.Nil(t, err)

	t.Run("proto struct should work", func(t *testing.T) {
		data := MustFromAttrPb[fixtures.Request](&UploadPlayerAttr{
			Extra: bs,
		})
		assert.Equal(t, "hedon", data.Name)
	})

	t.Run("invalid proto should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromAttrPb[fixtures.Request](&UploadPlayerAttr{
				Extra: []byte("hello"),
			})
		})
	})

	t.Run("invalid type should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromAttrPb[int64](&UploadPlayerAttr{
				Extra: bs,
			})
		})
	})
}
