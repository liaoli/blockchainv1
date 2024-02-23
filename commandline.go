package main

import (
	"fmt"
	"os"
)

type CLI struct {
	//bc *BlockChain
}

const Usage = `
    ./block create <address> "create block chain"
	./block addBlock --data DATA "add a block"
    ./block printChain "print block Chain"
    ./block getBalance <address> "Get balance by address"
`

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "create":
		fmt.Println("创建区块被调用")
		if len(os.Args) != 3 {
			fmt.Println("输入参数无效，请检查！")
			fmt.Println(Usage)
		}
		address := os.Args[2]
		cli.createBlockChain(address)
	case "addBlock":
		if len(os.Args) > 3 && os.Args[2] == "--data" {
			data := os.Args[3]
			if data == "" {
				fmt.Println("data should not be empty")
				os.Exit(1)
			}
			//cli.addBlock(data)
		}
	case "printChain":
		cli.PrintChain()
	case "getBalance":
		if len(os.Args) != 3 {
			fmt.Println("输入参数无效，请检查！")
			fmt.Println(Usage)
			os.Exit(1)
		}
		fmt.Println("调用获取余额命令")
		address := os.Args[2]
		cli.getBalance(address)
	}
}

func (cli CLI) createBlockChain(address string) {
	if address == "" {
		fmt.Println("传入的地址无效:", address)
	}

	NewBlockChain(address)

}

func (cli *CLI) addBlock(data string) {
	//cli.bc.AddBlock(data)
}

func (cli *CLI) PrintChain() {
	//it := NewBlockChainIterator(cli.bc)
	//
	//for {
	//	block := it.GetBlockAnMoveLeft()
	//	fmt.Println("============== ===============")
	//	fmt.Printf("PreHash : %x\n", block.PreHash)
	//	fmt.Printf("Hash : %x\n", block.Hash)
	//	//fmt.Printf("Data : %s\n", block.Data)
	//	fmt.Printf("merkleRoot : %x\n", block.merkleRoot)
	//	fmt.Printf("Nonce : %d\n", block.Nonce)
	//	fmt.Printf("Version : %d\n", block.Version)
	//	//时间格式化
	//	timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
	//	//fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
	//	fmt.Printf("TimeStamp : %s\n", timeFormat)
	//	fmt.Printf("Difficulty : %d\n", block.Difficulty)
	//	pow := NewProofOfWork(block)
	//	fmt.Printf("IsValid:%v \n", pow.IsValid())
	//
	//	if len(block.PreHash) == 0 {
	//		break
	//	}
	//}
}

func (cli *CLI) getBalance(address string) {
	bc, _ := GetBlockChainInstance()
	utxos := bc.FindMyUTXO(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("地址 %s total：%f ", address, total)
}
