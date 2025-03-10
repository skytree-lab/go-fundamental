package util

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/skytree-lab/go-fundamental/core"
)

var IsAlphanumeric = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func ConvertHexToDecimalInStringFormat(hexString string) string {
	i := new(big.Int)
	// if hexString with '0x' prefix, using fmt.Sscan()
	fmt.Sscan(hexString, i)
	// if hexString without '0x' prefix, using i.SetString()
	//i.SetString(hexString, 16)

	return fmt.Sprintf("%v", i)
}

func ConvertFloat64ToTokenAmount(amount float64, decimals int) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(amount)

	fp := math.Pow10(decimals)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(int64(fp)))
	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func PadLeft(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) >= length {
			return str[0:length]
		}
	}
}

func IsAnAddress(address string) bool {
	return len(address) == core.AddressFixedLength+2 && address[:2] == "0x" && IsAlphanumeric(address)
}

func IsValidTxHash(txHash string) bool {
	return len(txHash) == core.TxHashFixedLength && txHash[:2] == "0x" && IsAlphanumeric(txHash)
}

func ConvertTokenAmountToFloat64(amt string, tokenDecimal int32) float64 {
	amount, _ := decimal.NewFromString(amt)
	amount_converted := amount.Div(decimal.New(1, tokenDecimal))
	amountFloat, _ := amount_converted.Float64()
	return amountFloat
}

func ConvertFloatStringTokenAmountToBigInt(amt string, tokenDecimal int32) *big.Int {
	amount, _ := decimal.NewFromString(amt)
	amount_converted := amount.Mul(decimal.New(1, tokenDecimal))
	return amount_converted.BigInt()
}

func ConvertBigIntTokenAmountToFloat64(b *big.Int, tokenDecimal int32) float64 {
	amount := decimal.NewFromBigInt(b, 0)
	amount_converted := amount.Div(decimal.New(1, tokenDecimal))
	amountFloat, _ := amount_converted.Float64()
	return amountFloat
}

func GetBigIntFromString(v0 string) (n0 *big.Int, err error) {
	n0 = new(big.Int)
	n0, ok := n0.SetString(v0, 10)
	if !ok {
		err = errors.New("GetBigIntFromString err")
		Logger().Error(err.Error())
		return
	}
	return
}

func ConvertBigIntFromString(v0, v1 string) (n0 *big.Int, n1 *big.Int, err error) {
	n0 = new(big.Int)
	n0, ok := n0.SetString(v0, 10)
	if !ok {
		err = errors.New("ConvertBigIntFromString err")
		Logger().Error(err.Error())
		return
	}

	n1 = new(big.Int)
	n1, ok = n1.SetString(v1, 10)
	if !ok {
		err = errors.New("ConvertBigIntFromString err")
		Logger().Error(err.Error())
		return
	}
	return
}

func GenerateIncreaseID() (int64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		Logger().Error(fmt.Sprintf("GenerateIncreaseID err:%+v", err))
		return 0, err
	}
	// Generate a snowflake ID.
	id := node.Generate()

	return id.Int64(), nil
}

func RemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func CreateTransactionOpts(client *ethclient.Client, key *ecdsa.PrivateKey, chainId uint64, caller common.Address, amount *big.Int) (opts *bind.TransactOpts, err error) {
	nonce, err := client.PendingNonceAt(context.Background(), caller)
	if err != nil {
		errMsg := fmt.Sprintf("CreateTransactionOpts:client.PendingNonceAt err: %+v", err)
		Logger().Error(errMsg)
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		errMsg := fmt.Sprintf("CreateTransactionOpts:client.SuggestGasPrice err: %+v", err)
		Logger().Error(errMsg)
		return nil, err
	}

	srcChainID := big.NewInt(int64(chainId))
	opts, err = bind.NewKeyedTransactorWithChainID(key, srcChainID)
	if err != nil {
		errMsg := fmt.Sprintf("CreateTransactionOpts:NewKeyedTransactorWithChainID err: %+v", err)
		Logger().Error(errMsg)
		return nil, err
	}

	opts.Nonce = big.NewInt(int64(nonce))
	if amount != nil {
		opts.Value = amount // in wei
	} else {
		opts.Value = big.NewInt(0) // in wei
	}

	opts.GasLimit = uint64(0) // in units
	opts.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))

	return opts, nil
}

func TxWaitToSync(ctx context.Context, client *ethclient.Client, tx *types.Transaction) (*types.Receipt, bool, error) {
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		errMsg := fmt.Sprintf("TxWaitToSync:bind.WaitMine err: %+v", err)
		Logger().Error(errMsg)
		return nil, false, err
	}

	return receipt, receipt.Status == types.ReceiptStatusSuccessful, nil
}

func PrivateToAddress(key string) (string, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}
	addr := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return addr, nil
}

type EthCallResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

func GetEthCallPostData(to string, data string) string {
	tpl := `{"method":"eth_call","params":[{"from": null,"to":"%s","data":"%s"}, "latest"],"id":1,"jsonrpc":"2.0"}`
	return fmt.Sprintf(tpl, to, data)
}

func ReadContract(url string, postdata string) (*EthCallResult, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = " application/json"
	hc := GetHTTPClient()
	result := &EthCallResult{}
	body, err := HTTPReq("POST", url, hc, []byte(postdata), headers)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Pad left-pads s with spaces, to length n.
// If n is smaller than s, Pad is a no-op.
func Pad(s string, n int, r rune) (string, error) {
	return PadChar(s, n, r)
}

// PadChar left-pads s with the rune r, to length n.
// If n is smaller than s, PadChar is a no-op.
func PadChar(s string, n int, r rune) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("invalid length %d", n)
	}
	if len(s) > n {
		return s, nil
	}
	return strings.Repeat(string(r), n-len(s)) + s, nil
}

type Ticker struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []*TickerData `json:"data"`
}
type TickerData struct {
	InstID  string `json:"instId"`
	IdxPx   string `json:"idxPx"`
	High24H string `json:"high24h"`
	SodUtc0 string `json:"sodUtc0"`
	Open24H string `json:"open24h"`
	Low24H  string `json:"low24h"`
	SodUtc8 string `json:"sodUtc8"`
	Ts      string `json:"ts"`
}

func GetTokenPriceUSDT(okxurl string, base, quote string) (float64, error) {
	client := GetHTTPClient()
	url := okxurl + strings.ToUpper(base) + "-" + strings.ToUpper(quote)
	headers := make(map[string]string)
	headers["Content-Type"] = " application/json"
	headers["Accept"] = " application/json"
	data, err := HTTPReq("GET", url, client, nil, headers)
	if err != nil {
		return 0, err
	}
	ticker := &Ticker{}
	err = json.Unmarshal(data, ticker)
	if err != nil {
		return 0, err
	}
	if len(ticker.Data) == 0 {
		return 0, nil
	}
	price, err := strconv.ParseFloat(ticker.Data[0].IdxPx, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}
