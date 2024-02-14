package deposit

import (
	"ethkilla/src/account"
)

type DepositApp struct {
	account account.Account
	wallet  account.Wallet
}

type WithdrawalParams struct {
	Currency    string  `json:"ccy"`
	Amount      float64 `json:"amt"`
	Destination string  `json:"dest"`
	ToAddress   string  `json:"toAddr"`
	Fee         string  `json:"fee"`
	Chain       string  `json:"chain"`
}
