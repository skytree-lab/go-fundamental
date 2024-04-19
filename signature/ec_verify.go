package signature

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type EIP191 struct {
	msg       string
	signature string
	address   string
}

func hasValidLastByte(sig []byte) bool {
	return sig[64] == 0 || sig[64] == 1
}

func hasMatchingAddress(knownAddress string, recoveredAddress string) bool {
	return strings.EqualFold(knownAddress, recoveredAddress)
}

func signEIP191(message string) common.Hash {
	msg := []byte(message)
	formattedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
	return crypto.Keccak256Hash([]byte(formattedMsg))
}

func VerifySignature(addr string, data []byte, signature []byte) (bool, error) {
	eipChallenge := &EIP191{
		msg:       string(data),
		signature: string(signature),
		address:   addr,
	}
	decodedSig, err := hexutil.Decode(eipChallenge.signature)
	if err != nil {
		return false, err
	}

	if decodedSig[64] < 27 {
		if !hasValidLastByte(decodedSig) {
			err := errors.New("invalid last byte")
			return false, err
		}
	} else {
		decodedSig[64] -= 27 // shift byte?
	}

	hash := signEIP191(eipChallenge.msg)

	recoveredPublicKey, err := crypto.Ecrecover(hash.Bytes(), decodedSig)
	if err != nil {
		return false, err
	}

	secp256k1RecoveredPublicKey, err := crypto.UnmarshalPubkey(recoveredPublicKey)
	if err != nil {
		return false, err
	}

	recoveredAddress := crypto.PubkeyToAddress(*secp256k1RecoveredPublicKey).Hex()

	if hasMatchingAddress(eipChallenge.address, recoveredAddress) {
		return true, nil
	} else {
		errMsg := fmt.Sprintf("Recovered address %s does not match %s\n", recoveredAddress, eipChallenge.address)
		err := errors.New(errMsg)
		return false, err
	}
}
