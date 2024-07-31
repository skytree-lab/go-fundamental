package solana

import (
	"context"
	"encoding/json"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/machinebox/graphql"
	"github.com/skytree-lab/go-fundamental/util"
	"github.com/ybbus/jsonrpc/v3"
)

func GetTransaction(urls []string, signature string) (out *rpc.GetTransactionResult, err error) {
	type Tx struct {
		Commitment                     string `json:"commitment"`
		Encoding                       string `json:"encoding"`
		MaxSupportedTransactionVersion int    `json:"maxSupportedTransactionVersion"`
	}

	for _, url := range urls {
		rpcClient := jsonrpc.NewClient(url)
		resp, err := rpcClient.Call(context.Background(), "getTransaction", signature, &Tx{Commitment: "confirmed", Encoding: "json", MaxSupportedTransactionVersion: 0})
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTransaction err:%+v", err))
			continue
		}

		err = resp.GetObject(&out)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTransaction err:%+v", err))
			continue
		}
	}

	return
}

func GetBalance(urls []string, pubkey string) (uint64, error) {
	for _, url := range urls {
		rpcClient := jsonrpc.NewClient(url)
		resp, err := rpcClient.Call(context.Background(), "getBalance", pubkey)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("getBalance err:%+v", err))
			continue
		}

		balance := &GetBalanceResponse{}
		err = resp.GetObject(&balance)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetBalance err:%+v", err))
			continue
		}

		if balance == nil {
			util.Logger().Error("balance type err")
			continue
		}

		return balance.Value, nil
	}
	return 0, fmt.Errorf("getBalance, pubkey: %s", pubkey)
}

func GetTokenAccountBalance(urls []string, pubkey string) (string, float64, error) {
	for _, url := range urls {
		rpcClient := jsonrpc.NewClient(url)
		resp, err := rpcClient.Call(context.Background(), "getTokenAccountBalance", pubkey)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenAccountBalance err:%+v", err))
			continue
		}
		var balance *TokenBalance
		err = resp.GetObject(&balance)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenAccountBalance err:%+v", err))
			continue
		}

		if balance == nil {
			util.Logger().Error("balance type err")
			continue
		}

		return balance.Value.UIAmountString, balance.Value.UIAmount, nil
	}
	return "", 0, fmt.Errorf("getTokenAccountBalance err, pubkey: %s", pubkey)
}

func GetTokenAccountsByOwner(urls []string, pubkey string, mint string, programid string) (*GetTokenAccountsByOwnerResponse, error) {
	for _, url := range urls {
		rpcClient := jsonrpc.NewClient(url)
		min := &Mint{
			Mint: mint,
		}
		pro := &Program{
			ProgramId: programid,
		}
		encode := &Encoding{
			Encoding: "jsonParsed",
		}
		var resp *jsonrpc.RPCResponse
		var err error
		if mint != "" {
			resp, err = rpcClient.Call(context.Background(), "getTokenAccountsByOwner", pubkey, min, encode)
		} else if programid != "" {
			resp, err = rpcClient.Call(context.Background(), "getTokenAccountsByOwner", pubkey, pro, encode)
		}

		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenAccountsByOwner err:%+v", err))
			continue
		}

		var response *GetTokenAccountsByOwnerResponse
		err = resp.GetObject(&response)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("GetTokenAccountsByOwner err:%+v", err))
			continue
		}

		if response == nil {
			util.Logger().Error("response type err")
			continue
		}

		return response, nil
	}
	return nil, fmt.Errorf("getTokenAccountsByOwner err, pubkey:%s, mint:%s", pubkey, mint)
}

func GetTokeMeta(url string, mint string, key string) (out *Result, err error) {
	path := "/sol/v1/token/get_info?network=mainnet-beta&token_address="
	u := fmt.Sprintf("%s%s%s", url, path, mint)
	client := util.GetHTTPClient()
	header := make(map[string]string)
	header["x-api-key"] = key
	header["Content-Type"] = " application/json"
	header["Accept"] = " application/json"
	resp, err := util.HTTPReq("GET", u, client, nil, header)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokeMeta err:%+v", err))
		return nil, err
	}

	var meta MetaResponse
	err = json.Unmarshal(resp, &meta)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokeMeta err:%+v", err))
		return nil, err
	}

	if !meta.Success {
		err := fmt.Errorf("GetTokeMeta err:%s", mint)
		return nil, err
	}

	return meta.Result, nil
}

