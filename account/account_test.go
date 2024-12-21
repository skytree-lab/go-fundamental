package account

import (
	"fmt"
	"testing"
)

func Test_NewSuiAccount(t *testing.T) {
	acc := NewSuiAccount()
	fmt.Println(acc)
}
