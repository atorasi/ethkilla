package run

import (
	"ethkilla/src/account"
	"ethkilla/src/bridge"
	"ethkilla/src/other"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func EtherActions(index int, wallet account.Wallet, module string, etherClient *ethclient.Client) (string, error) {
	if module == "bungee" {
		log.Printf("Acc.%d | Preparing to Bungee Refuel", index)
		client := bridge.NewBridgeApp(etherClient, wallet)
		if _, err := client.BungeeRefuel(); err != nil {
			return "Bungee", err
		}
	} else {
		log.Printf("Acc.%d | Preparing to send native tokens", index)
		client := other.NewOtherClient(etherClient, wallet)
		if _, err := client.SelfTranaction(); err != nil {
			return "Bungee", err
		}
	}

	return "", nil
}
