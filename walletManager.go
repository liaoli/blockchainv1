package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

//负责对外，管理生成的钱包（公钥，私钥）
//私钥 ->公钥->地址

type WallerManager struct {
	//key 地址
	//value：wallet
	Wallets map[string]*Wallet
}

func NewWalletManager() *WallerManager {

	var wm WallerManager
	//加载数据库数据

	return &wm
}

func (wm *WallerManager) CreateWallet() string {

	w := NewWalletKeyPair()
	if w == nil {
		fmt.Println("NewWalletKeyPair fail")
		return ""
	}

	address := w.getAddress()

	//将密钥写入磁盘
	if !wm.saveFile() {
		return ""
	}

	return address
}

//
const walletFile = "wallet.dat"

func (wm *WallerManager) saveFile() bool {

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(wm)

	if err != nil {
		fmt.Println("encoder.Encode err:", err)
		return false
	}

	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)

	if err != nil {
		fmt.Println("ioutil.WriteFile err:", err)
		return false
	}

	return true
}
