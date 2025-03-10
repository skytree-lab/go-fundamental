package solana

import (
	"fmt"
	"testing"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

func Test_GetTokenMeta(t *testing.T) {
	meta, _ := GetTokenMeta("https://api.shyft.to", "4DWaPEVY3E3bkG2APWS13wRKRiQeCpz4G4ZGVuyCYJU9", []string{""})
	fmt.Println(meta)
}

func Test_GetPoolDecimals(t *testing.T) {
	resp, _ := GetTokenDecimal([]string{"https://solana-mainnet.g.alchemy.com/v2/alch-demo"}, "4DWaPEVY3E3bkG2APWS13wRKRiQeCpz4G4ZGVuyCYJU9")
	fmt.Println(resp)
}

func Test_GetPoolInfo(t *testing.T) {
	resp, _ := GetPoolInfo("https://programs.shyft.to", []string{""}, "4DWaPEVY3E3bkG2APWS13wRKRiQeCpz4G4ZGVuyCYJU9", "So11111111111111111111111111111111111111112")
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
	pubkey := "6huu25nWzFtBWPMQmWRzKLD4Wtfq11SSjZTU6oitLqdz"
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

	param, succeed, err := ParseRaydiumSwapInstructionParam(out, []string{url})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(succeed)
	fmt.Println(param)
}

func Test_TransferTransaction(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
	// sig := "2PSzrxAmn7fHtRhNXK6RCNFzFR2uvN2CpY2T8tsnLJaFiiBHVuqtmekukr7zqDNCekj9TN5jhU4zq32RiTbgosPZ"
	sig := "yEPhnF66CMGMjtCCUcnSJXakGsbGXwzsT1QxPoYP3gUNGcod5ZkMfJrXBmLDawsMEmAzGXDuzFowShTmAyepGTU"

	out, err := GetTransaction([]string{url}, sig)
	if err != nil {
		fmt.Println(err)
		return
	}

	params, succeed, err := ParseTransferSOLInstructionParam(out, []string{url})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(succeed)
	fmt.Println(params)
}

func Test_GetLatestBlockhash(t *testing.T) {
	urls := "https://api.mainnet-beta.solana.com"
	count, err := GetLatestBlockhash(urls)
	fmt.Println(count)
	fmt.Println(err)
}

func Test_TransferCpTransaction(t *testing.T) {
	url := "https://solana-mainnet.g.alchemy.com/v2/alch-demo"
	// sig := "2PSzrxAmn7fHtRhNXK6RCNFzFR2uvN2CpY2T8tsnLJaFiiBHVuqtmekukr7zqDNCekj9TN5jhU4zq32RiTbgosPZ"
	sig := "yEPhnF66CMGMjtCCUcnSJXakGsbGXwzsT1QxPoYP3gUNGcod5ZkMfJrXBmLDawsMEmAzGXDuzFowShTmAyepGTU"

	out, err := GetTransaction([]string{url}, sig)
	if err != nil {
		fmt.Println(err)
		return
	}

	params, succeed, err := ParseRaydiumCpSwapInstructionParam(out, []string{url})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(succeed)
	fmt.Println(params)
}

func Test_GetSPLTokenTotalSupply(t *testing.T) {
	urls := []string{"https://api.mainnet-beta.solana.com"}
	token := "FFZoPGvkgUqPCHyvKxyfeJQz5PqbRnWM2tS9BjWcpump"
	value, err := GetSPLTokenTotalSupply(urls, token)
	fmt.Println(value)
	fmt.Println(err)
}

func Test_GetSPLTokenLargestAccount(t *testing.T) {
	urls := []string{"https://api.mainnet-beta.solana.com"}
	token := "88WiNYkPTFvY1UwdoYo89C2JMRAJByTKPabzhJxpcq5q"
	count, err := GetSPLTokenTopAccounts(urls, token)
	fmt.Println(count)
	fmt.Println(err)
}

func Test_TransferSOL(t *testing.T) {
	urls := []string{"https://api.mainnet-beta.solana.com"}
	wsurl := "wss://api.mainnet-beta.solana.com"
	from := ""
	to := "6huu25nWzFtBWPMQmWRzKLD4Wtfq11SSjZTU6oitLqdz"
	acc := solana.MustPrivateKeyFromBase58(from)
	fmt.Println(acc.PublicKey().String())
	sig, err := TransferSOL(urls, wsurl, from, to, 5000000)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(sig)
	}
}

