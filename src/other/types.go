package other

import (
	"ethkilla/src/account"

	"github.com/ethereum/go-ethereum/ethclient"
)

type OtherClient struct {
	Wallet  account.Wallet
	Client  *ethclient.Client
	Account account.Account
}
