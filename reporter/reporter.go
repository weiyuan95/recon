package reporter

type Transfer struct {
	Txid         string
	Timestamp    string
	TransferType string
	TokenType    string
	From         string
	To           string
	Amount       string
}

type WatchedAddressInfo struct {
	Address string
}
