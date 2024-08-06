package chains

type ChainName string

const (
	Ethereum        ChainName = "ethereum"
	EthereumSepolia ChainName = "ethereumSepolia"
	Polygon         ChainName = "polygon"
)

func IsValidChain(chain ChainName) bool {
	switch chain {
	case Ethereum, EthereumSepolia, Polygon:
		return true
	default:
		return false
	}
}
