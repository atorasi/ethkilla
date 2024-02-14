package deposit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"ethkilla/src/account"
	"ethkilla/src/constants"
	"fmt"
	"net/http"
	"time"
)

func NewDepositApp(wallet account.Wallet) DepositApp {
	return DepositApp{
		wallet:  wallet,
		account: account.Account{},
	}
}

func (app DepositApp) OkxWithdraw(chainName string) (*http.Response, error) {
	value := app.account.RandomFloat(constants.SETTINGS.OkxValueMin, constants.SETTINGS.OkxValueMax, 6)

	payload, err := json.Marshal(&WithdrawalParams{
		Currency:    "ETH",
		Amount:      value,
		Destination: "4",
		ToAddress:   app.wallet.PublicKey.String(),
		Fee:         constants.CHAINS[chainName]["WITHDRAWAL_FEE"],
		Chain:       constants.CHAINS[chainName]["OKX_CHAIN"],
	})
	if err != nil {
		return nil, err
	}
	timestamp, signature := app.sign(
		"POST",
		"/api/v5/asset/withdrawal",
		string(payload),
		constants.SETTINGS.OxkSecret,
	)

	request, err := http.NewRequest("POST", "https://www.okx.com/api/v5/asset/withdrawal", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	request.Header = http.Header{
		"Content-Type":         {"application/json"},
		"OK-ACCESS-KEY":        {constants.SETTINGS.OxkAPIKey},
		"OK-ACCESS-PASSPHRASE": {constants.SETTINGS.OxkPassword},
		"OK-ACCESS-SIGN":       {signature},
		"OK-ACCESS-TIMESTAMP":  {timestamp},
	}

	response, err := http.DefaultClient.Do(request)
	if response.StatusCode >= 200 && response.StatusCode < 300 && err == nil {
		return response, nil
	}
	defer response.Body.Close()

	return nil, err
}

func (app DepositApp) sign(method, path, body, secretKey string) (string, string) {
	format := "2006-01-02T15:04:05.999Z07:00"
	t := time.Now().UTC().Format(format)
	ts := fmt.Sprint(t)
	s := ts + method + path + body
	p := []byte(s)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(p)
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}
