package ctrls

import (
	"encoding/json"
	"errors"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/theftcoin/server/confs"
	"github.com/tendermint/abci/types"
)

func (tca *TCApplication) validateInflators(dr DeliveryRequest) (uint32, error) {
	ok := confs.Conf.InflatorExists(string(dr.Data.From))
	if !ok {
		return CodeTypeUnauthorized, errors.New("You are not inflator.")
	}
	return CodeTypeOK, nil
}

func (tca *TCApplication) validateSend(dr DeliveryRequest) (uint32, error) {
	if dr.Data.To == nil {
		return CodeTypeUnauthorized, errors.New("The receiver's public key is empty.")
	}
	b := *dr.Data.To
	_, err := crypto.UnmarshalPublicKey(b)
	if err != nil {
		return CodeTypeEncodingError, errors.New("The receiver's public key is not correct." + err.Error())
	}

	if dr.Data.TaxHash == nil {
		return CodeTypeUnauthorized, errors.New("The tax is not included.")
	}

	if confs.Conf.IpfsTax != *dr.Data.TaxHash {
		return CodeTypeUnauthorized, errors.New("The tax is not a validated UUID.")
	}
	return CodeTypeOK, nil
}

func (tca *TCApplication) validateDelivery(dr DeliveryRequest) (uint32, error) {
	if dr.Data.Coins <= 0 {
		return CodeTypeUnauthorized, errors.New("Coins can not be the number of zero or negative.")
	}

	ver, err := dr.VerifySignature()
	if err != nil {
		return CodeTypeEncodingError, err
	}
	if !ver {
		return CodeTypeUnauthorized, errors.New("The signature does not validate the transaction.")
	}
	switch dr.Data.Action {
	case ADD_ACTION, REMOVE_ACTION:
		code, err := tca.validateInflators(dr)
		if err != nil {
			return code, err
		}
	case SEND_ACTION:
		code, err := tca.validateSend(dr)
		if err != nil {
			return code, err
		}
	}

	return CodeTypeOK, nil
}

func (tca *TCApplication) deliverAdd(dr DeliveryRequest) {
	from, _ := crypto.UnmarshalPublicKey(dr.Data.From)
	cj, _ := tca.state.GetCoins(from)
	coins := dr.Data.Coins + cj.Coins
	tca.state.SetCoins(from, coins)
}

func (tca *TCApplication) deliverRemove(dr DeliveryRequest) error {
	from, _ := crypto.UnmarshalPublicKey(dr.Data.From)
	cj, _ := tca.state.GetCoins(from)
	coins := cj.Coins - dr.Data.Coins
	if coins < 0 {
		return errors.New("You can not remove more than your requested.")
	}
	tca.state.SetCoins(from, coins)
	return nil
}

func (tca *TCApplication) deliverSend(dr DeliveryRequest) error {
	from, _ := crypto.UnmarshalPublicKey(dr.Data.From)
	fromCj, _ := tca.state.GetCoins(from)
	newFromCoins := fromCj.Coins - dr.Data.Coins
	if newFromCoins < 0 {
		return errors.New("You dont have enough money to send.")
	}
	tca.state.SetCoins(from, newFromCoins)

	taxCoins := dr.Data.Coins * float64(confs.Conf.Tax.Percentage) / 100
	toCoins := dr.Data.Coins - taxCoins

	to, _ := crypto.UnmarshalPublicKey(*dr.Data.To)
	toCj, _ := tca.state.GetCoins(to)
	newToCoins := toCj.Coins + toCoins
	tca.state.SetCoins(to, newToCoins)
	taxCj, _ := tca.state.GetCoins(confs.Conf.TaxReceiver)
	newTaxCoins := taxCj.Coins + taxCoins
	tca.state.SetCoins(confs.Conf.TaxReceiver, newTaxCoins)

	return nil
}

func (tca *TCApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	dr := DeliveryRequest{}
	err := json.Unmarshal(tx, &dr)
	if err != nil {
		return types.ResponseDeliverTx{Code: CodeTypeEncodingError, Log: "The json is not correct."}
	}

	code, err := tca.validateDelivery(dr)
	if err != nil {
		return types.ResponseDeliverTx{Code: code, Log: err.Error()}
	}

	switch dr.Data.Action {
	case ADD_ACTION:
		tca.deliverAdd(dr)
	case REMOVE_ACTION:
		err := tca.deliverRemove(dr)
		if err != nil {
			return types.ResponseDeliverTx{Code: CodeTypeUnauthorized, Log: err.Error()}
		}
	case SEND_ACTION:
		err := tca.deliverSend(dr)
		if err != nil {
			return types.ResponseDeliverTx{Code: CodeTypeUnauthorized, Log: err.Error()}
		}
	}

	return types.ResponseDeliverTx{Code: CodeTypeOK}
}
