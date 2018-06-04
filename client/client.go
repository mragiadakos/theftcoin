package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tendermint/abci/types"
)

type jsonRpcRequest struct {
	Method  string      `json:"method"`  //"method": "broadcast_tx_sync",
	Version string      `json:"jsonrpc"` //"jsonrpc": "2.0",
	Params  interface{} `json:"params"`  //"params": ,
	Id      string      `json:"id"`      //"id": "dontcare"
}

type jsonRpcResponseForDelivery struct {
	Method  string     `json:"method"`  //"method": "broadcast_tx_sync",
	Version string     `json:"jsonrpc"` //"jsonrpc": "2.0",
	Result  deliveryTx `json:"result"`  //"result": ,
	Id      string     `json:"id"`      //"id": "dontcare"
}

type deliveryTx struct {
	DeliveryTx types.ResponseDeliverTx `json:"deliver_tx"`
	Height     int                     `json:"height"`
	Hash       string                  `json:"hash"`
}

type jsonRpcResponseForQuery struct {
	Method  string        `json:"method"`  //"method": "broadcast_tx_sync",
	Version string        `json:"jsonrpc"` //"jsonrpc": "2.0",
	Result  responseQuery `json:"result"`  //"result": ,
	Id      string        `json:"id"`      //"id": "dontcare"
}

type responseQuery struct {
	Response types.ResponseQuery `json:"response"`
}

func newJsonRpcRequest(method string, js interface{}) jsonRpcRequest {
	jr := jsonRpcRequest{}
	jr.Method = method
	jr.Id = "dontcare"
	jr.Params = js
	jr.Version = "2.0"
	return jr
}

type Tx struct {
	Tx string `json:"tx"`
}

type AbciQuery struct {
	Data string `json:"data"`
}

func RpcBroadcastCommit(deliveryB []byte) (*types.ResponseDeliverTx, error) {
	tx := Tx{}
	tx.Tx = base64.StdEncoding.EncodeToString(deliveryB)
	jr := newJsonRpcRequest("broadcast_tx_commit", tx)
	bout, _ := json.Marshal(jr)
	resp, err := http.Post(Conf.NodeDaemon, "text/plain", bytes.NewBuffer(bout))
	if err != nil {
		return nil, err
	}
	bresp, _ := ioutil.ReadAll(resp.Body)
	jresp := jsonRpcResponseForDelivery{}
	json.Unmarshal(bresp, &jresp)

	return &jresp.Result.DeliveryTx, nil
}

func RpcQuery(b []byte) (*types.ResponseQuery, error) {
	aq := AbciQuery{}
	aq.Data = hex.EncodeToString(b)
	jr := newJsonRpcRequest("abci_query", aq)
	bout, _ := json.Marshal(jr)
	resp, err := http.Post(Conf.NodeDaemon, "text/plain", bytes.NewBuffer(bout))
	if err != nil {
		return nil, err
	}
	bresp, _ := ioutil.ReadAll(resp.Body)
	jresp := jsonRpcResponseForQuery{}
	json.Unmarshal(bresp, &jresp)
	return &jresp.Result.Response, nil
}
