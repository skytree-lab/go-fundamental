package solana

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

func Test_TransferJitoTransaction(t *testing.T) {
	urls := []string{"https://api.mainnet-beta.solana.com"}
	jitourls := []string{
		"https://ny.mainnet.block-engine.jito.wtf/api/v1/transactions",
	}
	accTo := solana.MustPublicKeyFromBase58("6huu25nWzFtBWPMQmWRzKLD4Wtfq11SSjZTU6oitLqdz")
	accFrom := solana.MustPrivateKeyFromBase58("")

	var signers []solana.PrivateKey
	signers = append(signers, accFrom)

	var instuctions []solana.Instruction
	instuctions = append(instuctions, system.NewTransferInstruction(5000000, accFrom.PublicKey(), accTo).Build())

	sig, err := SendJitoTransaction(urls, jitourls, instuctions, uint64(1000), uint64(1000), signers)
	fmt.Println(err)
	fmt.Println(sig)
}
