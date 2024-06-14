package account

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type SolanaAccount struct {
	PrivateKey string
	Address    string
}

func (sa *SolanaAccount) String() string {
	return fmt.Sprintf("%s:%s", sa.Address, sa.PrivateKey)
}

func NewSolnanaAccount() *SolanaAccount {
	a := solana.NewWallet()
	sa := &SolanaAccount{
		PrivateKey: a.PrivateKey.String(),
		Address:    a.PublicKey().String(),
	}

	return sa
}
