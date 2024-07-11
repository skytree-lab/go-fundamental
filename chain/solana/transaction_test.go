package solana

import (
	"context"
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
	uiAmount, amount, err := GetTokenAccountBalance([]string{url}, pubkey)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(uiAmount)
		fmt.Println(amount)
	}

	pubkey = "7UZ8VjMTYF1yBraryJscXQu8wREyHBomZD223PyrJn36"
	uiAmount, amount, err = GetTokenAccountBalance([]string{url}, pubkey)
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
	b, err := GetBalance([]string{url}, pubkey)
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
	resp, err := GetTokenAccountsByOwner([]string{url}, pubkey, mint, "")
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
	out, err := GetTransaction([]string{url}, sig)
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

	out, err := GetTransaction([]string{url}, sig)
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
	url := "https://solana-mainnet.g.alchemy.com/v2/"
	wsurl := "wss://api.mainnet-beta.solana.com"
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

func Test_TransferSpl(t *testing.T) {
	url := "https://api.mainnet-beta.solana.com"
	wsurl := "wss://api.mainnet-beta.solana.com"
	to := solana.MustPublicKeyFromBase58("J2YwmwqMfCE4LW1T4Xy5R9owBn1t2yhaMQ7TPrS2aZau")
	sender := solana.MustPrivateKeyFromBase58("")
	mint := solana.MustPublicKeyFromBase58("4G86CMxGsMdLETrYnavMFKPhQzKTvDBYGMRAdVtr72nu")
	serderSpl, _, err := solana.FindAssociatedTokenAddress(sender.PublicKey(), mint)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(serderSpl)

	sig, err := TransferToken(context.Background(), url, wsurl, 100000000, serderSpl, mint, to, &sender)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sig)
}

func Test_GetMultipleAccounts(t *testing.T) {
	urls := []string{
		"https://go.getblock.io/f79d62f3bd954dbda2fad5f2be26e91e",
		"https://solana-mainnet.g.alchemy.com/v2/v-8DDF-sKRNirIxSFn9rdszEKw_vu0i5",
		"https://solana-mainnet.g.alchemy.com/v2/1m6yZduJJzCft9xSRO6ut5bOSXF7CSfH",
	}
	acc1 := &PoolTokenPairAccount{
		BaseMint:  solana.MustPublicKeyFromBase58("ABLksYkz92eK1AbZvxwgfma6Zoz1fKnzhgVGpwBWNQyk"),
		QuoteMint: solana.MustPublicKeyFromBase58("6Q1hGQVEzL8dZCjn6Vb5jvJj61vozD8BRxDWQh6ZAAgY"),
	}
	acc2 := &PoolTokenPairAccount{
		BaseMint:  solana.MustPublicKeyFromBase58("C8EC7PEZehgcPzdRADLzk12M6bnYTgV5mtvt6Kdjok61"),
		QuoteMint: solana.MustPublicKeyFromBase58("DFcjeot4H54xfwP5GCtxs1SHBfLfXVESpgwN14gh1xWV"),
	}

	accs := []*PoolTokenPairAccount{
		acc1,
		acc2,
	}

	out, err := GetMultiAccounts(urls, accs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}
