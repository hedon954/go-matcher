package reviewdog

import (
	"fmt"
	"os"
)

func TestReviewDog() {
	f, _ := os.Open("README.md")
	if f != nil {
		f.Close()
	}
	fmt.Println("hello", f)
}
