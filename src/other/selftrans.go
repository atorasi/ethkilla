package other

import (
	"ethkilla/src/account"
	"ethkilla/src/constants"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewOtherClient(client *ethclient.Client, wallet account.Wallet) OtherClient {
	return OtherClient{
		Client:  client,
		Wallet:  wallet,
		Account: account.Account{},
	}
}

func (app OtherClient) SelfTranaction() (*types.Receipt, error) {
	balance, _ := app.Account.NativeBalance(app.Client, app.Wallet)
	value := app.Account.RandomValue(balance, constants.SETTINGS.SelfTransPercentMin, constants.SETTINGS.SelfTransPercentMax)

	txHash, err := app.Account.SendTransaction(app.Client, app.Wallet, app.Wallet.PublicKey, nil, value, "ETHEREUM")
	if err != nil {
		return nil, err
	}

	return txHash, nil
}
