package codereview

import (
	"fmt"
	"os"
)

func TestCodeReview() {
	f, _ := os.Open("README.md")
	if f != nil {
		fmt.Println("file exists", f)
		f.Close()
	}
	fmt.Println("codereview", f)
}
