package solana

import (
	"encoding/binary"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	jitorpc "github.com/jito-labs/jito-go-rpc"
	"github.com/mr-tron/base58"
)

func createSetComputeUnitPriceInstruction(microLamports uint64) solana.Instruction {
	data := make([]byte, 9)
	data[0] = 3 // Instruction index for SetComputeUnitPrice
	binary.LittleEndian.PutUint64(data[1:], microLamports)
	return solana.NewInstruction(
		solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111"),
		solana.AccountMetaSlice{},
		data,
	)
}

func BuildJitoTransactionAndSend(signers []solana.PrivateKey, instructions []solana.Instruction, jitotip uint64, priority uint64, jitourl string, urls []string) (tx string, err error) {
	jitoClient := jitorpc.NewJitoJsonRpcClient(jitourl, "")
	randomTipAccount, err := jitoClient.GetRandomTipAccount()
	if err != nil {
		return
	}
	jitoTipAccount, err := solana.PublicKeyFromBase58(randomTipAccount.Address)
	if err != nil {
		return
	}
	tipIns := system.NewTransferInstruction(
		jitotip,
		signers[0].PublicKey(),
		jitoTipAccount,
	).Build()

	var targetInstructions []solana.Instruction
	targetInstructions = append(targetInstructions, createSetComputeUnitPriceInstruction(priority))
	targetInstructions = append(targetInstructions, instructions...)
	targetInstructions = append(targetInstructions, tipIns)
	var latestBlockhash *LastestBlockHashResult
	for _, url := range urls {
		latestBlockhash, err = GetLatestBlockhash(url)
		if err != nil {
			continue
		}
		if latestBlockhash != nil {
			break
		}
	}
	if latestBlockhash == nil {
		return
	}
	transaction, err := solana.NewTransaction(
		targetInstructions,
		latestBlockhash.Value.Blockhash,
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		return
	}
	_, err = transaction.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		for _, payer := range signers {
			if payer.PublicKey().Equals(key) {
				return &payer
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	serializedTx, err := transaction.MarshalBinary()
	if err != nil {
		return
	}
	base58EncodedTx := base58.Encode(serializedTx)
	txnRequest := []string{base58EncodedTx}
	result, err := jitoClient.SendTxn(txnRequest, false)
	if err != nil {
		return
	}
	tx = strings.Trim(string(result), "\"")
	return
}

func SendJitoTransaction(urls []string, jitorpcurls []string, instructions []solana.Instruction, tip uint64, priority uint64, signers []solana.PrivateKey) (tx string, err error) {
	for _, jitourl := range jitorpcurls {
		tx, err = BuildJitoTransactionAndSend(signers, instructions, tip, priority, jitourl, urls)
		if err != nil {
			continue
		}
	}
	return
}
