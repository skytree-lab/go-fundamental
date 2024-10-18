package eth

import (
	"fmt"
	"testing"

	"github.com/skytree-lab/go-fundamental/util"
)

func Test_FetchPoolPrice(t *testing.T) {
	price, err := FetchPoolPrice([]string{"https://ethereum-sepolia.publicnode.com"}, "0x779877A7B0D9E8603169DdbD7836e478b4624789", 18, "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14", 18, "0xDD7CC9a0dA070fB8B60dC6680b596133fb4A7100")
	fmt.Println(err)
	fmt.Println(price)
}

func Test_GetAllowance(t *testing.T) {
	price, err := GetAllowance([]string{"http://127.0.0.1:8545/"}, "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14", "0x8f5f0cf2d3b2d635250d46e5d3aa4dcb1ced3b54", "0x3bFA4769FB09eefC5a80d6E87c3B9C650f7Ae48E")
	fmt.Println(err)
	fmt.Println(price)
}

func Test_Approve(t *testing.T) {
	urls := []string{"http://127.0.0.1:8545/"}
	chainid := uint64(31337)
	token := "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	ownerKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	spender := "0x50cf1849e32e6a17bbff6b1aa8b1f7b479ad6c12"
	spend := util.ConvertFloat64ToTokenAmount(1000, 18)
	hash, succeed, err := Approve(urls, chainid, token, ownerKey, spender, spend)
	fmt.Println(err)
	fmt.Println(succeed)
	fmt.Println(hash)
}
