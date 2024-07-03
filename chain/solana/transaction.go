package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/machinebox/graphql"
	"github.com/skytree-lab/go-fundamental/util"
	"github.com/ybbus/jsonrpc/v3"
)

func GetTransaction(url string, signature string) (out *rpc.GetTransactionResult, err error) {
	rpcClient := jsonrpc.NewClient(url)

	type Tx struct {
		Commitment                     string `json:"commitment"`
		Encoding                       string `json:"encoding"`
		MaxSupportedTransactionVersion int    `json:"maxSupportedTransactionVersion"`
	}

	resp, err := rpcClient.Call(context.Background(), "getTransaction", signature, &Tx{Commitment: "confirmed", Encoding: "json", MaxSupportedTransactionVersion: 0})
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTransaction err:%+v", err))
		return
	}

	err = resp.GetObject(&out)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTransaction err:%+v", err))
		return
	}
	return
}

func GetBalance(url string, pubkey string) (uint64, error) {
	rpcClient := jsonrpc.NewClient(url)
	resp, err := rpcClient.Call(context.Background(), "getBalance", pubkey)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("getBalance err:%+v", err))
		return 0, err
	}

	var balance *GetBalanceResponse
	err = resp.GetObject(&balance)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetBalance err:%+v", err))
		return 0, err
	}

	if balance == nil {
		util.Logger().Error("balance type err")
		return 0, nil
	}

	return balance.Value, nil
}

func GetTokenAccountBalance(url string, pubkey string) (string, float64, error) {
	rpcClient := jsonrpc.NewClient(url)
	resp, err := rpcClient.Call(context.Background(), "getTokenAccountBalance", pubkey)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokenAccountBalance err:%+v", err))
		return "", 0, err
	}
	var balance *TokenBalance
	err = resp.GetObject(&balance)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokenAccountBalance err:%+v", err))
		return "", 0, err
	}

	if balance == nil {
		util.Logger().Error("balance type err")
		return "", 0, nil
	}

	return balance.Value.UIAmountString, balance.Value.UIAmount, nil
}

func GetTokenAccountsByOwner(url string, pubkey string, mint string, programid string) (*GetTokenAccountsByOwnerResponse, error) {
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
		return nil, err
	}

	var response *GetTokenAccountsByOwnerResponse
	err = resp.GetObject(&response)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokenAccountsByOwner err:%+v", err))
		return nil, err
	}

	if response == nil {
		util.Logger().Error("response type err")
		return nil, nil
	}

	if len(response.Value) <= 0 {
		return nil, nil
	}

	return response, nil
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
