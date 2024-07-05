package solana

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
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
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
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

func Test_GetBalance(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
	pubkey := "J27ma1MPBRvmPJxLqBqQGNECMXDm9L6abFa4duKiPosa"
	b, err := GetBalance(url, pubkey)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(b)
	}
}

func Test_GettokenAccountByOwner(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
	pubkey := "J27ma1MPBRvmPJxLqBqQGNECMXDm9L6abFa4duKiPosa"
	mint := "2FPyTwcZLUg1MDrwsyoP4D6s1tM7hAkHYRjkNb5w6Pxk"
	resp, err := GetTokenAccountsByOwner(url, pubkey, mint, "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Value[0].Account.Data.Parsed.Info.TokenAmount.UIAmount)
		fmt.Println(resp.Value[0].Account.Data.Parsed.Info.TokenAmount.Amount)
	}
}

func Test_SwapTransaction(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
	sig := "5Qcrof1h7VmL9P7g1M62zxMkPPbMMfdWYDaXEBSGJEEX6Xi6PLDuqnfibAB9KhBRcexcDjC1VFBVK6gXzEVJNWwW"
	out, err := GetTransaction(url, sig)
	if err != nil {
		fmt.Println(err)
	}

	param, succeed, err := ParseRaydiumSwapInstructionParam(out)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(succeed)
	fmt.Println(param)
}

func Test_TransferTransaction(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
	sig := "2PSzrxAmn7fHtRhNXK6RCNFzFR2uvN2CpY2T8tsnLJaFiiBHVuqtmekukr7zqDNCekj9TN5jhU4zq32RiTbgosPZ"
	// sig := "2ud2sUFqwdmYptgSBNZCvZ514tVSrQRTZnzjEErCVix5eZVUhymZfk7qE9QiZGM9PfiDqS4pH2GcfgzAZV2LJikK"

	out, err := GetTransaction(url, sig)
	if err != nil {
		fmt.Println(err)
		return
	}

	params, succeed, err := ParseTransferSOLInstructionParam(out)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(succeed)
	fmt.Println(params)
}

func Test_TransferSOL(t *testing.T) {
	url := "https://api.testnet.solana.com"
	wsurl := "wss://api.testnet.solana.com"
	from := ""
	to := "6huu25nWzFtBWPMQmWRzKLD4Wtfq11SSjZTU6oitLqdz"

	acc := solana.MustPrivateKeyFromBase58(from)
	fmt.Println(acc.PublicKey().String())
	sig, err := TransferSOL(url, wsurl, from, to, 100000000)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(sig)
	}
}
