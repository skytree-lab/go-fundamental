package eth

import (
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

func Weth9Deposit(urls []string, chainid uint64, weth9address string, key string, amount *big.Int) (hash string, succeed bool, err error) {
	var client *ethclient.Client
	var weth9 *contract.Weth9Transactor
	var tx *types.Transaction
	var opts *bind.TransactOpts
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		errMsg := fmt.Sprintf("Weth9Deposit err:%+v", err)
		util.Logger().Error(errMsg)
		return
	}
	addr, err := util.PrivateToAddress(key)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("Weth9Deposit err: %+v", err))
		return
	}

	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Deposit Dial err: %+v", err))
			continue
		}

		opts, err = util.CreateTransactionOpts(client, privateKey, chainid, common.HexToAddress(addr), amount)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Deposit CreateTransactionOpts err: %+v", err))
			continue
		}

		weth9, err = contract.NewWeth9Transactor(common.HexToAddress(weth9address), client)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Deposit NewWeth9Transactor err: %+v", err))
			continue
		}

		tx, err = weth9.Deposit(opts)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Deposit Deposit err: %+v", err))
			continue
		}
		_, succeed, err = util.TxWaitToSync(opts.Context, client, tx)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Deposit transaction err: %+v", err))
			continue
		}

		if succeed {
			hash = tx.Hash().String()
			return
		}
	}
	return
}

func Weth9Withdraw(urls []string, chainid uint64, weth9address string, key string, amount *big.Int) (hash string, succeed bool, err error) {
	var client *ethclient.Client
	var weth9 *contract.Weth9Transactor
	var tx *types.Transaction
	var opts *bind.TransactOpts
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		errMsg := fmt.Sprintf("Weth9Withdraw err:%+v", err)
		util.Logger().Error(errMsg)
		return
	}
	addr, err := util.PrivateToAddress(key)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("Weth9Withdraw err: %+v", err))
		return
	}

	for _, url := range urls {
		client, err = ethclient.Dial(url)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Withdraw Dial err: %+v", err))
			continue
		}

		opts, err = util.CreateTransactionOpts(client, privateKey, chainid, common.HexToAddress(addr), nil)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Withdraw CreateTransactionOpts err: %+v", err))
			continue
		}

		weth9, err = contract.NewWeth9Transactor(common.HexToAddress(weth9address), client)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Withdraw NewWeth9Transactor err: %+v", err))
			continue
		}

		tx, err = weth9.Withdraw(opts, amount)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Withdraw Withdraw err: %+v", err))
			continue
		}
		_, succeed, err = util.TxWaitToSync(opts.Context, client, tx)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Weth9Withdraw transaction err: %+v", err))
			continue
		}

		if succeed {
			hash = tx.Hash().String()
			return
		}
	}
	return
}
