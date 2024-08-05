package tron

import (
	"chaintx/reporter"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type TransferItem struct {
	TransactionHash, TransferFromAddress, TransferToAddress string
	Confirmed                                               bool
	Amount                                                  int // In SUN
	Timestamp                                               int
}

type TransferResult struct {
	Total int
	Data  []TransferItem
}

func Transfers(tokenAddress string) ([]reporter.Transfer, error) {
	// The limit has a cap of 50 which is not documented
	resultLimit := 50
	start := 0
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	url := fmt.Sprintf(
		"https://apilist.tronscanapi.com/api/trx/transfer?sort=-timestamp&count=true&limit=%d&start=%d&address=%s&filterTokenValue=0",
		resultLimit,
		start,
		tokenAddress,
	)

	get, err := client.Get(url)
	if err != nil {
		return []reporter.Transfer{}, errors.New("failed to fetch data")
	}

	body, err := io.ReadAll(get.Body)
	defer get.Body.Close()

	var data TransferResult
	unmarshalErr := json.Unmarshal(body, &data)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
		return []reporter.Transfer{}, errors.New("failed to unmarshal data")
	}

	var result []reporter.Transfer

	for len(data.Data) > 0 {
		time.Sleep(1000 * time.Millisecond)

		for _, transfer := range data.Data {
			result = append(result, reporter.Transfer{
				Txid:         transfer.TransactionHash,
				From:         transfer.TransferFromAddress,
				To:           transfer.TransferToAddress,
				Timestamp:    strconv.Itoa(transfer.Timestamp),
				Amount:       strconv.Itoa(transfer.Amount),
				TransferType: "native",
			})

		}

		start += len(data.Data)
		url = fmt.Sprintf(
			"https://apilist.tronscanapi.com/api/trx/transfer?sort=-timestamp&count=true&limit=%d&start=%d&address=%s&filterTokenValue=0",
			resultLimit,
			start,
			tokenAddress,
		)

		get, err = client.Get(url)
		if err != nil {
			return []reporter.Transfer{}, errors.New("failed to fetch data")
		}

		body, err = io.ReadAll(get.Body)
		defer get.Body.Close()

		unmarshalErr = json.Unmarshal(body, &data)
		if unmarshalErr != nil {
			fmt.Println(unmarshalErr)
			return []reporter.Transfer{}, errors.New("failed to unmarshal data")
		}
	}

	return result, nil
}
