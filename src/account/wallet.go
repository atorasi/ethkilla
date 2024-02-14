package account

import (
	"ethkilla/utils"

	"github.com/ethereum/go-ethereum/crypto"
)

func SliceOfAccs() ([]Wallet, error) {
	keys, err := utils.ReadFile(`..\data\pKeys.txt`)
	if err != nil {
		return nil, err
	}

	accounts, err := Wallets(keys)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func Wallets(pks []string) ([]Wallet, error) {
	var accounts []Wallet
	for index, pk := range pks {
		account := newAccount(pk, index+1)

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func newAccount(private string, index int) Wallet {
	if private[0:2] == "0x" {
		private = private[2:]
	}
	privateKey, _ := crypto.HexToECDSA(private)

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)

	return Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Index:      index,
	}
}
