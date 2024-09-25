package solana

import "github.com/gagliardetto/solana-go"

var (
	RaydiumLiquidityPoolv4ProgramID = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	RaydiumCpmmPoolProgramID        = solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")
	NativeSOL                       = "11111111111111111111111111111111"
)

type PoolTokenPairAccount struct {
	BaseMint      solana.PublicKey
	BaseDecimals  int
	BaseAmount    uint64
	QuoteMint     solana.PublicKey
	QuoteDecimals int
	QuoteAmount   uint64
}

type TokenAccountInfo struct {
	Mint    solana.PublicKey
	Account solana.PublicKey
}

type TransferSOLInstructionParam struct {
	Source      string
	Destination string
	Amount      uint64
}

type RaydiumSwapInstructionParam struct {
	SwapIn  *RaydiumSwapInnerInstructionData
	SwapOut *RaydiumSwapInnerInstructionData
}

type RaydiumSwapInnerInstructionData struct {
	Source       string
	Destionation string
	Owner        string
	Amount       uint64
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

type Mint struct {
	Mint string `json:"mint"`
}

type Program struct {
	ProgramId string `json:"programId"`
}

type Encoding struct {
	Encoding string `json:"encoding"`
}

type GetBalanceResponse struct {
	Context struct {
		APIVersion string `json:"apiVersion"`
		Slot       int    `json:"slot"`
	} `json:"context"`
	Value uint64 `json:"value"`
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
