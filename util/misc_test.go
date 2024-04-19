package util

import (
	"fmt"
	"testing"
)

func TestRemoveSlice(t *testing.T) {
	var a []int
	a = append(a, 0)
	a = append(a, 11)
	a = append(a, 31)
	a = append(a, 21)

	b := RemoveIndex(a, 3)
	fmt.Println(b)

	tpl := GetEthCallPostData("0x123", "234")
	fmt.Println(tpl)

	lpadstr, _ := Pad("71", 64, '0')
	fmt.Println(lpadstr)
}
