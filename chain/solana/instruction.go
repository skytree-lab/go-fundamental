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

func ParseTransferSOLInstructionParam(out *rpc.GetTransactionResult) (params []*TransferSOLInstructionParam, err error) {
	if out == nil || out.Transaction == nil {
		util.Logger().Error("out GetTransaction nil")
		return
	}

	tx, err := out.Transaction.GetTransaction()
	if err != nil {
		util.Logger().Error(fmt.Sprintf("out GetTransaction err:%v", err))
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

		if len(tx.Message.AccountKeys) <= int(instruction.ProgramIDIndex) {
			continue
		}

		if tx.Message.AccountKeys[instruction.ProgramIDIndex] != solana.SystemProgramID {
			continue
		}

		sourceIdx := instruction.Accounts[0]
		destIdx := instruction.Accounts[1]

		param := &TransferSOLInstructionParam{
			Source:      tx.Message.AccountKeys[sourceIdx].String(),
			Destination: tx.Message.AccountKeys[destIdx].String(),
			Amount:      bin.LE.Uint64(datas[4:12]),
		}

		params = append(params, param)
	}

	return
}

func ParseRaydiumSwapInstructionParam(out *rpc.GetTransactionResult) (param *RaydiumSwapInstructionParam, err error) {
	if out == nil || out.Transaction == nil {
		util.Logger().Error("out GetTransaction nil")
		return
	}

	tx, err := out.Transaction.GetTransaction()
	if err != nil {
		util.Logger().Error(fmt.Sprintf("out GetTransaction err:%v", err))
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

		if len(tx.Message.AccountKeys) <= int(instruction.ProgramIDIndex) {
			continue
		}

		if tx.Message.AccountKeys[instruction.ProgramIDIndex] != RaydiumLiquidityPoolv4ProgramID {
			continue
		}

		userSourceIdx := instruction.Accounts[15]
		userDestinationIdx := instruction.Accounts[16]

		for _, inner := range out.Meta.InnerInstructions {
			if inner.Index != uint16(idx) {
				continue
			}

			for _, innerInstruction := range inner.Instructions {
				if len(tx.Message.AccountKeys) <= int(innerInstruction.ProgramIDIndex) {
					continue
				}

				if tx.Message.AccountKeys[innerInstruction.ProgramIDIndex] != solana.TokenProgramID {
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

				if tx.Message.AccountKeys[innerInstruction.ProgramIDIndex] != solana.TokenProgramID {
					continue
				}

				swapdata := &RaydiumSwapInnerInstructionData{
					Source:       tx.Message.AccountKeys[innerInstruction.Accounts[0]].String(),
					Destionation: tx.Message.AccountKeys[innerInstruction.Accounts[1]].String(),
					Owner:        tx.Message.AccountKeys[innerInstruction.Accounts[2]].String(),
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