func Test_TransferSpl(t *testing.T) {
	urls := []string{"https://api.mainnet-beta.solana.com"}
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

	sig, err := TransferToken(urls, wsurl, 100000000, serderSpl, mint, to, &sender)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sig)
}

func Test_GetMultipleAccounts(t *testing.T) {
	urls := []string{
		"https://solana-mainnet.g.alchemy.com/v2/alch-demo",
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

	for _, acc := range out.Value {
		var poolCoinBalance token.Account
		err = bin.NewBinDecoder(acc.Data.GetBinary()).Decode(&poolCoinBalance)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(poolCoinBalance.Mint)
		fmt.Println(poolCoinBalance.Owner)
		fmt.Println("------------------")
	}

	fmt.Println(out)
}

func Test_ProcessTransactionWithAddressLookups(t *testing.T) {
	// sig := "5EnN3PHQou83bnsLJrWALSTs9LK4gxVPMJZa6k3c77yLuNRDPSY42FYa49FYt4u53SGQ3Ti32pxjY6sCPxzXV9FU"
	// sig := "58khGDwFBpRgV1zqe8z43Q58rqcPatmWMiWzYca4mqapCCQk4MZexxUd4Ju3DTEo4FUfaBhBQAEjqt4jU3CiGSWj"
	sig := "WpWnyWgMVAhwCxqsQMmiCjxuzLVDLC1Xb15qtNrcy5NTVcmZDNjwDdmcfHaLhqZmYeqtUy6sH6B6GubtNRcyyGL"
	urls := []string{
		"https://solana-mainnet.g.alchemy.com/v2/alch-demo",
	}
	out, err := GetTransaction(urls, sig)
	if err != nil {
		fmt.Println(err)
		return
	}

	if out == nil || out.Transaction == nil {
		err := fmt.Errorf("tx not found, txid:%s", sig)
		fmt.Println(err)
		return
	}

	tx, err := out.Transaction.GetTransaction()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, key := range tx.Message.AccountKeys {
		fmt.Println(key)
	}

	err = ProcessTransactionWithAddressLookups(tx, urls)
	if err != nil {
		fmt.Println(err)
		return
	}

	slice, err := tx.Message.GetAllKeys()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("-------------------")
	for _, key := range slice {
		fmt.Println(key)
	}
}

func Test_GetTokenMint(t *testing.T) {
	mint := "7qvaWYSZgzx2ingvoLpHPJVDBDKDzRfBtRrFdjzJwLeH"
	// mint := "9gCFWhH4NUgqrVGRD4v7kBUC9nunq6DuEAJumRGNiD8"
	urls := []string{
		"https://solana-mainnet.g.alchemy.com/v2/alch-demo",
	}
	tokenmint, err := GetTokenMint(urls, mint)
	if err != nil {
		fmt.Println(err)
		return
	}

	if tokenmint == nil {
		err := fmt.Errorf("token not found, mint:%s", mint)
		fmt.Println(err)
		return
	}

	fmt.Println(tokenmint)
}

func Test_GetTokenDecimals(t *testing.T) {
	// mint := "7qvaWYSZgzx2ingvoLpHPJVDBDKDzRfBtRrFdjzJwLeH"
	mint := "9gCFWhH4NUgqrVGRD4v7kBUC9nunq6DuEAJumRGNiD8"
	urls := []string{
		"https://solana-mainnet.g.alchemy.com/v2/alch-demo",
	}
	tokenmint, err := GetTokenMint(urls, mint)
	if err != nil {
		fmt.Println(err)
		return
	}

	if tokenmint == nil {
		err := fmt.Errorf("token not found, mint:%s", mint)
		fmt.Println(err)
		return
	}

	fmt.Println(tokenmint.Decimals)
}
