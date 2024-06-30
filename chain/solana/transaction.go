package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/skytree-lab/go-fundamental/util"
	"github.com/streamingfast/solana-go/rpc"
	"github.com/ybbus/jsonrpc/v3"
)

func GetTransaction(c *rpc.Client, signature string, commitmentType *rpc.CommitmentType) (out *rpc.GetTransactionResponse, err error) {
	opts := map[string]interface{}{
		"encoding":                       "json",
		"maxSupportedTransactionVersion": 0,
	}
	if commitmentType != nil {
		opts["Commitment"] = *commitmentType
	}
	params := []interface{}{signature, opts}
	err = c.DoRequest(&out, "getTransaction", params...)
	return
}

type TokenBalance struct {
	Context *Context `json:"context"`
	Value   *Value   `json:"value"`
}
type Context struct {
	Slot int `json:"slot"`
}
type Value struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
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

type Mint struct {
	Mint string `json:"mint"`
}

type Encoding struct {
	Encoding string `json:"encoding"`
}

type GetTokenAccountsByOwnerResponse struct {
	Context struct {
		Slot int `json:"slot"`
	} `json:"context"`
	Value []struct {
		Account struct {
			Data struct {
				Parsed struct {
					Info struct {
						IsNative    bool   `json:"isNative"`
						Mint        string `json:"mint"`
						Owner       string `json:"owner"`
						State       string `json:"state"`
						TokenAmount struct {
							Amount         string  `json:"amount"`
							Decimals       int     `json:"decimals"`
							UIAmount       float64 `json:"uiAmount"`
							UIAmountString string  `json:"uiAmountString"`
						} `json:"tokenAmount"`
					} `json:"info"`
					Type string `json:"type"`
				} `json:"parsed"`
				Program string `json:"program"`
				Space   int    `json:"space"`
			} `json:"data"`
			Executable bool   `json:"executable"`
			Lamports   int    `json:"lamports"`
			Owner      string `json:"owner"`
		} `json:"account"`
		Pubkey string `json:"pubkey"`
	} `json:"value"`
}

func GetTokenAccountsByOwner(url string, pubkey string, mint string) (string, float64, error) {
	rpcClient := jsonrpc.NewClient(url)
	min := &Mint{
		Mint: mint,
	}

	encode := &Encoding{
		Encoding: "jsonParsed",
	}
	resp, err := rpcClient.Call(context.Background(), "getTokenAccountsByOwner", pubkey, min, encode)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokenAccountsByOwner err:%+v", err))
		return "", 0, err
	}

	var balance *GetTokenAccountsByOwnerResponse
	err = resp.GetObject(&balance)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("GetTokenAccountsByOwner err:%+v", err))
		return "", 0, err
	}

	if balance == nil {
		util.Logger().Error("balance type err")
		return "", 0, nil
	}
	if len(balance.Value) <= 0 {
		return "", 0, nil
	}

	return balance.Value[0].Account.Data.Parsed.Info.TokenAmount.UIAmountString, balance.Value[0].Account.Data.Parsed.Info.TokenAmount.UIAmount, nil
}

type MetaResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Result  *Result `json:"result"`
}

type Result struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
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

type RaydiumLiquidityPoolv4 struct {
	BaseDecimal     int    `json:"baseDecimal"`
	BaseMint        string `json:"baseMint"`
	BaseVault       string `json:"baseVault"`
	LpMint          string `json:"lpMint"`
	LpVault         string `json:"lpVault"`
	MarketId        string `json:"marketId"`
	MarketProgramId string `json:"marketProgramId"`
	OpenOrders      string `json:"openOrders"`
	QuoteDecimal    int    `json:"quoteDecimal"`
	QuoteMint       string `json:"quoteMint"`
	QuoteVault      string `json:"quoteVault"`
	TargetOrders    string `json:"targetOrders"`
	WithdrawQueue   string `json:"withdrawQueue"`
	Pubkey          string `json:"pubkey"`
}

type PoolInfoResponse struct {
	RaydiumLiquidityPoolv4 []*RaydiumLiquidityPoolv4 `json:"Raydium_LiquidityPoolv4"`
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
