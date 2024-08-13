package typeconv

import (
	"testing"

	"github.com/hedon954/go-matcher/fixtures"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestFromProto(t *testing.T) {
	t.Run("empty bytes should failed", func(t *testing.T) {
		data, err := FromProto[fixtures.Request]([]byte{})
		assert.Equal(t, "protobuf data is empty", err.Error())
		assert.Nil(t, data)
	})

	t.Run("proto struct should work", func(t *testing.T) {
		req := fixtures.Request{
			Name: "hedon",
		}
		bs, err := proto.Marshal(&req)
		assert.Nil(t, err)
		data, err := FromProto[fixtures.Request](bs)
		assert.Nil(t, err)
		assert.Equal(t, "hedon", data.Name)
	})

	t.Run("invalid proto should fail", func(t *testing.T) {
		data, err := FromProto[fixtures.Request]([]byte("hello"))
		assert.NotNil(t, err)
		assert.Nil(t, data)
	})

	t.Run("basic type should failed", func(t *testing.T) {
		data, err := FromProto[int64]([]byte("1"))
		assert.NotNil(t, err)
		assert.Nil(t, data)
	})
}

//nolint:dupl
func TestMustFromProto(t *testing.T) {
	req := fixtures.Request{
		Name: "hedon",
	}
	bs, err := proto.Marshal(&req)
	assert.Nil(t, err)

	t.Run("proto struct should work", func(t *testing.T) {
		data := MustFromProto[fixtures.Request](bs)
		assert.Equal(t, "hedon", data.Name)
	})

	t.Run("invalid proto should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromProto[fixtures.Request]([]byte("hello"))
		})
	})

	t.Run("invalid type should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustFromProto[int64](bs)
		})
	})
}
