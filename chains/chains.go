package chains

type ChainName string

const (
	Ethereum        ChainName = "ethereum"
	EthereumSepolia ChainName = "ethereumSepolia"
	Polygon         ChainName = "polygon"
	Tron            ChainName = "tron"
	TronShasta      ChainName = "tronShasta"
)

func IsValidChain(chain ChainName) bool {
	switch chain {
	case Ethereum, EthereumSepolia, Polygon, Tron, TronShasta:
		return true
	default:
		return false
	}
}
