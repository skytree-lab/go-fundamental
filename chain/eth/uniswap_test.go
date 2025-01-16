package eth

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_GetPoolAddress(t *testing.T) {
	urls := []string{"HTTP://127.0.0.1:8545"}
	factory := "0x1F98431c8aD98523631AE4a59f267346ea31F984"

	uniaddr := "0xA35923162C49cF95e6BF26623385eb431ad920D3"
	weth9addr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

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
	urls := []string{"HTTP://127.0.0.1:8545"}
	chainid := uint64(31337)
	factory := "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	router := "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	key := ""
	// amount := util.ConvertFloat64ToTokenAmount(0.1, 18)
	amount := new(big.Int)
	amount.SetString("40412950793038993432951", 10)

	token0 := "0xA35923162C49cF95e6BF26623385eb431ad920D3"
	token1 := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

	ok, err := HandleApprove(urls, chainid, token0, key, router, amount)
	fmt.Println(err)
	fmt.Println(ok)

	hash, succeed, err := SwapInUni(urls, chainid, factory, router, key, amount, token0, token1)
	fmt.Println(err)
	fmt.Println(succeed)
	fmt.Println(hash)
}

func Test_SwapInUnibase(t *testing.T) {
	urls := []string{"HTTP://127.0.0.1:8545"}
	chainid := uint64(31337)
	factory := "0x33128a8fC17869897dcE68Ed026d694621f6FDfD"
	router := "0x2626664c2603336E57B271c5C0b26F421741e481"
	key := ""
	amount1 := new(big.Int)
	amount, _ := amount1.SetString("51912787592331905195651300", 10)
	// amount := util.ConvertFloat64ToTokenAmount(0.002263498835798428, 18)

	token1 := "0x4200000000000000000000000000000000000006"
	token0 := "0x2Da56AcB9Ea78330f947bD57C54119Debda7AF71"

	ok, err := HandleApprove(urls, chainid, token0, key, router, amount)
	fmt.Println(err)
	fmt.Println(ok)

	hash, succeed, err := SwapInUniBase(urls, chainid, factory, router, key, amount, token0, token1)
	fmt.Println(err)
	fmt.Println(succeed)
	fmt.Println(hash)
}

func Test_ParseUniSwapTransaction(t *testing.T) {
	urls := []string{"HTTP://127.0.0.1:8545"}
	tx := "0xcfa14d0e43afa00302cfa21a14099b9ee8f1d6c13e1f8af44141e7338fbdd39d"
	status, amount0, amount1, err := ParseUniTransaction(urls, tx)
	fmt.Println(status)
	fmt.Println(amount0.Text(10))
	fmt.Println(amount1.Text(10))
	fmt.Println(err)
}
