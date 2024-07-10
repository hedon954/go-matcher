package race

import (
	"testing"
)

func TestConcurrentOpMap(t *testing.T) {
	// 使用 go test -race 来检测潜在的并发问题的前提是测试会运行到相关代码。
	ConcurrentOpMap()
}
