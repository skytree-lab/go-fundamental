package solana

import (
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
)

func Test_sol_transaction(t *testing.T) {
	endpoint := rpc.TestNet_RPC
	wallet := "7HZaCWazgTuuFuajxaaxGYbGnyVKwxvsJKue1W4Nvyro"
	bs, err := GetBalances(endpoint, wallet)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bs)
}

func Test_sol_metadata(t *testing.T) {
	endpoint := rpc.TestNet_RPC
	mint := "CpMah17kQEL2wqyMKt3mZBdTnZbkbfx4nqmQMFDP5vwp"
	// mint := "5HzzumbGepQduUNY2exqamvA9XY7iB86kCFPyrYLZtnb"
	bs, err := GetMetaData(endpoint, mint)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bs)
}
