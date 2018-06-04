package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/theftcoin/server/confs"
	abcicli "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/types"
)

type KeyJson struct {
	PublicKey  string // hex
	PrivateKey []byte
}

/*
func deliver(b []byte) (uint32, error) {
	resp, err := RpcBroadcastCommit(b)
	if err != nil {
		return CodeTypeClientError, err
	}
	if resp.Code > CodeTypeOK {
		return resp.Code, errors.New(resp.Log)
	}
	return CodeTypeOK, nil
}

func query(b []byte) (*QueryResponse, uint32, error) {
	resp, err := RpcQuery(b)
	if err != nil {
		return nil, CodeTypeClientError, err
	}
	if resp.Code > CodeTypeOK {
		return nil, resp.Code, errors.New(resp.Log)
	}

	qresp := QueryResponse{}
	json.Unmarshal(resp.Value, &qresp)

	return &qresp, CodeTypeOK, nil
}
*/

func deliver(b []byte) (uint32, error) {
	client := abcicli.NewSocketClient(confs.Conf.AbciDaemon, false)
	defer func() {
		client.Stop()
	}()
	err := client.Start()
	if err != nil {
		return CodeTypeClientError, err
	}
	resp, err := client.DeliverTxSync(b)
	if err != nil {
		return CodeTypeClientError, err
	}
	if resp.Code > CodeTypeOK {
		return resp.Code, errors.New(resp.Log)
	}
	return CodeTypeOK, nil
}

func query(b []byte) (*QueryResponse, uint32, error) {
	client := abcicli.NewSocketClient(confs.Conf.AbciDaemon, false)
	defer func() {
		client.Stop()
	}()
	err := client.Start()
	if err != nil {
		return nil, CodeTypeClientError, err
	}
	req := types.RequestQuery{}
	req.Data = b
	resp, err := client.QuerySync(req)
	if err != nil {
		return nil, CodeTypeClientError, err
	}

	if resp.Code > CodeTypeOK {
		return nil, resp.Code, errors.New(resp.Log)
	}
	qresp := QueryResponse{}
	json.Unmarshal(resp.Value, &qresp)
	return &qresp, CodeTypeOK, nil
}

func Add(from crypto.PrivKey, coins float64) (uint32, error) {
	var err error
	dd := DeliveryData{}
	dd.From, err = from.GetPublic().Bytes()
	if err != nil {
		return CodeTypeClientError, err
	}
	dd.Action = ADD_ACTION
	dd.Coins = coins
	b, _ := json.Marshal(dd)
	dr := DeliveryRequest{}
	dr.Signature, err = from.Sign(b)
	if err != nil {
		return CodeTypeClientError, err
	}
	dr.Date = time.Now().UTC()
	dr.Data = dd
	b, _ = json.Marshal(dr)
	return deliver(b)

}

func Remove(from crypto.PrivKey, coins float64) (uint32, error) {
	var err error
	dd := DeliveryData{}
	dd.From, err = from.GetPublic().Bytes()
	if err != nil {
		return CodeTypeClientError, err
	}
	dd.Action = REMOVE_ACTION
	dd.Coins = coins
	b, _ := json.Marshal(dd)
	dr := DeliveryRequest{}
	dr.Signature, err = from.Sign(b)
	if err != nil {
		return CodeTypeClientError, err
	}
	dr.Date = time.Now().UTC()
	dr.Data = dd
	b, _ = json.Marshal(dr)
	return deliver(b)
}

func Send(from crypto.PrivKey, toPublicKey []byte, taxHash string, coins float64) (uint32, error) {
	var err error
	dd := DeliveryData{}
	dd.From, err = from.GetPublic().Bytes()
	if err != nil {
		return CodeTypeClientError, err
	}
	dd.Action = SEND_ACTION
	dd.To = &toPublicKey
	dd.Coins = coins
	dd.TaxHash = &taxHash
	b, _ := json.Marshal(dd)
	dr := DeliveryRequest{}
	dr.Signature, err = from.Sign(b)
	dr.Date = time.Now().UTC()

	if err != nil {
		return CodeTypeClientError, err
	}
	dr.Data = dd
	b, _ = json.Marshal(dr)
	return deliver(b)
}

func Query(from crypto.PrivKey, userAddr *[]byte) (*QueryResponse, uint32, error) {
	var err error
	q := QueryRequest{}
	data := QueryData{}
	data.From, err = from.GetPublic().Bytes()
	if err != nil {
		return nil, CodeTypeClientError, err
	}

	if userAddr != nil {
		data.User = userAddr
	}
	data.Date = time.Now().UTC()
	b, _ := json.Marshal(data)

	q.Data = data
	q.Signature, err = from.Sign(b)
	if err != nil {
		return nil, CodeTypeClientError, err
	}
	b, _ = json.Marshal(q)
	return query(b)
}

func fileKey(filename string) (crypto.PrivKey, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Error: " + err.Error())
	}

	kj := KeyJson{}
	err = json.Unmarshal(b, &kj)
	if err != nil {
		return nil, errors.New("Error: json problem with the key " + err.Error())
	}

	edKey, err := crypto.UnmarshalPrivateKey(kj.PrivateKey)
	if err != nil {
		return nil, errors.New("Error: private key decoding problem with the key " + err.Error())
	}
	return edKey, nil
}
