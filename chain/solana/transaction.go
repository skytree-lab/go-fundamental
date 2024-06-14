package solana

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	solws "github.com/gagliardetto/solana-go/rpc/ws"
)

func ClaimAirDropFromTestNet(pubKey string) error {
	acc := solana.MustPublicKeyFromBase58(pubKey)
	client := rpc.New(rpc.TestNet_RPC)
	sig, err := client.RequestAirdrop(context.TODO(), acc, solana.LAMPORTS_PER_SOL*1, rpc.CommitmentFinalized)
	if err != nil {
		return err
	}
	spew.Dump(sig)
	return nil
}

func Transfer(url, wsurl string, keyFrom string, to string, amount uint64) error {
	accountFrom, err := solana.PrivateKeyFromBase58(keyFrom)
	if err != nil {
		return err
	}
	accountTo := solana.MustPublicKeyFromBase58(to)
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
	instructions := []solana.Instruction{
		instruction,
	}
	tx, err := solana.NewTransaction(instructions, recent.Value.Blockhash, solana.TransactionPayer(accountFrom.PublicKey()))
	if err != nil {
		return err
	}

	getter := func(key solana.PublicKey) *solana.PrivateKey {
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
