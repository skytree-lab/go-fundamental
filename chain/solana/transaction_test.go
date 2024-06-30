package solana

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/streamingfast/solana-go/rpc"
)

func Test_GetTokenMeta(t *testing.T) {
	meta, _ := GetTokeMeta("https://api.shyft.to", "4DWaPEVY3E3bkG2APWS13wRKRiQeCpz4G4ZGVuyCYJU9", "")
	fmt.Println(meta)
}

func Test_GetPoolInfo(t *testing.T) {
	resp, _ := GetPoolInfo("https://programs.shyft.to", "", "4DWaPEVY3E3bkG2APWS13wRKRiQeCpz4G4ZGVuyCYJU9", "So11111111111111111111111111111111111111112")
	fmt.Println(resp)
}

func Test_GettokenAccountBalance(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/"
	pubkey := "2ZsNAdu4kzkRPs89P4EZjvRzq1BfdTgBhMrtDkWAUg2X"
	uiAmount, amount, err := GetTokenAccountBalance(url, pubkey)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(uiAmount)
		fmt.Println(amount)
	}

	pubkey = "7UZ8VjMTYF1yBraryJscXQu8wREyHBomZD223PyrJn36"
	uiAmount, amount, err = GetTokenAccountBalance(url, pubkey)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(uiAmount)
		fmt.Println(amount)
	}
}

func Test_GettokenAccountByOwner(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/"
	pubkey := "J27ma1MPBRvmPJxLqBqQGNECMXDm9L6abFa4duKiPosa"
	mint := "2FPyTwcZLUg1MDrwsyoP4D6s1tM7hAkHYRjkNb5w6Pxk"
	uiAmount, amount, err := GetTokenAccountsByOwner(url, pubkey, mint)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(uiAmount)
		fmt.Println(amount)
	}
}

func Test_Transaction(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5"
	sig := "2PSzrxAmn7fHtRhNXK6RCNFzFR2uvN2CpY2T8tsnLJaFiiBHVuqtmekukr7zqDNCekj9TN5jhU4zq32RiTbgosPZ"
	// sig := "2ud2sUFqwdmYptgSBNZCvZ514tVSrQRTZnzjEErCVix5eZVUhymZfk7qE9QiZGM9PfiDqS4pH2GcfgzAZV2LJikK"
	rp := rpc.NewClient(url)
	conf := rpc.CommitmentConfirmed
	out, err := GetTransaction(rp, sig, &conf)
	if err != nil {
		fmt.Println(err)
	}
	for _, instruction := range out.Transaction.Message.Instructions {
		datas, err := base58.Decode(instruction.Data)
		if err != nil {
			continue
		}

		if len(datas) < 8 {
			continue
		}

		t := new(big.Int)
		var typeBuf []byte
		for i := 3; i >= 0; i-- {
			typeBuf = append(typeBuf, datas[i])
		}

		t.SetBytes(typeBuf[:])
		instype := t.Uint64()
		if instype != 2 {
			continue
		}

		sourceIdx := instruction.Accounts[0]
		destIdx := instruction.Accounts[1]
		fmt.Println(out.Transaction.Message.AccountKeys[sourceIdx].String())
		fmt.Println(out.Transaction.Message.AccountKeys[destIdx].String())

		var amountBuf []byte
		for i := 7; i >= 4; i-- {
			amountBuf = append(amountBuf, datas[i])
		}

		amount := new(big.Int)
		amount.SetBytes(amountBuf[:])
		fmt.Println(amount.Uint64())
	}
}
