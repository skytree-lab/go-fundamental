package util

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
)

func Test_decode_input(t *testing.T) {
	type swapInputs struct {
		PoolId        uint32
		DstChainId    uint32
		Amount        *big.Int
		Recipient     [32]byte
		RefundAddress common.Address
		ConsumerData  []byte
	}
	abiStr := "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"factory_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CALLED_ON_NATIVE\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"INVALID_AMOUNT\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NOT_POOL_OWNER\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"collector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PoolFeeCollected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"newFeeRatio\",\"type\":\"uint16\"}],\"name\":\"PoolFeeRatioUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"PoolLimitUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"isIncrease\",\"type\":\"bool\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lpAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PoolLiquidityModified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"dstChainId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lpAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RemoveRemoteLiquidity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"dstChainId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Swapped\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NATIVE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REMOTE_REMOVE_LIQUIDITY\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REMOTE_SWAP_OUT\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"collectLPFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"decreaseLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"refundAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"dstChainId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"recipient\",\"type\":\"bytes32\"}],\"name\":\"decreaseLiquidityRemote\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"dstChainId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"recipient\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"consumerData\",\"type\":\"bytes\"}],\"name\":\"estimateBNFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"increaseLiquidity\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"feeRatio\",\"type\":\"uint16\"}],\"name\":\"setFeeRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"setSwapLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"poolId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"dstChainId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"recipient\",\"type\":\"bytes32\"},{\"internalType\":\"addresspayable\",\"name\":\"refundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"consumerData\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"crossId\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

	// example of transaction input data
	txInput := "0xed866fb900000000000000000000000000000000000000000000000000000000000000070000000000000000000000000000000000000000000000000000000000066eed0000000000000000000000000000000000000000000000000000000005f5e1000000000000000000000000003e732199984e53fec585eedd600be21aa435ceb60000000000000000000000003e732199984e53fec585eedd600be21aa435ceb600000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000"

	var data swapInputs
	inputMap := DecodeInputData(abiStr, txInput)

	err := mapstructure.Decode(inputMap, &data)
	if err != nil || len(inputMap) == 0 {
		t.Errorf("DecodeInputData (%s %s) was incorrect", abiStr, txInput)
	}

	if data.PoolId != 7 {
		t.Error("DecodeInputData mapstructure was incorrect")
	}

}
