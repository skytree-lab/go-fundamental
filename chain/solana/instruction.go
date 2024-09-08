package solana

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/skytree-lab/go-fundamental/util"
)

func ParseTransferSOLInstructionParam(out *rpc.GetTransactionResult, urls []string) (params []*TransferSOLInstructionParam, succeed bool, err error) {
	if out == nil || out.Transaction == nil || out.Meta == nil {
		util.Logger().Error("out GetTransaction nil")
		return
	}

	if out.Meta.Err == nil {
		succeed = true
	} else {
		succeed = false
	}

	tx, err := out.Transaction.GetTransaction()
	if err != nil {
		util.Logger().Error(fmt.Sprintf("out GetTransaction err:%v", err))
		return
	}

	err = ProcessTransactionWithAddressLookups(tx, urls)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("parse err:%v", err))
		return
	}

	accoountKeys, err := tx.Message.GetAllKeys()
	if err != nil {
		util.Logger().Error(fmt.Sprintf("parse err:%v", err))
		return
	}

	for _, instruction := range tx.Message.Instructions {
		datas := []byte(instruction.Data)
		if len(datas) < 12 {
			continue
		}

		instype := bin.LE.Uint32(datas[0:4])
		if instype != system.Instruction_Transfer {
			continue
		}

		if len(accoountKeys) <= int(instruction.ProgramIDIndex) {
			continue
		}

		if accoountKeys[instruction.ProgramIDIndex] != solana.SystemProgramID {
			continue
		}

		sourceIdx := instruction.Accounts[0]
		destIdx := instruction.Accounts[1]

		param := &TransferSOLInstructionParam{
			Source:      accoountKeys[sourceIdx].String(),
			Destination: accoountKeys[destIdx].String(),
			Amount:      bin.LE.Uint64(datas[4:12]),
		}

		params = append(params, param)
	}

	return
}

func ParseRaydiumSwapInstructionParam(out *rpc.GetTransactionResult, urls []string) (param *RaydiumSwapInstructionParam, succeed bool, err error) {
	if out == nil || out.Transaction == nil || out.Meta == nil {
		util.Logger().Error("out GetTransaction nil")
		return
	}

	if out.Meta.Err == nil {
		succeed = true
	} else {
		succeed = false
	}

	tx, err := out.Transaction.GetTransaction()
	if err != nil {
		util.Logger().Error(fmt.Sprintf("out GetTransaction err:%v", err))
		return
	}

	err = ProcessTransactionWithAddressLookups(tx, urls)
	if err != nil {
		util.Logger().Error(fmt.Sprintf("parse err:%v", err))
		return
	}

	accoountKeys, err := tx.Message.GetAllKeys()
	if err != nil {
		util.Logger().Error(fmt.Sprintf("parse err:%v", err))
		return
	}

	tempParam := &RaydiumSwapInstructionParam{}
	for idx, instruction := range tx.Message.Instructions {
		datas := []byte(instruction.Data)
		if len(datas) < 17 {
			continue
		}

		instype := uint8(datas[0])
		if instype != 9 {
			continue
		}

		if len(accoountKeys) <= int(instruction.ProgramIDIndex) {
			continue
		}

		if accoountKeys[instruction.ProgramIDIndex] != RaydiumLiquidityPoolv4ProgramID {
			continue
		}

		userSourceIdx := instruction.Accounts[15]
		userDestinationIdx := instruction.Accounts[16]

		for _, inner := range out.Meta.InnerInstructions {
			if inner.Index != uint16(idx) {
				continue
			}

			for _, innerInstruction := range inner.Instructions {
				if len(accoountKeys) <= int(innerInstruction.ProgramIDIndex) {
					continue
				}

				if accoountKeys[innerInstruction.ProgramIDIndex] != solana.TokenProgramID {
					continue
				}

				if len(innerInstruction.Accounts) != 3 {
					continue
				}

				dataBytes := []byte(innerInstruction.Data)
				if len(dataBytes) != 9 {
					continue
				}

				instype := uint8(dataBytes[0])
				if instype != token.Instruction_Transfer {
					continue
				}

				if accoountKeys[innerInstruction.ProgramIDIndex] != solana.TokenProgramID {
					continue
				}

				swapdata := &RaydiumSwapInnerInstructionData{
					Source:       accoountKeys[innerInstruction.Accounts[0]].String(),
					Destionation: accoountKeys[innerInstruction.Accounts[1]].String(),
					Owner:        accoountKeys[innerInstruction.Accounts[2]].String(),
					Amount:       bin.LE.Uint64(dataBytes[1:9]),
				}
				if innerInstruction.Accounts[0] == userSourceIdx {
					tempParam.SwapIn = swapdata
				} else if innerInstruction.Accounts[1] == userDestinationIdx {
					tempParam.SwapOut = swapdata
				}
			}
		}
	}

	if tempParam.SwapIn != nil && tempParam.SwapOut != nil {
		param = tempParam
	}

	return
}
