package account

import (
	"context"
	"errors"
	"ethkilla/src/constants"
	"ethkilla/utils"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (acc Account) SendTransaction(
	client *ethclient.Client,
	wallet Wallet,
	contract common.Address,
	calldata []byte,
	value *big.Int,
	chainName string,
) (txHash *types.Receipt, err error) {
	log.Printf("Acc.%d | Building transaction", wallet.Index)
	message := ethereum.CallMsg{
		From:  wallet.PublicKey,
		To:    &contract,
		Data:  calldata,
		Value: value,
	}

	gasPrice, gasLimit, gasErr := acc.EstimateFees(client, message)
	if gasErr != nil && calldata == nil {
		gasLimit = 21000
		gasErr = nil
	}
	nonce, nonceErr := client.PendingNonceAt(context.Background(), wallet.PublicKey)
	if nonceErr != nil || gasErr != nil {
		return nil, nonceErr
	}

	chainID, _ := client.ChainID(context.Background())
	tx := types.NewTransaction(nonce, contract, value, gasLimit, gasPrice, calldata)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), wallet.PrivateKey)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100)*time.Second)
	defer cancel()

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	txHash, err = acc.waitReciept(client, signedTx)
	if err != nil || txHash != nil {
		return nil, errors.New("transaction waiting error")
	}

	log.Printf("Acc.%d | Transaction Hash: %s/tx/%v", wallet.Index, constants.CHAINS[chainName]["SCAN"], txHash.TxHash)
	utils.SendTelegramMessage("Acc.%d | Transaction Hash: %s/tx/%v", wallet.Index, constants.CHAINS[chainName]["SCAN"], txHash.TxHash)

	if constants.SETTINGS.NeedDelayAct {
		delay := acc.RandomInt(constants.SETTINGS.DelayActMin, constants.SETTINGS.DelayActMax)
		log.Printf("Acc.%d | sleep for %d seconds before the next activity", wallet.Index, delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	return txHash, err
}

func (acc Account) Approve(
	client *ethclient.Client, wallet Wallet, token,
	approveFor common.Address, chainName string,
) (*types.Receipt, error) {
	contractAbi, _ := acc.ReadAbi(constants.USDC_ABI)
	valueToApprove, _ := acc.TokenBalance(client, wallet, token)

	calldata, err := contractAbi.Pack("approve", approveFor, valueToApprove)
	if err != nil {
		return nil, err
	}

	return acc.SendTransaction(client, wallet, token, calldata, big.NewInt(0), chainName)
}

func (acc Account) ReadAbi(abiString string) (abi.ABI, error) {
	myContractAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return abi.ABI{}, err
	}

	return myContractAbi, nil
}

func (acc Account) TokenBalance(client *ethclient.Client, wallet Wallet, tokenAddress common.Address) (*big.Int, error) {
	contractAbi, _ := acc.ReadAbi(constants.USDC_ABI)
	calldata, _ := contractAbi.Pack("balanceOf", wallet.PublicKey)

	res, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: calldata,
	},
		nil)
	if err != nil {
		return nil, err
	}

	result, _ := contractAbi.Unpack("balanceOf", res)

	return result[0].(*big.Int), nil
}

func (acc Account) NativeBalance(client *ethclient.Client, wallet Wallet) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), wallet.PublicKey, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (acc Account) EstimateFees(client *ethclient.Client, message ethereum.CallMsg) (gasPrice *big.Int, gasLimit uint64, err error) {
	gasPrice, priceErr := client.SuggestGasPrice(context.Background())
	if gasPrice.Cmp(big.NewInt(int64(constants.SETTINGS.MaxGwei*1000000000))) != -1 {
		log.Println("| Waiting for lower gwei. Sleep 15 seconds")
		time.Sleep(time.Duration(15) * time.Second)
		return acc.EstimateFees(client, message)
	}
	if priceErr != nil {
		return nil, 0, priceErr
	}
	gasLimit, limitErr := client.EstimateGas(context.Background(), message)
	if limitErr != nil {
		return gasPrice, 0, err
	}

	gasPrice.Mul(gasPrice, big.NewInt(11))
	gasPrice.Div(gasPrice, big.NewInt(10))

	return gasPrice, uint64(float64(gasLimit) * 1.54), nil
}

func (acc Account) RandomValue(value *big.Int, min, max int) *big.Int {
	floatValue := new(big.Float).SetInt(value)
	randomPercent := float64(acc.RandomInt(min, max)) / 100

	percent := new(big.Float).SetFloat64(randomPercent)
	result := new(big.Float).Mul(floatValue, percent)

	intResult, _ := result.Int(nil)
	return intResult
}

func (acc Account) RandomInt(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := rand.Intn(max-min+1) + min

	return randomNum
}

func (acc Account) RandomFloat(min, max float64, places int) float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := min + rand.Float64()*(max-min)

	factor := math.Pow(10, float64(places))
	return math.Round(randomNum*factor) / factor
}

func (acc Account) ToWei(amount float64) *big.Int {
	weiAmount := big.NewFloat(amount)
	weiAmount = new(big.Float).Mul(weiAmount, big.NewFloat(1e18))

	wei, _ := weiAmount.Int(nil)
	return wei
}

func (acc Account) waitReciept(client *ethclient.Client, signedTx *types.Transaction) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	mined, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return nil, err
	}

	return mined, nil
}
