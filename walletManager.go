package main

import (
	"bytes"
	"crypto/elliptic"
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
	wm = WallerManager{map[string]*Wallet{}}
	//加载数据库数据

	if !wm.loadFile() {
		return nil
	}

	return &wm
}

func (wm *WallerManager) CreateWallet() string {

	w := NewWalletKeyPair()
	if w == nil {
		fmt.Println("NewWalletKeyPair fail")
		return ""
	}

	address := w.getAddress()

	wm.Wallets[address] = w

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
	//注册接口函数
	gob.Register(elliptic.P256())
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

func (wm *WallerManager) loadFile() bool {
	if !isFileExist(walletFile) {
		fmt.Println("文件不存在，不用加载")
		return true
	}

	content, err := ioutil.ReadFile(walletFile)

	if err != nil {
		fmt.Println("ioutil.ReadFile err :", err)
		return false
	}

	//解码
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))

	err = decoder.Decode(wm)

	if err != nil {
		fmt.Println("decoder.Decode err :", err)
		return false
	}
	return true
}

func (wm WallerManager) listAddress() []string {
	var addresses []string
	for address, _ := range wm.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}
