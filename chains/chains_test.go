package chains

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChains(t *testing.T) {
	assert.Equal(t, Ethereum, ChainName("ethereum"))
	assert.Equal(t, EthereumSepolia, ChainName("ethereumSepolia"))
	assert.Equal(t, Polygon, ChainName("polygon"))
	assert.Equal(t, Tron, ChainName("tron"))
	assert.Equal(t, TronShasta, ChainName("tronShasta"))
}