func GetPoolInfo(url string, key string, tokenA string, tokenB string) (*PoolInfoResponse, error) {
	path := "/v0/graphql/?api_key="
	u := fmt.Sprintf("%s%s%s", url, path, key)
	client := graphql.NewClient(u)
	q := fmt.Sprintf(`
query MyQuery {
  Raydium_LiquidityPoolv4(
    where: {
    baseMint: {_eq: "%s"},
    quoteMint: {_eq: "%s"}}
  ) {
    baseDecimal
    baseMint
    baseVault
    lpMint
    lpVault
    marketId
    marketProgramId
    openOrders
    quoteDecimal
    quoteMint
    quoteVault
    targetOrders
    withdrawQueue
    pubkey
  }
}`, tokenA, tokenB)
	req := graphql.NewRequest(q)
	var resp PoolInfoResponse
	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("graphql run err:%v", err))
		return nil, err
	}

	return &resp, nil
}

func BuildTransacion(ctx context.Context, clientRPC *rpc.Client, signers []solana.PrivateKey, instrs ...solana.Instruction) (*solana.Transaction, error) {
	recent, err := clientRPC.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		instrs,
		recent.Value.Blockhash,
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, payer := range signers {
				if payer.PublicKey().Equals(key) {
					return &payer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ExecuteInstructions(
	ctx context.Context,
	clientRPC *rpc.Client,
	signers []solana.PrivateKey,
	instrs ...solana.Instruction,
) (string, error) {
	tx, err := BuildTransacion(ctx, clientRPC, signers, instrs...)
	if err != nil {
		return "", err
	}

	sig, err := clientRPC.SendTransactionWithOpts(
		ctx,
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return "", err
	}
	return sig.String(), nil
}

func ExecuteInstructionsAndWaitConfirm(
	ctx context.Context,
	clientRPC *rpc.Client,
	RPCWs string,
	signers []solana.PrivateKey,
	instrs ...solana.Instruction,
) (string, error) {
	tx, err := BuildTransacion(ctx, clientRPC, signers, instrs...)
	if err != nil {
		return "", err
	}

	clientWS, err := ws.Connect(ctx, RPCWs)
	if err != nil {
		return "", err
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		clientRPC,
		clientWS,
		tx,
	)
	if err != nil {
		return "", err
	}

	return sig.String(), nil
}

func TransferSOL(urls []string, wsUrl string, fromKey string, to string, amount uint64) (sig string, err error) {
	for _, url := range urls {
		rpcClient := rpc.New(url)
		wsClient, err := ws.Connect(context.Background(), wsUrl)
		if err != nil {
			continue
		}

		accFrom := solana.MustPrivateKeyFromBase58(fromKey)
		accTo := solana.MustPublicKeyFromBase58(to)

		recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			continue
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(amount, accFrom.PublicKey(), accTo).Build(),
			},
			recent.Value.Blockhash,
			solana.TransactionPayer(accFrom.PublicKey()),
		)
		if err != nil {
			continue
		}

		_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
			if accFrom.PublicKey().Equals(key) {
				return &accFrom
			}
			return nil
		})
		if err != nil {
			continue
		}

		// Send transaction, and wait for confirmation:
		signature, err := confirm.SendAndConfirmTransaction(context.TODO(), rpcClient, wsClient, tx)
		if err != nil {
			continue
		}

		sig = signature.String()
		return sig, nil
	}
	return
}

func CheckAccount(url string, payer solana.PublicKey, publicKey solana.PublicKey, fromTokenAddr, toTokenAddr string) (map[string]solana.PublicKey, []solana.Instruction, error) {
	var mints []solana.PublicKey
	mints = append(mints, solana.MustPublicKeyFromBase58(fromTokenAddr))
	if fromTokenAddr != toTokenAddr {
		mints = append(mints, solana.MustPublicKeyFromBase58(toTokenAddr))
	}

	existingAccounts, missingAccounts, err := GetTokenAccountsFromMints(context.Background(), url, publicKey, mints...)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokenAccountsFromMints err:%v", err))
		return nil, nil, err
	}

	instrs := []solana.Instruction{}
	if len(missingAccounts) != 0 {
		for mint := range missingAccounts {
			if mint == NativeSOL {
				continue
			}
			inst, err := associatedtokenaccount.NewCreateInstruction(
				payer,
				publicKey,
				solana.MustPublicKeyFromBase58(mint),
			).ValidateAndBuild()
			if err != nil {
				util.Logger().Error(fmt.Sprintf("NewCreateInstruction err:%v", err))
				return nil, nil, err
			}
			instrs = append(instrs, inst)
		}
		for k, v := range missingAccounts {
			existingAccounts[k] = v
		}
	}
	return existingAccounts, instrs, nil
}

