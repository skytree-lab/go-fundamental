package eth

import (
	"fmt"
	"testing"
)

func Test_FetchPoolPrice(t *testing.T) {
	price, err := FetchPoolPrice([]string{"https://ethereum-sepolia.publicnode.com"}, "0x779877A7B0D9E8603169DdbD7836e478b4624789", 18, "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14", 18, "0xDD7CC9a0dA070fB8B60dC6680b596133fb4A7100")
	fmt.Println(err)
	fmt.Println(price)
}
