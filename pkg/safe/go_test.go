package safe_test

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/hedon954/go-matcher/pkg/safe"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

type logger struct {
	sync.Mutex
	buffer bytes.Buffer
}

func (l *logger) Error(msg string, v ...any) {
	l.Lock()
	defer l.Unlock()
	l.buffer.WriteString(msg)
	if len(v) > 0 {
		l.buffer.WriteString(": ")
		l.buffer.WriteString(cast.ToString(v))
	}
}

func TestSafe(t *testing.T) {
	l := &logger{}
	safe.SetLogger(l)

	safe.Go(func() {
		panic("panic in safe.Go")
	})

	start := time.Now().UnixMilli()
	safe.Call(func() {
		time.Sleep(10 * time.Millisecond)
		panic("panic in safe.Call")
	})
	safe.Wait()
	end := time.Now().UnixMilli()

	assert.True(t, end-start >= 10)

	buffer := l.buffer.String()
	assert.Contains(t, buffer, "safe call occurs panic: panic in safe.Call")
	assert.Contains(t, buffer, "safe go occurs panic: panic in safe.Go")
}
