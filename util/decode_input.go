package util

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func DecodeInputData(abiStr string, inputData string) map[string]interface{} {

	// load contract ABI
	abi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil
	}

	// decode txInput method signature
	decodedSig, err := hex.DecodeString(inputData[2:10])
	if err != nil {
		return nil
	}

	// recover Method from signature and ABI
	method, err := abi.MethodById(decodedSig)
	if err != nil {
		return nil
	}

	// decode txInput Payload
	decodedData, err := hex.DecodeString(inputData[10:])
	if err != nil {
		return nil
	}

	// unpack method inputs
	inputMap := make(map[string]interface{}, 0)
	err = method.Inputs.UnpackIntoMap(inputMap, decodedData)
	if err != nil {
		return nil
	}
	return inputMap

}
