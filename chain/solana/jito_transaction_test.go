package solana

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

func Test_TransferJitoTransaction(t *testing.T) {
	urls := []string{
		"https://api.mainnet-beta.solana.com",
	}

	jitourls := []string{
		"https://mainnet.block-engine.jito.wtf/api/v1",
	}
	accTo := solana.MustPublicKeyFromBase58("6huu25nWzFtBWPMQmWRzKLD4Wtfq11SSjZTU6oitLqdz")
	accFrom := solana.MustPrivateKeyFromBase58("")

	var signers []solana.PrivateKey
	signers = append(signers, accFrom)

	var instuctions []solana.Instruction
	instuctions = append(instuctions, system.NewTransferInstruction(1000, accFrom.PublicKey(), accTo).Build())

	sig, err := SendJitoTransaction(urls, jitourls, instuctions, uint64(10000), signers)
	fmt.Println(err)
	fmt.Println(sig)
}
