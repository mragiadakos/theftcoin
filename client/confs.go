package main

import (
	"encoding/json"
	"errors"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

type configuration struct {
	NodeDaemon     string
	IpfsConnection string
}

const (
	CodeTypeOK            uint32 = 0
	CodeTypeEncodingError uint32 = 1
	CodeTypeBadNonce      uint32 = 2
	CodeTypeUnauthorized  uint32 = 3
	CodeTypeClientError   uint32 = 4
)

var Conf = configuration{}

func init() {
	Conf.NodeDaemon = "http://localhost:46657"
	Conf.IpfsConnection = "127.0.0.1:5001"
}

type ActionStruct string

const (
	ADD_ACTION    = ActionStruct("add")
	REMOVE_ACTION = ActionStruct("remove")
	SEND_ACTION   = ActionStruct("send")
)

type DeliveryData struct {
	From    []byte  // public key
	To      *[]byte // public key
	Action  ActionStruct
	TaxHash *string
	Coins   float64
}

type DeliveryRequest struct {
	Signature []byte
	Date      time.Time
	Data      DeliveryData
}

func (dr *DeliveryRequest) VerifySignature() (bool, error) {
	pub, err := crypto.UnmarshalPublicKey(dr.Data.From)
	if err != nil {
		return false, errors.New("The sender's public key is not correct")
	}
	b, _ := json.Marshal(dr.Data)
	ver, err := pub.Verify(b, dr.Signature)
	if err != nil {
		return false, errors.New("The signature's format is not correct.")
	}
	return ver, nil
}

type QueryData struct {
	From  []byte // public key
	Date  time.Time
	Nonce string
	User  *[]byte
}

type QueryRequest struct {
	Signature []byte
	Data      QueryData
}

func (qr *QueryRequest) VerifySignature() (bool, error) {
	pub, err := crypto.UnmarshalPublicKey(qr.Data.From)
	if err != nil {
		return false, errors.New("The public key is not correct")
	}
	b, _ := json.Marshal(qr.Data)
	ver, err := pub.Verify(b, qr.Signature)
	if err != nil {
		return false, errors.New("The signature's format is not correct.")
	}
	return ver, nil
}

type QueryResponse struct {
	Coins float64
}
