package eth

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/skytree-lab/go-fundamental/util"
)

func Test_GetPoolAddress(t *testing.T) {
	urls := []string{"https://eth-sepolia.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5"}
	factory := "0x0227628f3F023bb0B980b67D528571c95c6DaC1c"

	uniaddr := "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"
	weth9addr := "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14"

	pooladdr, valid, err := GetPoolAddress(urls, factory, common.HexToAddress(weth9addr), common.HexToAddress(uniaddr))
	fmt.Println(err)
	fmt.Println(valid)
	fmt.Println(pooladdr.Hex())

	pooladdr, valid, err = GetPoolAddress(urls, factory, common.HexToAddress(uniaddr), common.HexToAddress(weth9addr))
	fmt.Println(err)
	fmt.Println(valid)
	fmt.Println(pooladdr.Hex())
}

func Test_SwapInUni(t *testing.T) {
	urls := []string{"https://eth-sepolia.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5"}
	chainid := uint64(11155111)
	factory := "0x0227628f3F023bb0B980b67D528571c95c6DaC1c"
	router := "0x65669fE35312947050C450Bd5d36e6361F85eC12"
	key := ""
	// amount := util.ConvertFloat64ToTokenAmount(0.1, 18)
	amount := util.ConvertFloat64ToTokenAmount(0.002263498835798428, 18)

	uniaddr := "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"
	weth9addr := "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14"

	hash, succeed, err := SwapInUni(urls, chainid, factory, router, key, amount, uniaddr, weth9addr)
	fmt.Println(err)
	fmt.Println(succeed)
	fmt.Println(hash)
}

func Test_ParseUniSwapTransaction(t *testing.T) {
	urls := []string{"https://eth-sepolia.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5"}
	tx := "0xcfa14d0e43afa00302cfa21a14099b9ee8f1d6c13e1f8af44141e7338fbdd39d"
	status, amount0, amount1, err := ParseUniTransaction(urls, tx)
	fmt.Println(status)
	fmt.Println(amount0.Text(10))
	fmt.Println(amount1.Text(10))
	fmt.Println(err)
}
