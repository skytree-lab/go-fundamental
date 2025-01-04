package util

import (
	"fmt"
	"math"
	"math/big"
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

func Test_GetTokenPriceUSDT(t *testing.T) {
	urlbase := "https://www.okx.com/api/v5/market/index-tickers?instId="
	base := "sol"
	quote := "usdt"
	price, err := GetTokenPriceUSDT(urlbase, base, quote)
	fmt.Println(err)
	fmt.Println(price)
}

func Test_ConvertBigIntTokenAmountToFloat64(t *testing.T) {
	a := 10000 * math.Pow10(18)
	fmt.Println(a)
	amount := big.NewInt(int64(1 * math.Pow10(18)))
	amt := ConvertBigIntTokenAmountToFloat64(amount, 18)
	fmt.Println(amt)
}
