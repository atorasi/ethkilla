package run

import (
	"ethkilla/src/account"
	"ethkilla/src/bridge"
	"ethkilla/src/constants"
	"ethkilla/src/deposit"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func SideActions(index int, wallet account.Wallet, client, ethClient *ethclient.Client) (string, error) {
	sideChain := constants.SETTINGS.SideChain

	var module string
	if constants.SETTINGS.NeedOkx {
		module, err := okxWithdrawal(index, sideChain, client, wallet)
		if err != nil {
			return module, err
		}
	}

	module, err := relayDeposit(index, sideChain, client, ethClient, wallet)
	if err != nil {
		return module, err
	}

	return "", nil
}

func okxWithdrawal(index int, sideChain string, client *ethclient.Client, wallet account.Wallet) (string, error) {
	balanceSideStart, err := account.Account.NativeBalance(account.Account{}, client, wallet)
	if err != nil {
		return "Get balance", err
	}

	log.Printf("Acc.%d | Preparing to Okx withdrawal", index)
	depositClient := deposit.NewDepositApp(wallet)
	if _, err := depositClient.OkxWithdraw(sideChain); err != nil {
		return "OKX", err
	}
	log.Printf("Acc.%d | Succesfully withdrew from OKX, waiting for funds", index)

	for {
		newBalance, err := account.Account.NativeBalance(account.Account{}, client, wallet)
		if err != nil {
			return "Get balance", err
		}
		if newBalance != balanceSideStart {
			log.Printf("Acc.%d | Funds deposited", index)
			break
		}
		time.Sleep(time.Duration(30) * time.Second)
		log.Printf("Acc.%d | Didnt get funds yet, sleep 30 seconds.", index)

	}

	return "", nil
}

func relayDeposit(index int, sideChain string, client, ethclient *ethclient.Client, wallet account.Wallet) (string, error) {
	balanceEthStart, err := account.Account.NativeBalance(account.Account{}, ethclient, wallet)
	if err != nil {
		return "Get balance", err
	}

	log.Printf("Acc.%d | Preparing to Relay App", index)
	bridgeClient := bridge.NewBridgeApp(client, wallet)
	if _, err := bridgeClient.RelayBridge(sideChain); err != nil {
		return "Relay", err
	}

	for newBalance := balanceEthStart; newBalance == balanceEthStart; {
		newBalance, err = account.Account.NativeBalance(account.Account{}, ethclient, wallet)
		if err != nil {
			return "Get balance", err
		}
		log.Printf("Acc.%d | Didnt get funds yet, sleep 30 seconds.", index)
		time.Sleep(time.Duration(30) * time.Second)
	}

	return "", nil
}
