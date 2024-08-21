package evm

import (
	"errors"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"recon/chains"
)

type evmChain struct {
	Name   chains.ChainName
	RpcUrl []string
}

var evmChains = map[chains.ChainName]evmChain{
	chains.Ethereum: {
		Name: chains.Ethereum,
		RpcUrl: []string{
			"https://eth-pokt.nodies.app",
			"https://ethereum.blockpi.network/v1/rpc/public",
			"https://eth-mainnet.public.blastapi.io",
			"https://eth.drpc.org",
			"https://1rpc.io/eth",
		},
	},
	chains.EthereumSepolia: {
		Name: chains.EthereumSepolia,
		RpcUrl: []string{
			"https://ethereum-sepolia-rpc.publicnode.com",
			"https://rpc-sepolia.rockx.com",
			"https://eth-sepolia.public.blastapi.io",
		},
	},
	chains.Polygon: {
		Name: chains.Polygon,
		RpcUrl: []string{
			"https://polygon-pokt.nodies.app",
			"https://polygon.llamarpc.com",
			"wss://polygon.drpc.org",
		},
	},
}

// A cache of clients for each chain
var clients = make(map[chains.ChainName]*ethclient.Client)

// GetClient returns a client for the given chain
func GetClient(chain chains.ChainName) (*ethclient.Client, error) {
	if !chains.IsValidChain(chain) {
		return nil, errors.New("invalid chain")
	}

	client := clients[chain]
	if client == nil {
		client = createClient(chain)
		clients[chain] = client
	}
	return client, nil
}

func createClient(chain chains.ChainName) *ethclient.Client {
	// Connect to the chain
	client, err := ethclient.Dial(evmChains[chain].RpcUrl[0])
	if err != nil {
		log.Println("Failed to connect to the chain:", err)
		panic("Could not create a client for the chain.")
	}
	return client
}
