package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/skytree-lab/go-fundamental/chain/eth/contract"
	"github.com/skytree-lab/go-fundamental/util"
)

func Parse(urls []string, hash string, addr string) (*big.Int, error) {
	txHash := common.HexToHash(hash)
	for _, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}

		tx, _, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil || tx == nil {
			continue
		}

		if tx.To().String() == addr {
			return tx.Value(), nil
		}
		return nil, nil
	}
	return nil, nil
}

func GetBalance(urls []string, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	for _, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}

		balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			continue
		}

		return balance, nil
	}
	return nil, nil
}

func GetTokenBalance(urls []string, token string, addr string) (amount *big.Int, err error) {
	var client *ethclient.Client
	var instance *contract.Usdt
	tokenAddr := common.HexToAddress(token)
	user := common.HexToAddress(addr)
	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenBalance Dial err is: %+v", err))
			continue
		}
		instance, err = contract.NewUsdt(tokenAddr, client)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenBalance NewUsdt err is: %+v", err))
			continue
		}
		amount, err = instance.BalanceOf(&bind.CallOpts{}, user)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenBalance BalanceOf err is: %+v", err))
			continue
		}
		return
	}
	return
}

func TransferETH(urls []string, fromKey string, to string, amount *big.Int) (string, error) {
	toAddress := common.HexToAddress(to)
	value := amount
	privateKey, err := crypto.HexToECDSA(fromKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", nil
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasLimit := uint64(21000) // in units

	for _, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			continue
		}

		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			continue
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			continue
		}

		var data []byte
		tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			continue
		}

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			return "", err
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			continue
		}

		return signedTx.Hash().Hex(), nil
	}
	return "", nil
}

func TransferErc20Token(urls []string, fromKey string, to string, amount *big.Int, token string, chainid uint64) (hash string, succeed bool, err error) {
	var client *ethclient.Client
	var erc20 *contract.UsdtTransactor
	var tx *types.Transaction
	var opts *bind.TransactOpts
	tokenAddr := common.HexToAddress(token)
	privateKey, err := crypto.HexToECDSA(fromKey)
	if err != nil {
		errMsg := fmt.Sprintf("TransferErc20Token err, reason=[%s]", err)
		util.Logger().Error(errMsg)
		return
	}
	addr, err := util.PrivateToAddress(fromKey)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("TransferErc20Token Dial err is: %+v", err))
		return
	}

	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("TransferErc20Token Dial err is: %+v", err))
			continue
		}

		opts, err = util.CreateTransactionOpts(client, privateKey, chainid, common.HexToAddress(addr), nil)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("TransferErc20Token CreateTransactionOpts err is: %+v", err))
			continue
		}

		erc20, err = contract.NewUsdtTransactor(tokenAddr, client)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("TransferErc20Token NewUsdtTransactor err is: %+v", err))
			continue
		}

		tx, err = erc20.Transfer(opts, common.HexToAddress(to), amount)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("TransferErc20Token Transfer err is: %+v", err))
			continue
		}
		_, succeed, err = util.TxWaitToSync(opts.Context, client, tx)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("TransferErc20Token Transfer err is: %+v", err))
			continue
		}

		if succeed {
			hash = tx.Hash().String()
			return
		}
	}
	return
}

func FetchErc20TokenMeta(urls []string, token string) (decimal int, symbol string, name string, err error) {
	var client *ethclient.Client
	var instance *contract.Usdt
	var decimalBig *big.Int
	tokenAddr := common.HexToAddress(token)
	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("FetchErc20TokenMeta Dial err is: %+v", err))
			continue
		}
		instance, err = contract.NewUsdt(tokenAddr, client)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("FetchErc20TokenMeta NewUsdt err is: %+v", err))
			continue
		}
		decimalBig, err = instance.Decimals(&bind.CallOpts{})
		if err != nil {
			util.Logger().Error(fmt.Sprintf("FetchErc20TokenMeta Decimals err is: %+v", err))
			continue
		}
		decimal = int(decimalBig.Int64())
		name, err = instance.Name(&bind.CallOpts{})
		if err != nil {
			util.Logger().Error(fmt.Sprintf("FetchErc20TokenMeta Name err is: %+v", err))
			continue
		}
		symbol, err = instance.Symbol(&bind.CallOpts{})
		if err != nil {
			util.Logger().Error(fmt.Sprintf("FetchErc20TokenMeta Name err is: %+v", err))
			continue
		}
		return
	}
	return
}

func FetchPoolPrice(urls []string, base string, baseDecimal int, quote string, quoteDecimal int, pool string) (price float64, err error) {
	var baseAmount *big.Int
	var quoteAmount *big.Int
	baseAmount, err = GetTokenBalance(urls, base, pool)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("FetchPoolPrice GetTokenBalance err: %+v", err))
		return
	}

	quoteAmount, err = GetTokenBalance(urls, quote, pool)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("FetchPoolPrice GetTokenBalance err: %+v", err))
		return
	}

	if baseAmount.Uint64() == 0 || quoteAmount.Uint64() == 0 {
		return
	}

	price = util.BigIntDiv(baseAmount, baseDecimal, quoteAmount, quoteDecimal)
	return
}
