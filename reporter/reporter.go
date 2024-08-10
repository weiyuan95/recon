package reporter

import "chaintx/chains"

type Transfer struct {
	Chain        chains.ChainName
	Address      string
	Txid         string
	Timestamp    string
	TransferType string
	TokenType    string
	From         string
	To           string
	Amount       string
}

type WatchedAddressInfo struct {
	Chain   chains.ChainName
	Address string
}
