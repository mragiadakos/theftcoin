package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ipfs/go-ipfs-api"

	"github.com/mragiadakos/theftcoin/server/confs"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

type KeyJson struct {
	PublicKeyHex string
	PrivateKey   []byte
}

func CreateTaxPublicKey() {
	fmt.Println("Creating the tax files:")
	fileTax := "tax.json"
	fileTaxPriv := "tax_priv.json"
	privk, pubk, _ := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	kj := KeyJson{}
	pubB, _ := pubk.Bytes()
	hexPubk := hex.EncodeToString(pubB)
	kj.PublicKeyHex = hexPubk
	kj.PrivateKey, _ = privk.Bytes()
	b, _ := json.Marshal(kj)
	ioutil.WriteFile(fileTaxPriv, b, 0644)
	fmt.Println("The taxer's private key file " + fileTaxPriv)

	tax := confs.Tax{}
	tax.PublicKeyHex = hexPubk
	tax.Percentage = 10
	b, _ = json.Marshal(tax)
	ioutil.WriteFile("tax.json", b, 0644)
	fmt.Println("The taxer's public key file " + fileTax)

	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(fileTax)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println("IPFS Hash for tax: " + hash)

}

func CreateWatchersPublicKey() {
	fmt.Println("Creating the watchers files:")
	fileWatchers := "watchers.json"
	fileWatcherPriv := "watcher_priv.json"
	privk, pubk, _ := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	kj := KeyJson{}
	pubB, _ := pubk.Bytes()
	hexPubk := hex.EncodeToString(pubB)
	kj.PublicKeyHex = hexPubk
	kj.PrivateKey, _ = privk.Bytes()
	b, _ := json.Marshal(kj)
	ioutil.WriteFile(fileWatcherPriv, b, 0644)
	fmt.Println("The watcher's private key " + fileWatcherPriv)

	watchers := []confs.Watcher{}
	w := confs.Watcher{}
	w.PublicKeyHex = hexPubk
	watchers = append(watchers, w)
	b, _ = json.Marshal(watchers)
	ioutil.WriteFile(fileWatchers, b, 0644)
	fmt.Println("The watchers' file for the server is " + fileWatchers)

	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(fileWatchers)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println("IPFS Hash for watchers: " + hash)
}

func CreateInflatorsPublicKey() {
	fmt.Println("Creating the inflators files:")
	fileInflators := "inflators.json"
	fileInflatorPriv := "inflator_priv.json"
	privk, pubk, _ := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	kj := KeyJson{}
	pubB, _ := pubk.Bytes()
	hexPubk := hex.EncodeToString(pubB)
	kj.PublicKeyHex = hexPubk
	kj.PrivateKey, _ = privk.Bytes()
	b, _ := json.Marshal(kj)
	ioutil.WriteFile(fileInflatorPriv, b, 0644)
	fmt.Println("The inflator's private key " + fileInflatorPriv)

	infs := []confs.Inflator{}
	i := confs.Inflator{}
	i.PublicKeyHex = hexPubk
	infs = append(infs, i)
	b, _ = json.Marshal(infs)
	ioutil.WriteFile(fileInflators, b, 0644)
	fmt.Println("The inflators' file for the server is " + fileInflators)

	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(fileInflators)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println("IPFS Hash for inflators: " + hash)
}