func GetTokenAccountsFromMints(ctx context.Context, url string, owner solana.PublicKey,
	mints ...solana.PublicKey) (map[string]solana.PublicKey, map[string]solana.PublicKey, error) {

	duplicates := map[string]bool{}
	tokenAccounts := []solana.PublicKey{}
	tokenAccountInfos := []TokenAccountInfo{}
	for _, m := range mints {
		if ok := duplicates[m.String()]; ok {
			continue
		}
		duplicates[m.String()] = true
		a, _, err := solana.FindAssociatedTokenAddress(owner, m)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("FindAssociatedTokenAddress err:%v", err))
			return nil, nil, err
		}
		// Use owner address for NativeSOL mint
		if m.String() == NativeSOL {
			a = owner
		}
		tokenAccounts = append(tokenAccounts, a)
		tokenAccountInfos = append(tokenAccountInfos, TokenAccountInfo{
			Mint:    m,
			Account: a,
		})
	}
	res, err := getMultiAccounts(url, tokenAccounts)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("getMultiAccounts err:%v", err))
		return nil, nil, err
	}

	missingAccounts := map[string]solana.PublicKey{}
	existingAccounts := map[string]solana.PublicKey{}
	for i, a := range res.Value {
		tai := tokenAccountInfos[i]
		if a == nil {
			missingAccounts[tai.Mint.String()] = tai.Account
			continue
		}
		if tai.Mint.String() == NativeSOL {
			existingAccounts[tai.Mint.String()] = owner
			continue
		}
		var ta token.Account
		err = bin.NewBinDecoder(a.Data.GetBinary()).Decode(&ta)
		if err != nil {
			util.Logger().Error(fmt.Sprintf("NewBinDecoder err:%v", err))
			return nil, nil, err
		}
		existingAccounts[tai.Mint.String()] = tai.Account
	}

	return existingAccounts, missingAccounts, nil
}

func TransferToken(urls []string, wsUrl string, amount uint64, senderSPLTokenAccount, mint, recipient solana.PublicKey, sender *solana.PrivateKey) (sig string, err error) {
	for _, url := range urls {
		rpcClient := rpc.New(url)
		wsClient, err := ws.Connect(context.Background(), wsUrl)
		if err != nil {
			continue
		}

		existingAccounts, instrs, err := CheckAccount(url, sender.PublicKey(), recipient, mint.String(), mint.String())
		if err != nil {
			continue
		}

		var instructions []solana.Instruction
		instructions = append(instructions, instrs...)

		recipientSPLTokenAccount, ok := existingAccounts[mint.String()]
		if !ok {
			continue
		}

		var signaturers []solana.PublicKey
		signaturers = append(signaturers, sender.PublicKey())

		transfer, err := token.NewTransferInstruction(amount,
			senderSPLTokenAccount,
			recipientSPLTokenAccount,
			sender.PublicKey(),
			signaturers).ValidateAndBuild()
		if err != nil {
			continue
		}

		instructions = append(instructions, transfer)
		recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			continue
		}
		trx, err := solana.NewTransaction(instructions,
			recent.Value.Blockhash,
			solana.TransactionPayer(sender.PublicKey()))
		if err != nil {
			continue
		}

		_, err = trx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
			if key == sender.PublicKey() {
				return sender
			}
			return nil
		})
		if err != nil {
			continue
		}

		signature, err := confirm.SendAndConfirmTransaction(context.Background(), rpcClient, wsClient, trx)
		if err != nil {
			continue
		}
		sig = signature.String()
		return sig, nil
	}
	return
}

func getMultiAccounts(url string, accs []solana.PublicKey) (out *rpc.GetMultipleAccountsResult, err error) {
	rpcClient := jsonrpc.NewClient(url)
	params := []interface{}{accs}
	resp, err := rpcClient.Call(context.Background(), "getMultipleAccounts", params)
	if err != nil {
		return
	}

	err = resp.GetObject(&out)
	if err != nil {
		return
	}
	return
}

func GetMultiAccounts(urls []string, accounts []*PoolTokenPairAccount) (resp *rpc.GetMultipleAccountsResult, err error) {
	var accs []solana.PublicKey
	for _, acc := range accounts {
		accs = append(accs, acc.BaseMint)
		accs = append(accs, acc.QuoteMint)
	}

	for _, url := range urls {
		resp, err = getMultiAccounts(url, accs)
		if err != nil {
			continue
		}

		if resp == nil {
			continue
		}
		return
	}
	return
}
