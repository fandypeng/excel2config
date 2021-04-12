package helper

import (
	"fmt"
	"testing"
)

func TestGenerateRandomStr(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(GenerateRandomStr(32))
	}
}
