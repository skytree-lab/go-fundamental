package signature

import (
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func Test_rc_recovery_addr2(t *testing.T) {
	privateKeyString := ""
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		log.Fatal(err)
	}
	// Sign message
	signature, err := PersonalSign("test", privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signature)
	addr := ""
	verified, err := VerifySignature(addr, []byte("test"), []byte(signature))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(verified)
}
