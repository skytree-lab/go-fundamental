package util

import (
	"fmt"
	"testing"
)

func Test_csa(t *testing.T) {
	var strs []string
	strs = append(strs, "aaa")
	strs = append(strs, "bbb")
	strs = append(strs, "ccc")
	strs = append(strs, "ddd")

	csa := CreateCycledStringArray(strs)
	for i := 0; i < 100; i++ {
		fmt.Println(csa.GetCurrentString())
		fmt.Println(csa.CurIdx)
		fmt.Println(csa.GetCurrentString())
		fmt.Println(csa.CurIdx)
		fmt.Println(csa.GetCurrentString())
		fmt.Println(csa.CurIdx)
		fmt.Println(csa.GetCurrentString())
		fmt.Println(csa.CurIdx)
	}
}
