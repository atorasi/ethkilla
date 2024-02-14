package bridge

import (
	"ethkilla/src/constants"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (app Bridge) BungeeRefuel() (*types.Receipt, error) {
	contractAddr := common.HexToAddress(constants.CHAINS["ETHEREUM"]["BUNGEECONTRACT"])
	contractAbi, _ := app.Account.ReadAbi(constants.BUNGEE_ABI)

	_chain := constants.SETTINGS.RefuelTo[app.Account.RandomInt(0, len(constants.SETTINGS.RefuelTo)-1)]

	destinationChainId, err := strconv.Atoi(constants.CHAINS[_chain]["CHAIN_ID"])
	if err != nil {
		return nil, err
	}

	calldata, err := contractAbi.Pack("depositNativeToken", big.NewInt(int64(destinationChainId)), app.Wallet.PublicKey)
	if err != nil {
		return nil, err
	}

	value := app.Account.RandomFloat(constants.SETTINGS.BungeeValueMin, constants.SETTINGS.BungeeValueMax, 8)

	txHash, err := app.Account.SendTransaction(
		app.Client, app.Wallet, contractAddr, calldata, app.Account.ToWei(value), "ETHEREUM",
	)
	if err != nil {
		return nil, err
	}

	return txHash, nil
}
