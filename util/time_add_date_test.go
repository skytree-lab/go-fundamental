package util

import (
	"fmt"
	"testing"
	"time"
)

func Test_addDate(tt *testing.T) {
	t := time.Date(2024, 10, 31, 23, 59, 59, 0, time.UTC)
	at := AddDate(t, 0, 1, 0)
	fmt.Println(at)
}
