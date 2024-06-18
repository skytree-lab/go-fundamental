package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	sol "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/tokenregistry"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	solws "github.com/gagliardetto/solana-go/rpc/ws"
	streamingfastSol "github.com/streamingfast/solana-go"
	associatedtokenaccount "github.com/streamingfast/solana-go/programs/associated-token-account"
	"github.com/streamingfast/solana-go/programs/token"
	streamingfastRpc "github.com/streamingfast/solana-go/rpc"
	streamingfastWs "github.com/streamingfast/solana-go/rpc/ws"
)

type Balance struct {
	Parsed  Parsed `json:"parsed"`
	Program string `json:"program"`
}
type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
type Info struct {
	IsNative    bool        `json:"isNative"`
	Mint        string      `json:"mint"`
	Owner       string      `json:"owner"`
	State       string      `json:"state"`
	TokenAmount TokenAmount `json:"tokenAmount"`
}
type Parsed struct {
	Info Info   `json:"info"`
	Type string `json:"type"`
}

func ClaimAirDropFromTestNet(pubKey string) error {
	acc := sol.MustPublicKeyFromBase58(pubKey)
	client := rpc.New(rpc.TestNet_RPC)
	sig, err := client.RequestAirdrop(context.TODO(), acc, sol.LAMPORTS_PER_SOL*1, rpc.CommitmentFinalized)
	if err != nil {
		return err
	}
	spew.Dump(sig)
	return nil
}

func TransferSol(url, wsurl string, keyFrom string, to string, amount uint64) error {
	accountFrom, err := sol.PrivateKeyFromBase58(keyFrom)
	if err != nil {
		return err
	}
	accountTo := sol.MustPublicKeyFromBase58(to)
	rpcClient := rpc.New(url)
	wsClient, err := solws.Connect(context.Background(), wsurl)
	if err != nil {
		return err
	}

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return err
	}

	instruction := system.NewTransferInstruction(amount, accountFrom.PublicKey(), accountTo).Build()
	instructions := []sol.Instruction{
		instruction,
	}
	tx, err := sol.NewTransaction(instructions, recent.Value.Blockhash, sol.TransactionPayer(accountFrom.PublicKey()))
	if err != nil {
		return err
	}

	getter := func(key sol.PublicKey) *sol.PrivateKey {
		if accountFrom.PublicKey().Equals(key) {
			return &accountFrom
		}
		return nil
	}

	_, err = tx.Sign(getter)
	if err != nil {
		return err
	}

	sig, err := confirm.SendAndConfirmTransaction(context.TODO(), rpcClient, wsClient, tx)
	if err != nil {
		return err
	}
	spew.Dump(sig)

	return nil
}

func GetMetaData(url, mintAddress string) (*tokenregistry.TokenMeta, error) {
	client := rpc.New(url)
	resp, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		sol.TokenMetadataProgramID,
		&rpc.GetProgramAccountsOpts{
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 32,
						Bytes:  []byte(mintAddress),
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp empty... cannot find account")
	}

	for idx, keyedAcct := range resp {
		acct := keyedAcct.Account
		var t *tokenregistry.TokenMeta
		decoder := bin.NewBinDecoder(acct.Data.GetBinary())
		err := decoder.Decode(&t)
		if err != nil {
			continue
		}
		msg := fmt.Sprintf("idx:%d mint:%s registration:%s symbol:%s name:%s", idx, t.MintAddress.String(), t.RegistrationAuthority.String(), t.Symbol.String(), t.Name.String())
		fmt.Println(msg)
		if t.MintAddress.String() == mintAddress || t.RegistrationAuthority.String() == mintAddress {
			return t, nil
		}
	}

	return nil, errors.New("not found")
}

func GetBalances(url string, wallet string) ([]*Balance, error) {
	client := rpc.New(url)
	pubKey := sol.MustPublicKeyFromBase58(wallet)
	var balances []*Balance
	out, err := client.GetTokenAccountsByOwner(
		context.TODO(),
		pubKey,
		&rpc.GetTokenAccountsConfig{
			ProgramId: sol.TokenProgramID.ToPointer(),
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: sol.EncodingJSONParsed,
		},
	)

	if err != nil {
		return balances, err
	}

	for _, rawAccount := range out.Value {
		data, err := rawAccount.Account.Data.MarshalJSON()
		if err != nil {
			return balances, err
		}

		balance := &Balance{}
		json.Unmarshal(data, balance)
		balances = append(balances, balance)
	}

	return balances, nil
}

func TransferToken(url, wsurl string, keyFrom string, to string, amount uint64, mint string) (string, error) {
	rpcClient := streamingfastRpc.NewClient(url)
	wsClient := streamingfastWs.NewClient(wsurl, false)

	sender := &streamingfastSol.Account{
		PrivateKey: streamingfastSol.MustPrivateKeyFromBase58(keyFrom),
	}

	mintPub := streamingfastSol.MustPublicKeyFromBase58(mint)
	toPub := streamingfastSol.MustPublicKeyFromBase58(to)
	senderAta := associatedtokenaccount.MustGetAssociatedTokenAddress(mintPub, token.PROGRAM_ID, sender.PublicKey())
	_, tx, err := token.TransferToken(context.TODO(), rpcClient, wsClient, amount, senderAta, mintPub, toPub, sender)
	return tx, err
}
