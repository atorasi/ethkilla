package account

import (
	"ethkilla/src/constants"
	"ethkilla/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewClient(index int, proxyList []string, chainName string) (*ethclient.Client, error) {
	if constants.SETTINGS.UseProxy {
		proxies := strings.Split(proxyList[index], "@")
		userpass := strings.Split(proxies[0], ":")

		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: "http",
					Host:   proxies[1],
					User:   url.UserPassword(userpass[0], userpass[1]),
				}),
			},
		}

		rpcClient, err := rpc.DialHTTPWithClient(constants.CHAINS[chainName]["RPC"], httpClient)
		sideClient := ethclient.NewClient(rpcClient)
		if err != nil {
			return nil, err
		}
		return sideClient, nil
	} else {
		sideClient, err := ethclient.Dial(constants.CHAINS[constants.SETTINGS.SideChain]["RPC"])
		if err != nil {
			return nil, err
		}
		return sideClient, nil
	}
}

func SliceOfAccs() ([]Wallet, error) {
	keys, err := utils.ReadFile(`..\data\pKeys.txt`)
	if err != nil {
		return nil, err
	}

	accounts, err := Wallets(keys)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func Wallets(pks []string) ([]Wallet, error) {
	var accounts []Wallet
	for index, pk := range pks {
		account := newAccount(pk, index+1)

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func newAccount(private string, index int) Wallet {
	if private[0:2] == "0x" {
		private = private[2:]
	}
	privateKey, _ := crypto.HexToECDSA(private)

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)

	return Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Index:      index,
	}
}
