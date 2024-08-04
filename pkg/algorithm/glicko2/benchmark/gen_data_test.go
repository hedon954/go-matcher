package benchmark

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2/example"
)

// Test_Generate_Group_Data 用于生成 groups 模拟数据，确保多种匹配算法用的同一套数据
func Test_Generate_Group_Data(t *testing.T) {
	groupCounts := []int{100, 1000, 10000, 100000, 1000000}
	for _, gc := range groupCounts {
		filename := fmt.Sprintf("groups-%d.json", gc)
		_ = os.Remove(filename)
		f, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
			return
		}
		groups := make([]*example.Group, gc)
		for i := 0; i < gc; i++ {
			var players []*example.Player
			count := rand.Intn(5) + 1
			for j := 0; j < count; j++ {
				p := example.NewPlayer(uuid.NewString(), false, 0, 0,
					glicko2.Args{
						MMR: 1000 + float64(rand.Intn(1000)),
						RD:  0,
						V:   0,
					})
				players = append(players, p)
			}
			groups[i] = example.NewGroup(fmt.Sprintf("Group%d", i+1), players)
		}
		bs, _ := json.Marshal(groups)
		_, _ = f.Write(bs)
		_ = f.Close()
	}
}
