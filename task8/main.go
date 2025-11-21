package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"test8/counter"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, chainID := Connect()
	_, _ = client, chainID
	// GetBlockInfo(client, chainID)
	// SendTransaction(client, chainID)
	CallContract(client, chainID)
}

/**
 * @brief 连接测试网
 *
 * @param void
 * @return *ethclient.Client, *big.Int
 */
func Connect() (*ethclient.Client, *big.Int) {
	// 连接到测试网节点
	// client, _ := ethclient.Dial("https://sepolia.infura.io/v3/8a555900313d47f88d80e297034f1433")
	client, _ := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/1fhgxXSLQBfWGEydzGHFp")

	// 查询链ID
	chainID, _ := client.ChainID(context.Background())
	log.Println("Chain ID:", chainID) // 11155111

	return client, chainID
}

/**
 * @brief 查询区块信息
 *
 * @param client *ethclient.Client
 * @param chainID *big.Int
 * @return void
 */
func GetBlockInfo(client *ethclient.Client, chainID *big.Int) {
	// 查询指定区块信息
	block, _ := client.BlockByNumber(context.Background(), nil)
	log.Println("Block Hash:", block.Hash().Hex())
	log.Println("Block Timestamp:", block.Time())
	log.Println("Block Transaction:", block.Transactions().Len())
}

/**
 * @brief 发送交易-ETH转账
 *
 * @param client *ethclient.Client
 * @param chainID *big.Int
 * @return void
 */
func SendTransaction(client *ethclient.Client, chainID *big.Int) {
	privateKey, _ := crypto.HexToECDSA("0x-"[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	toAddress := common.HexToAddress("0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2")
	value := big.NewInt(10000000000000000) // 0.01 ETH
	gasLimit := uint64(21000)              // in units
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	txSign, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	err := client.SendTransaction(context.Background(), txSign)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("tx sent: %s", txSign.Hash().Hex())
}

/**
 * @brief 使用 Go 绑定代码调用智能合约
 *
 * @param client *ethclient.Client
 * @param chainID *big.Int
 * @return void
 */
func CallContract(client *ethclient.Client, chainID *big.Int) {
	privateKey, _ := crypto.HexToECDSA("0x-"[2:])
	tokenAddress := common.HexToAddress("0x7c53302282998Bc18d629434382a0de334F9764D")
	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	ins, _ := counter.NewCounter(tokenAddress, client)
	tx, _ := ins.Increment(auth)
	log.Println("Transaction Hash:", tx.Hash().Hex())
}
