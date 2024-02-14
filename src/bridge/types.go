package bridge

import (
	"ethkilla/src/account"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Bridge struct {
	Client  *ethclient.Client
	Wallet  account.Wallet
	Account account.Account
}

type RelayPayload struct {
	DestinationChainID int    `json:"destinationChainId"`
	OriginChainID      int    `json:"originChainId"`
	Sourse             string `json:"source"`
	User               string `json:"user"`
	Txs                []struct {
		To    string `json:"to"`
		Value string `json:"value"`
		Data  string `json:"data"`
	} `json:"txs"`
}

type RelayResponse struct {
	Steps []struct {
		ID          string `json:"id"`
		Action      string `json:"action"`
		Description string `json:"description"`
		Kind        string `json:"kind"`
		Items       []struct {
			Status string `json:"status"`
			Data   struct {
				From    string `json:"from"`
				To      string `json:"to"`
				Data    string `json:"data"`
				Value   string `json:"value"`
				ChainID int    `json:"chainId"`
			} `json:"data"`
			Check struct {
				Endpoint string `json:"endpoint"`
				Method   string `json:"method"`
			} `json:"check"`
		} `json:"items"`
	} `json:"steps"`
	Fees struct {
		Gas            string `json:"gas"`
		Relayer        string `json:"relayer"`
		RelayerGas     string `json:"relayerGas"`
		RelayerService string `json:"relayerService"`
	} `json:"fees"`
	Balances struct {
		UserBalance     string `json:"userBalance"`
		RequiredToSolve string `json:"requiredToSolve"`
	} `json:"balances"`
}
