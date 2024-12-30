package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/skytree-lab/go-fundamental/chain/eth/contract"
	"github.com/skytree-lab/go-fundamental/util"
)

func SwapInUni(urls []string, chainid uint64, uinifactory string, unirouter string, key string, swapValue *big.Int, token0 *coreEntities.Token, token1 *coreEntities.Token, slippage float64) (hash string, succeed bool, err error) {
	wallet := helper.InitWallet(key)
	if wallet == nil {
		err = errors.New("SwapInUni helper.InitWallet err")
		return
	}

	approvehash, approvesucceed, err := Approve(urls, chainid, token0.Address.Hex(), key, unirouter, swapValue)
	if err != nil {
		return
	}
	if !approvesucceed {
		return
	}

	status := checkTransactionStatus(urls, approvehash)
	if status == 0 {
		return
	}

	d := time.Now().Add(time.Minute * time.Duration(5)).Unix()
	deadline := big.NewInt(d)
	amountOut := big.NewInt(0)

	exactInputSingleParams := contract.ISwapRouterExactInputSingleParams{
		TokenIn:           common.HexToAddress(token0.Address.String()),
		TokenOut:          common.HexToAddress(token1.Address.String()),
		Fee:               big.NewInt(int64(constants.FeeMedium)),
		Recipient:         wallet.PublicKey,
		Deadline:          deadline,
		AmountIn:          swapValue,
		AmountOutMinimum:  amountOut,
		SqrtPriceLimitX96: big.NewInt(0),
	}

	var client *ethclient.Client
	var router *contract.UnirouterTransactor
	var tx *types.Transaction
	var opts *bind.TransactOpts
	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			continue
		}
		opts, err = util.CreateTransactionOpts(client, wallet.PrivateKey, chainid, wallet.PublicKey, nil)
		if err != nil {
			continue
		}
		router, err = contract.NewUnirouterTransactor(common.HexToAddress(unirouter), client)
		if err != nil {
			continue
		}
		tx, err = router.ExactInputSingle(opts, exactInputSingleParams)
		if err != nil {
			continue
		}
		_, succeed, err = util.TxWaitToSync(opts.Context, client, tx)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("SwapInUni transaction err: %+v", err))
			continue
		}
		if succeed {
			hash = tx.Hash().String()
			return
		}
	}
	return
}

func GetPoolAddress(urls []string, uinifactory string, token0, token1 common.Address, fee *big.Int) (poolAddr common.Address, valid bool, err error) {
	var client *ethclient.Client
	var f *contract.Unifactory
	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			continue
		}
		f, err = contract.NewUnifactory(common.HexToAddress(uinifactory), client)
		if err != nil {
			continue
		}
		poolAddr, err = f.GetPool(nil, token0, token1, fee)
		if err != nil {
			continue
		}
		if poolAddr == (common.Address{}) {
			continue
		}
		valid = true
		return
	}
	return
}

func FetchPoolPrice(urls []string, base string, baseDecimal int, quote string, quoteDecimal int, pooladdr string) (price float64, err error) {
	var baseAmount *big.Int
	var quoteAmount *big.Int
	baseAmount, err = GetTokenBalance(urls, base, pooladdr)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("FetchPoolPrice GetTokenBalance err: %+v", err))
		return
	}

	quoteAmount, err = GetTokenBalance(urls, quote, pooladdr)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("FetchPoolPrice GetTokenBalance err: %+v", err))
		return
	}

	if baseAmount.Uint64() == 0 || quoteAmount.Uint64() == 0 {
		return
	}

	baseVal := util.ConvertTokenAmountToFloat64(baseAmount.String(), int32(baseDecimal))
	quoteVal := util.ConvertTokenAmountToFloat64(quoteAmount.String(), int32(quoteDecimal))
	price = quoteVal / baseVal
	return
}

func checkTransactionStatus(urls []string, tx string) (status uint64) {
	for i := 0; i < 300; i++ {
		time.Sleep(1 * time.Second)
		url := urls[i%len(urls)]
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}
		receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx))
		if err != nil {
			continue
		}
		if receipt == nil {
			continue
		}
		status = receipt.Status
		return
	}
	return
}
