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

func CreateTransactionOpts(client *ethclient.Client, key *ecdsa.PrivateKey, chainId uint64, caller common.Address) (opts *bind.TransactOpts, err error) {
	nonce, err := client.PendingNonceAt(context.Background(), caller)
	if err != nil {
		errMsg := fmt.Sprintf("CreateTransactionOpts:client.PendingNonceAt err: %+v", err)
		util.Logger().Error(errMsg)
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		errMsg := fmt.Sprintf("CreateTransactionOpts:client.SuggestGasPrice err: %+v", err)
		util.Logger().Error(errMsg)
		return nil, err
	}

	srcChainID := big.NewInt(int64(chainId))
	opts, err = bind.NewKeyedTransactorWithChainID(key, srcChainID)
	if err != nil {
		errMsg := fmt.Sprintf("CreateTransactionOpts:NewKeyedTransactorWithChainID err: %+v", err)
		util.Logger().Error(errMsg)
		return nil, err
	}

	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0) // in wei
	opts.GasLimit = uint64(0)  // in units
	opts.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))

	return opts, nil
}

func TxWaitToSync(ctx context.Context, client *ethclient.Client, tx *types.Transaction) (*types.Receipt, bool, error) {
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		errMsg := fmt.Sprintf("TxWaitToSync:bind.WaitMine err: %+v", err)
		util.Logger().Error(errMsg)
		return nil, false, err
	}

	return receipt, receipt.Status == types.ReceiptStatusSuccessful, nil
}
