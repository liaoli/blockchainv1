package main

type BlockChain struct {
	//定义一个区块链的结构
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	//对区块链进行初始化，并把创世块添加到区块链
	genesisBlock := GenesisBlock(genesisInfo, []byte{})
	bc := BlockChain{[]*Block{genesisBlock}}
	return &bc
}

func (receiver *BlockChain) AddBlock(data string) {

	//1.产生区块链
	//  数据和区块的哈希值
	//  通过数组的下标，拿到最后一个区块的哈希值，这个哈希值就是我们新区块的前哈希值
	blockLen := len(receiver.Blocks)

	lastBlock := receiver.Blocks[blockLen-1]
	prevBlockHash := lastBlock.Hash
	block := NewBlock(data, prevBlockHash)
	receiver.Blocks = append(receiver.Blocks, block)

}
