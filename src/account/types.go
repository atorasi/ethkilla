package account

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
)

type Account struct{}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  common.Address
	Index      int
}
