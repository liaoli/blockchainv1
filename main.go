package main

import "fmt"

func main() {
	//fmt.Printf("hello world!")
	//block := GenesisBlock(genesisInfo, []byte{})
	//fmt.Printf("PreHash : %x\n", block.PreHash)
	//fmt.Printf("Hash : %x\n", block.Hash)
	//fmt.Printf("Data : %s\n", block.Data)

	bc := NewBlockChain()
	bc.AddBlock("aa send 1 btc to c")
	bc.AddBlock("bb send 1 btc to c")

	for i, block := range bc.Blocks {
		fmt.Println("==============block height：", i, "===============")
		fmt.Printf("PreHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("merkleRoot : %x\n", block.merkleRoot)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Difficulty : %d\n", block.Difficulty)
		pow := NewProofOfWork(*block)
		fmt.Printf("IsValid:%v \n", pow.IsValid())
	}
}
