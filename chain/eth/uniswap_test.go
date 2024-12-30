package eth

import (
	"fmt"
	"math/big"
	"testing"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/skytree-lab/go-fundamental/util"
)

func Test_GetPoolAddress(t *testing.T) {
	urls := []string{"https://eth-sepolia.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5"}
	factory := "0x0227628f3F023bb0B980b67D528571c95c6DaC1c"

	uniaddr := "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"
	weth9addr := "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14"

	pooladdr, valid, err := GetPoolAddress(urls, factory, common.HexToAddress(weth9addr), common.HexToAddress(uniaddr), big.NewInt(int64(constants.FeeMedium)))
	fmt.Println(err)
	fmt.Println(valid)
	fmt.Println(pooladdr.Hex())

	pooladdr, valid, err = GetPoolAddress(urls, factory, common.HexToAddress(uniaddr), common.HexToAddress(weth9addr), big.NewInt(int64(constants.FeeMedium)))
	fmt.Println(err)
	fmt.Println(valid)
	fmt.Println(pooladdr.Hex())
}

func Test_SwapInUni(t *testing.T) {
	urls := []string{"https://eth-sepolia.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5"}
	chainid := uint64(11155111)
	factory := "0x0227628f3F023bb0B980b67D528571c95c6DaC1c"
	router := "0x65669fE35312947050C450Bd5d36e6361F85eC12"
	key := "16bd0714f3c50f3f24b6db7053e9876604a4ebe590b2aaec63cf8bd16d604e62"
	slippage := 0.005
	// amount := util.ConvertFloat64ToTokenAmount(0.1, 18)
	amount := util.ConvertFloat64ToTokenAmount(0.002263498835798428, 18)

	uniaddr := "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"
	weth9addr := "0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14"
	UNI := coreEntities.NewToken(uint(chainid), common.HexToAddress(uniaddr), 18, "UNI", "Uniswap")
	WETH := coreEntities.NewToken(uint(chainid), common.HexToAddress(weth9addr), 18, "WETH", "Wrapped Ether")

	hash, succeed, err := SwapInUni(urls, chainid, factory, router, key, amount, UNI, WETH, slippage)
	fmt.Println(err)
	fmt.Println(succeed)
	fmt.Println(hash)
}
