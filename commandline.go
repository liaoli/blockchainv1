package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type CLI struct {
	//bc *BlockChain
}

const Usage = `
    ./block create <address> "create block chain"
	./block addBlock --data DATA "add a block"
    ./block printChain "print block Chain"
    ./block getBalance <address> "Get balance by address"
    ./block send <From> <To> <Amount> <Miner> <Data> "transaction"
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
	case "send":
		fmt.Println("send 命令被调用")
		if len(os.Args) != 7 {
			fmt.Println("输入参数无效，请检查")
			return
		}
		from := os.Args[2]
		to := os.Args[3]
		amount, _ := strconv.ParseFloat(os.Args[4], 64)
		miner := os.Args[5]
		data := os.Args[6]

		cli.send(from, to, amount, miner, data)
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
	bc, err := GetBlockChainInstance()

	if err != nil {
		fmt.Println("print err:", err)
		return
	}
	defer bc.db.Close()
	it := bc.NewIterator()

	for {
		block := it.GetBlockAnMoveLeft()
		fmt.Println("============== ===============")
		fmt.Printf("PreHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		//fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Version : %d\n", block.Version)
		//时间格式化
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		//fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficulty : %d\n", block.Difficulty)
		pow := NewProofOfWork(&block)
		fmt.Printf("IsValid:%v \n", pow.IsValid())

		if len(block.PreHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	bc, err := GetBlockChainInstance()

	if err != nil {
		fmt.Println("print err:", err)
		return
	}
	defer bc.db.Close()
	utxos := bc.FindMyUTXO(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("地址 %s total：%f ", address, total)
}

func (cli *CLI) send(from, to string, amount float64, miner, data string) {
	//fmt.Println("from:", from)
	//fmt.Println("to:", to)
	//fmt.Println("amount:", amount)
	//fmt.Println("miner:", miner)
	//fmt.Println("data:", data)
	bc, err := GetBlockChainInstance()

	if err != nil {
		fmt.Println("print err:", err)
		return
	}
	defer bc.db.Close()
	//每次send时，都会添加一个块，
	//区块：创建挖矿交易，创建普通交易
	//执行addblock

	//1。创建挖矿交易
	coinbaseTx := NewCoinbaseTx(miner, data)

	//创建有效交易数组，将有效交易添加进来

	txs := []*Transaction{coinbaseTx}
	//2.创建普通交易
	tx := NewTransaction(from, to, amount, *bc)
	if tx != nil {
		fmt.Println("找到一笔有效的转账交易！")
		txs = append(txs, tx)
	} else {
		fmt.Println("找到一笔无效交易")
	}

	err = bc.AddBlock(txs)

	if err != nil {
		fmt.Println("添加区块失败，交易失败")
	}
	fmt.Println("添加区块成功，交易成果")

}
