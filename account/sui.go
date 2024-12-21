package account

import (
	"fmt"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/tyler-smith/go-bip39"
)

type SuiAccount struct {
	Address    string `json:"address"`
	Mnemonic   string `json:"mnemonic"`
	PrivateKey string `json:"private"`
	PublicKey  string `json:"public"`
}

func (sa *SuiAccount) String() string {
	return fmt.Sprintf("%s:%s", sa.Address, sa.Mnemonic)
}

func NewSuiAccount() *SuiAccount {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("[+]Mnemonic:", mnemonic)
	acc, _ := account.NewAccountWithMnemonic(mnemonic)
	fmt.Printf("[+]privateKey = %x\n", acc.KeyPair.PrivateKey()[:32])
	fmt.Printf("[+] publicKey = %x\n", acc.KeyPair.PublicKey())
	fmt.Printf("[+]   address = %v\n", acc.Address)
	wallet := &SuiAccount{}
	wallet.Mnemonic = mnemonic
	wallet.Address = acc.Address
	wallet.PrivateKey = fmt.Sprintf("%x", acc.KeyPair.PrivateKey()[:32])
	wallet.PublicKey = fmt.Sprintf("%x", acc.KeyPair.PublicKey())

	return wallet
}
