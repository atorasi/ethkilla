package bridge

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"ethkilla/src/account"
	"ethkilla/src/constants"
	"io"
	"math/big"
	"strconv"

	http "github.com/Danny-Dasilva/fhttp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewBridgeApp(client *ethclient.Client, wallet account.Wallet) Bridge {
	return Bridge{
		Client:  client,
		Wallet:  wallet,
		Account: account.Account{},
	}
}

func (app Bridge) RelayBridge(chainName string) (*types.Receipt, error) {
	nativeBalance, err := app.Account.NativeBalance(app.Client, app.Wallet)
	if err != nil {
		return nil, err
	}
	value := app.Account.RandomValue(nativeBalance, constants.SETTINGS.RelayPercentMin, constants.SETTINGS.RelayPercentMax)

	txData, err := app.parceCalldata(chainName, value)
	if err != nil {
		return nil, err
	}

	toAddr := common.HexToAddress(txData.Steps[0].Items[0].Data.To)
	calldata, _ := hex.DecodeString(txData.Steps[0].Items[0].Data.Data[2:])

	txHash, err := app.Account.SendTransaction(app.Client, app.Wallet, toAddr, calldata, value, chainName)
	if err != nil {
		return nil, err
	}

	return txHash, nil
}

func (app Bridge) parceCalldata(chainName string, value *big.Int) (RelayResponse, error) {
	chainID, _ := strconv.Atoi(constants.CHAINS[chainName]["CHAIN_ID"])

	payload, err := json.Marshal(&RelayPayload{
		User:   app.Wallet.PublicKey.String(),
		Sourse: "relay.link",
		Txs: []struct {
			To    string "json:\"to\""
			Value string "json:\"value\""
			Data  string "json:\"data\""
		}{
			{
				To:    app.Wallet.PublicKey.String(),
				Value: value.String(),
				Data:  "0x",
			},
		},
		OriginChainID:      chainID,
		DestinationChainID: 1,
	})
	if err != nil {
		return RelayResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://api.relay.link/execute/call", bytes.NewBuffer(payload))
	if err != nil {
		return RelayResponse{}, err
	}
	request.Header = http.Header{
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
		"Host":         {"api.relay.link"},
		"Origin":       {"https://www.relay.link"},
	}
	response, err := http.DefaultClient.Do(request)
	if response.StatusCode != 200 || err != nil {
		return RelayResponse{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return RelayResponse{}, err
	}

	var requestData RelayResponse
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return RelayResponse{}, err
	}

	return requestData, nil
}
