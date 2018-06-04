package ctrls

import (
	"encoding/json"
	"errors"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/theftcoin/server/confs"
	"github.com/tendermint/abci/types"
)

func (tca *TCApplication) validateQuery(qr QueryRequest) (uint32, error) {
	now := time.Now().UTC()
	since := now.Sub(qr.Data.Date)
	if since > time.Duration(time.Duration(confs.Conf.WaitingRequestTime)*time.Second) {
		return CodeTypeUnauthorized, errors.New("Request passed its time.")
	}

	if qr.Data.User != nil {
		if !confs.Conf.WatcherExists(string(qr.Data.From)) {
			return CodeTypeUnauthorized, errors.New("You are not a watcher.")
		}
	}
	ver, err := qr.VerifySignature()
	if err != nil {
		return CodeTypeEncodingError, err
	}
	if !ver {
		return CodeTypeUnauthorized, errors.New("The signature does not validate the query.")
	}
	return CodeTypeOK, nil
}

func (tca *TCApplication) Query(qreq types.RequestQuery) types.ResponseQuery {
	qr := QueryRequest{}
	err := json.Unmarshal(qreq.Data, &qr)
	if err != nil {

		return types.ResponseQuery{Code: CodeTypeEncodingError, Log: "The query request is not json."}
	}
	code, err := tca.validateQuery(qr)
	if err != nil {
		return types.ResponseQuery{Code: code, Log: err.Error()}
	}

	qresp := QueryResponse{}
	if qr.Data.User == nil {

		from, _ := crypto.UnmarshalPublicKey(qr.Data.From)
		cj, _ := tca.state.GetCoins(from)
		qresp.Coins = cj.Coins
	} else {
		user, _ := crypto.UnmarshalPublicKey(*qr.Data.User)
		cj, _ := tca.state.GetCoins(user)
		qresp.Coins = cj.Coins
	}
	b, _ := json.Marshal(qresp)
	resp := types.ResponseQuery{Code: CodeTypeOK, Value: b}
	return resp
}
