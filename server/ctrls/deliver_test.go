package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ipfs/go-ipfs-api"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/theftcoin/server/confs"
	"github.com/stretchr/testify/assert"
)

type testUtils struct{}

func (tu *testUtils) inflatorCoins(t *testing.T, from crypto.PrivKey, action ActionStruct, coins float64) DeliveryRequest {
	var err error
	dd := DeliveryData{}
	dd.Action = action
	dd.Coins = coins
	dd.From, err = from.GetPublic().Bytes()
	assert.Nil(t, err)
	b, _ := json.Marshal(dd)
	dr := DeliveryRequest{}
	dr.Signature, err = from.Sign(b)
	assert.Nil(t, err)
	dr.Data = dd
	return dr
}

func (tu *testUtils) sendCoins(t *testing.T, from crypto.PrivKey, to crypto.PubKey, taxhash string, coins float64) DeliveryRequest {
	var err error
	dd := DeliveryData{}
	dd.Action = SEND_ACTION
	if to != nil {
		b, _ := to.Bytes()
		dd.To = &b
	}
	dd.Coins = coins
	dd.TaxHash = &taxhash
	dd.From, err = from.GetPublic().Bytes()
	assert.Nil(t, err)
	b, _ := json.Marshal(dd)
	dr := DeliveryRequest{}
	dr.Signature, err = from.Sign(b)
	assert.Nil(t, err)
	dr.Data = dd
	return dr
}

func (tu *testUtils) addInflator(t *testing.T, b []byte) string {
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	inf := confs.Inflator{}
	inf.PublicKeyHex = hex.EncodeToString(b)
	infs := []confs.Inflator{inf}
	infB, err := json.Marshal(infs)
	assert.Nil(t, err)
	hash, err := sh.BlockPut(infB)
	assert.Nil(t, err)
	return hash
}

func TestAnyTransactionFailSignature(t *testing.T) {
	tu := testUtils{}
	privk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	assert.Nil(t, err)
	dd := DeliveryData{}
	dd.Action = ADD_ACTION
	dd.From, err = pubk.Bytes()
	assert.Nil(t, err)
	b, _ = json.Marshal(dd)
	dr := DeliveryRequest{}
	dr.Signature, err = privk.Sign(b)
	assert.Nil(t, err)
	dr.Data = dd
	b, _ = json.Marshal(dr)
	app := NewTCApplication()
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestAnyTransactionSendingNegativeOrEqualToZeroCoins(t *testing.T) {
	app := NewTCApplication()
	tu := testUtils{}

	privk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	var expectedCoins float64 = -111
	dr := tu.inflatorCoins(t, privk, ADD_ACTION, expectedCoins)
	b, _ = json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)

	expectedCoins = 0
	dr = tu.inflatorCoins(t, privk, ADD_ACTION, expectedCoins)
	b, _ = json.Marshal(dr)
	resp = app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestAddCoinSuccesfully(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	privk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	var expectedCoins float64 = 111
	dr := tu.inflatorCoins(t, privk, ADD_ACTION, expectedCoins)
	b, _ = json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, resp.Code)
	cj, err := app.state.GetCoins(privk.GetPublic())
	assert.Nil(t, err)
	assert.Equal(t, expectedCoins, cj.Coins)
}

func TestRemoveCoinFailNegativeCoin(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	privk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	dr := tu.inflatorCoins(t, privk, REMOVE_ACTION, 111)
	b, _ = json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestRemoveCoinsSuccessfully(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	privk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	dr := tu.inflatorCoins(t, privk, ADD_ACTION, 111)
	b, _ = json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, resp.Code)

	dr = tu.inflatorCoins(t, privk, REMOVE_ACTION, 111)
	b, _ = json.Marshal(dr)
	resp = app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, resp.Code)

	cj, err := app.state.GetCoins(privk.GetPublic())
	assert.Nil(t, err)
	assert.Equal(t, float64(0), cj.Coins)
}

func TestSendCoinsFailMissingTo(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	fromPrivk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	//_, toPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	dr := tu.sendCoins(t, fromPrivk, nil, "", 111)

	b, _ := json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestSendCoinsFailTaxHash(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	fromPrivk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	_, toPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	dr := tu.sendCoins(t, fromPrivk, toPubk, "1111111111111111111111", 111)

	b, _ := json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestSendCoinsSuccessfullyTaxed(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()

	// creating the users
	fromPrivk, fromPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	fromPubkB, err := fromPubk.Bytes()
	assert.Nil(t, err)

	hash := tu.addInflator(t, fromPubkB)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	_, toPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	_, taxPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	// submitting the tax
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	taxPubkB, _ := taxPubk.Bytes()
	hexPub := hex.EncodeToString(taxPubkB)
	tax := confs.Tax{Percentage: 10, PublicKeyHex: hexPub}
	b, _ := json.Marshal(tax)
	taxHash, err := sh.BlockPut(b)
	assert.Nil(t, err)
	confs.Conf.IpfsTax = taxHash
	confs.Conf.SubmitTax()
	assert.Equal(t, confs.Conf.Tax, tax)

	money := 111.0
	// adding the from as an inflator so we have money
	dr := tu.inflatorCoins(t, fromPrivk, ADD_ACTION, money)
	b, _ = json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, resp.Code)

	// sending the money
	dr = tu.sendCoins(t, fromPrivk, toPubk, taxHash, money)
	b, _ = json.Marshal(dr)
	resp = app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, resp.Code)

	fromCj, err := app.state.GetCoins(fromPubk)
	assert.Nil(t, err)
	assert.Equal(t, float64(0), fromCj.Coins)

	taxedCoins := money * float64(tax.Percentage) / 100
	receiversCoins := money - taxedCoins
	receiverCj, err := app.state.GetCoins(toPubk)
	assert.Nil(t, err)
	assert.Equal(t, receiversCoins, receiverCj.Coins)

	taxerCj, err := app.state.GetCoins(taxPubk)
	assert.Nil(t, err)
	assert.Equal(t, taxedCoins, taxerCj.Coins)

}

func TestSendCoinsFailOnTryingSendingMoreThanHeHave(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()

	// creating the users
	fromPrivk, fromPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	fromPubkB, err := fromPubk.Bytes()
	assert.Nil(t, err)

	hash := tu.addInflator(t, fromPubkB)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	_, toPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	_, taxPubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	// submitting the tax
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	taxPubkB, _ := taxPubk.Bytes()
	hexPub := hex.EncodeToString(taxPubkB)
	tax := confs.Tax{Percentage: 10, PublicKeyHex: hexPub}
	b, _ := json.Marshal(tax)
	taxHash, err := sh.BlockPut(b)
	assert.Nil(t, err)
	confs.Conf.IpfsTax = taxHash
	confs.Conf.SubmitTax()
	assert.Equal(t, confs.Conf.Tax, tax)

	money := 11.0
	// adding the from as an inflator so we have money
	dr := tu.inflatorCoins(t, fromPrivk, ADD_ACTION, money)
	b, _ = json.Marshal(dr)
	resp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, resp.Code)

	// sending the money plus extra 100 coins
	dr = tu.sendCoins(t, fromPrivk, toPubk, taxHash, money+100)
	b, _ = json.Marshal(dr)
	resp = app.DeliverTx(b)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}
