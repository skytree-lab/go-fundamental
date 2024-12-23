package solana

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	jitorpc "github.com/jito-labs/jito-go-rpc"
	"github.com/mr-tron/base58"
	"github.com/skytree-lab/go-fundamental/util"
)

func checkBundleStatus(jitoClient *jitorpc.JitoJsonRpcClient, bundleId string) (sig string, err error) {
	maxAttempts := 60
	pollInterval := 5 * time.Second

	var statusResponse *jitorpc.BundleStatusResponse
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(pollInterval)

		statusResponse, err = jitoClient.GetBundleStatuses([]string{bundleId})
		if err != nil {
			util.Logger().Error(fmt.Sprintf("Attempt %d: Failed to get bundle status: %v", attempt, err))
			continue
		}
		if len(statusResponse.Value) == 0 {
			util.Logger().Error(fmt.Sprintf("Attempt %d: No bundle status available", attempt))
			continue
		}

		bundleStatus := statusResponse.Value[0]
		switch bundleStatus.ConfirmationStatus {
		case "processed":
			util.Logger().Info("Bundle has been processed by the cluster. Continuing to poll...")
		case "confirmed":
			util.Logger().Info("Bundle has been confirmed by the cluster. Continuing to poll...")
		case "finalized":
			util.Logger().Info(fmt.Sprintf("Bundle has been finalized by the cluster in slot %d.\n", bundleStatus.Slot))
			if bundleStatus.Err.Ok == nil {
				util.Logger().Info("Bundle executed successfully.")
				util.Logger().Info("Transaction URLs:")
				for _, txID := range bundleStatus.Transactions {
					sig = txID
					util.Logger().Info(fmt.Sprintf("- https://solscan.io/tx/%s\n", txID))
					return
				}
			} else {
				util.Logger().Info(fmt.Sprintf("Bundle execution failed with error: %v\n", bundleStatus.Err.Ok))
			}
			return
		default:
			util.Logger().Info(fmt.Sprintf("Unexpected status: %s. Please check the bundle manually.\n", bundleStatus.ConfirmationStatus))
			return
		}
	}
	return
}

func BuildJitoTransactionAndSend(urls []string, jitourl string, signers []solana.PrivateKey, instructions []solana.Instruction, tip uint64) (tx string, err error) {
	jitoClient := jitorpc.NewJitoJsonRpcClient(jitourl, "")
	randomTipAccount, err := jitoClient.GetRandomTipAccount()
	if err != nil {
		return
	}
	tipAccount, err := solana.PublicKeyFromBase58(randomTipAccount.Address)
	if err != nil {
		return
	}
	tipIns := system.NewTransferInstruction(
		tip,
		signers[0].PublicKey(),
		tipAccount,
	).Build()

	var targetInstructions []solana.Instruction
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
	var txSignature string
	encodedTx := base58.Encode(serializedTx)
	txnRequest := []string{encodedTx}
	bundleRequest := [][]string{txnRequest}
	result, err := jitoClient.SendBundle(bundleRequest)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("Failed to send bundle: %v", err))
		return
	}
	if err := json.Unmarshal(result, &txSignature); err != nil {
		log.Fatalf("Failed to unmarshal bundle ID: %v", err)
	}
	tx, err = checkBundleStatus(jitoClient, txSignature)
	return
}

func SendJitoTransaction(urls []string, jitorpcbases []string, instructions []solana.Instruction, tip uint64, signers []solana.PrivateKey) (tx string, err error) {
	for _, jitobase := range jitorpcbases {
		tx, err = BuildJitoTransactionAndSend(urls, jitobase, signers, instructions, tip)
		if err != nil {
			continue
		}
	}
	return
}
