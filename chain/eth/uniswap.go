package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/skytree-lab/go-fundamental/chain/eth/contract"
	"github.com/skytree-lab/go-fundamental/util"
)

func HandleApprove(urls []string, chainid uint64, token0 string, key string, router string, swapValue *big.Int) (succeed bool, err error) {
	owner, err := util.PrivateToAddress(key)
	if err != nil {
		return
	}
	allowance, _ := GetAllowance(urls, token0, owner, router)
	if allowance != nil && allowance.Cmp(swapValue) >= 0 {
		succeed = true
		return
	}
	approvehash, approvesucceed, err := Approve(urls, chainid, token0, key, router, swapValue)
	if err != nil || !approvesucceed {
		return
	}
	status := checkTransactionStatus(urls, approvehash)
	if status == 0 {
		return
	}
	succeed = true
	return
}

func SwapInUni(urls []string, chainid uint64, uinifactory string, unirouter string, key string, swapValue *big.Int, token0 string, token1 string) (hash string, succeed bool, err error) {
	wallet := helper.InitWallet(key)
	if wallet == nil {
		err = errors.New("SwapInUni helper.InitWallet err")
		return
	}
	d := time.Now().Add(time.Minute * time.Duration(5)).Unix()
	deadline := big.NewInt(d)
	exactInputSingleParams := contract.ISwapRouterExactInputSingleParams{
		TokenIn:           common.HexToAddress(token0),
		TokenOut:          common.HexToAddress(token1),
		Fee:               big.NewInt(int64(constants.FeeMedium)),
		Recipient:         wallet.PublicKey,
		Deadline:          deadline,
		AmountIn:          swapValue,
		AmountOutMinimum:  big.NewInt(0),
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

func SwapInUniBase(urls []string, chainid uint64, uinifactory string, unirouter string, key string, swapValue *big.Int, token0 string, token1 string) (hash string, succeed bool, err error) {
	wallet := helper.InitWallet(key)
	if wallet == nil {
		err = errors.New("SwapInUni helper.InitWallet err")
		return
	}

	extractInputSingleParams := contract.IV3SwapRouterExactInputSingleParams{
		TokenIn:           common.HexToAddress(token0),
		TokenOut:          common.HexToAddress(token1),
		Fee:               big.NewInt(int64(constants.FeeMedium)),
		Recipient:         wallet.PublicKey,
		AmountIn:          swapValue,
		AmountOutMinimum:  big.NewInt(0),
		SqrtPriceLimitX96: big.NewInt(0),
	}

	var client *ethclient.Client
	var router *contract.UnibaserouterTransactor
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
		router, err = contract.NewUnibaserouterTransactor(common.HexToAddress(unirouter), client)
		if err != nil {
			continue
		}
		tx, err = router.ExactInputSingle(opts, extractInputSingleParams)
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

func GetPoolAddress(urls []string, uinifactory string, token0, token1 common.Address) (poolAddr common.Address, valid bool, err error) {
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
		poolAddr, err = f.GetPool(nil, token0, token1, big.NewInt(int64(constants.FeeMedium)))
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

func ParseUniTransaction(urls []string, tx string) (status uint64, amount0 *big.Int, amount1 *big.Int, err error) {
	var receipt *types.Receipt
	var client *ethclient.Client
	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			continue
		}
		receipt, err = client.TransactionReceipt(context.Background(), common.HexToHash(tx))
		if receipt == nil {
			continue
		}
		amount0, amount1, err = parseUniSwapLog(receipt.Logs)
		if err != nil {
			continue
		}
		status = receipt.Status
		return
	}
	return
}

func parseUniSwapLog(logs []*types.Log) (amount0 *big.Int, amount1 *big.Int, err error) {
	var contractAbi abi.ABI
	for _, log := range logs {
		if len(log.Topics) == 0 || log.Topics[0].Hex() != uniPoolSwappedTopic {
			continue
		}

		contractAbi, err = abi.JSON(strings.NewReader(string(contract.UnipoolABI)))
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Not found abi json, err:%+v", err))
			return
		}

		uniswapEvent := &contract.UnipoolSwap{}
		err = contractAbi.UnpackIntoInterface(uniswapEvent, "Swap", log.Data)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("parseUniSwapLog failed to UnpackIntoInterface: %+v", err))
			return
		}
		amount0 = uniswapEvent.Amount0
		amount1 = uniswapEvent.Amount1
		return
	}
	return
}
