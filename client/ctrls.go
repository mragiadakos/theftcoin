package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/urfave/cli"
)

var GenerateKeyCommand = cli.Command{
	Name:    "generate",
	Aliases: []string{"g"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "filename",
			Usage: "the filename that the key will be saved",
		},
	},
	Usage: "generate the key in a file",
	Action: func(c *cli.Context) error {
		filename := c.String("filename")
		if len(filename) == 0 {
			return errors.New("Error: filename is missing")
		}
		privk, _, _ := crypto.GenerateKeyPair(crypto.Ed25519, 0)
		kj := KeyJson{}
		b, _ := privk.GetPublic().Bytes()
		kj.PublicKey = hex.EncodeToString(b)
		kj.PrivateKey, _ = crypto.MarshalPrivateKey(privk)
		b, _ = json.Marshal(kj)
		err := ioutil.WriteFile(filename, b, 0644)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		fmt.Println("The generate was successful")
		return nil
	},
}

var AddCommand = cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename that contains the key in json file",
		},
		cli.Float64Flag{
			Name:  "coins",
			Usage: "the the number of coins you want to add in your account",
		},
	},
	Usage: "add coins to your account as an inflator",
	Action: func(c *cli.Context) error {
		key := c.String("key")
		if len(key) == 0 {
			return errors.New("Error: the key is missing")
		}

		coins := c.Float64("coins")
		if coins <= 0 {
			return errors.New("Error: the coins are not allowed to be 0 or less")
		}

		privk, err := fileKey(key)
		if err != nil {
			return errors.New("Error client:" + err.Error())
		}

		_, err = Add(privk, coins)
		if err != nil {
			return errors.New("Error:" + err.Error())
		}
		fmt.Println("The add was successful")
		return nil
	},
}

var RemoveCommand = cli.Command{
	Name:    "remove",
	Aliases: []string{"r"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename that contains the key in json file",
		},
		cli.Float64Flag{
			Name:  "coins",
			Usage: "the the number of coins you want to add in your account",
		},
	},
	Usage: "remove coins from your account as an inflator",
	Action: func(c *cli.Context) error {
		key := c.String("key")
		if len(key) == 0 {
			return errors.New("Error: the key is missing")
		}

		coins := c.Float64("input")
		if coins <= 0 {
			return errors.New("Error: the coins are not allowed to be 0 or less")
		}

		privk, err := fileKey(key)
		if err != nil {
			return errors.New("Error client:" + err.Error())
		}

		_, err = Remove(privk, coins)
		if err != nil {
			return errors.New("Error:" + err.Error())
		}
		fmt.Println("The remove was successful")
		return nil
	},
}

var SendCommand = cli.Command{
	Name:    "send",
	Aliases: []string{"s"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename that contains the key in json file",
		},
		cli.StringFlag{
			Name:  "receiver",
			Usage: "the receiver's public key",
		},
		cli.StringFlag{
			Name:  "tax",
			Usage: "the IPFS hash of the tax",
		},
		cli.Float64Flag{
			Name:  "coins",
			Usage: "the the number of coins you want to add in your account",
		},
	},
	Usage: "send coins to another account as an inflator",
	Action: func(c *cli.Context) error {
		key := c.String("key")
		if len(key) == 0 {
			return errors.New("Error: the key is missing")
		}

		receiver := c.String("receiver")
		if len(receiver) == 0 {
			return errors.New("Error: the receiver is missing")
		}

		taxHash := c.String("tax")
		if len(taxHash) == 0 {
			return errors.New("Error: The tax is not included.")
		}

		coins := c.Float64("coins")
		if coins <= 0 {
			return errors.New("Error: the coins are not allowed to be 0 or less")
		}

		fromPrivk, err := fileKey(key)
		if err != nil {
			return errors.New("Error client:" + err.Error())
		}

		b, err := hex.DecodeString(receiver)
		if err != nil {
			return errors.New("Error client:" + err.Error())
		}

		_, err = Send(fromPrivk, b, taxHash, coins)
		if err != nil {
			return errors.New("Error:" + err.Error())
		}
		fmt.Println("The send was successful")
		return nil
	},
}

var QueryCommand = cli.Command{
	Name:    "query",
	Aliases: []string{"q"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename that contains the key in json file",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "the user to check if has watcher",
		},
	},
	Usage: "send coins to another account as an inflator",
	Action: func(c *cli.Context) error {
		key := c.String("key")
		if len(key) == 0 {
			return errors.New("Error: the key is missing")
		}

		var userB *[]byte
		user := c.String("user")
		if len(user) > 0 {
			b, err := hex.DecodeString(user)
			if err != nil {
				return errors.New("Error: the user's public key is not hex")
			}
			userB = &b
		}

		privk, err := fileKey(key)
		if err != nil {
			return errors.New("Error client:" + err.Error())
		}

		qresp, _, err := Query(privk, userB)
		if err != nil {
			return errors.New("Error:" + err.Error())
		}
		fmt.Println("Coins: ", qresp.Coins)
		return nil
	},
}
