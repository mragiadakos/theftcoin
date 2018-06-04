package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs-api"

	"github.com/tendermint/abci/types"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/theftcoin/server/confs"
	"github.com/stretchr/testify/assert"
)

func TestQuerySuccessfully(t *testing.T) {
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
	dresp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, dresp.Code)
	qr := QueryRequest{}
	qr.Data.Date = time.Now().UTC()
	qr.Data.From, err = pubk.Bytes()
	assert.Nil(t, err)

	b, _ = json.Marshal(qr.Data)
	qr.Signature, err = privk.Sign(b)
	assert.Nil(t, err)

	b, _ = json.Marshal(qr)

	req := types.RequestQuery{}
	req.Data = b

	resp := app.Query(req)
	assert.Equal(t, CodeTypeOK, resp.Code)

	qresp := QueryResponse{}
	err = json.Unmarshal(resp.Value, &qresp)
	assert.Nil(t, err)

	assert.Equal(t, expectedCoins, qresp.Coins)
}

func TestQueryFailSignature(t *testing.T) {
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
	dresp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, dresp.Code)
	qr := QueryRequest{}
	qr.Data.Date = time.Now().UTC()
	qr.Data.From, err = pubk.Bytes()
	assert.Nil(t, err)

	b, _ = json.Marshal(qr.Data)
	qr.Signature, err = privk.Sign(b)
	assert.Nil(t, err)
	qr.Data.Nonce = "1"
	b, _ = json.Marshal(qr)

	req := types.RequestQuery{}
	req.Data = b

	resp := app.Query(req)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)

}

func TestQueryFailOnTime(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()

	confs.Conf.WaitingRequestTime = 1

	privk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	var expectedCoins float64 = 111
	dr := tu.inflatorCoins(t, privk, ADD_ACTION, expectedCoins)
	b, _ = json.Marshal(dr)
	dresp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, dresp.Code)
	qr := QueryRequest{}
	qr.Data.Date = time.Now().UTC()
	qr.Data.From, err = pubk.Bytes()
	assert.Nil(t, err)

	b, _ = json.Marshal(qr.Data)
	qr.Signature, err = privk.Sign(b)
	assert.Nil(t, err)
	b, _ = json.Marshal(qr)

	req := types.RequestQuery{}
	req.Data = b
	time.Sleep(1 * time.Second)
	resp := app.Query(req)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestQueryFailNotWatcher(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	fromPrivk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	var expectedCoins float64 = 111
	dr := tu.inflatorCoins(t, fromPrivk, ADD_ACTION, expectedCoins)
	b, _ = json.Marshal(dr)
	dresp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, dresp.Code)

	otherPrivk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	qr := QueryRequest{}
	data := QueryData{}
	data.Date = time.Now().UTC()
	data.From, err = otherPrivk.GetPublic().Bytes()
	assert.Nil(t, err)

	b, _ = fromPrivk.GetPublic().Bytes()
	data.User = &b
	assert.Nil(t, err)

	br, err := json.Marshal(data)
	assert.Nil(t, err)

	qr.Data = data
	qr.Signature, err = otherPrivk.Sign(br)

	assert.Nil(t, err)

	b, _ = json.Marshal(qr)
	req := types.RequestQuery{}
	req.Data = b

	resp := app.Query(req)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestQuerySuccesfullWatcher(t *testing.T) {
	tu := testUtils{}
	app := NewTCApplication()
	fromPrivk, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	b, _ := pubk.Bytes()
	hash := tu.addInflator(t, b)
	confs.Conf.IpfsInflators = hash
	confs.Conf.SubmitInflators()

	var expectedCoins float64 = 111
	dr := tu.inflatorCoins(t, fromPrivk, ADD_ACTION, expectedCoins)
	b, _ = json.Marshal(dr)
	dresp := app.DeliverTx(b)
	assert.Equal(t, CodeTypeOK, dresp.Code)

	otherPrivk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	qr := QueryRequest{}
	data := QueryData{}
	data.Date = time.Now().UTC()
	data.From, err = otherPrivk.GetPublic().Bytes()
	assert.Nil(t, err)

	buser, _ := fromPrivk.GetPublic().Bytes()
	data.User = &buser
	assert.Nil(t, err)
	br, err := json.Marshal(data)
	assert.Nil(t, err)

	qr.Data = data
	qr.Signature, err = otherPrivk.Sign(br)
	assert.Nil(t, err)
	hexPub := hex.EncodeToString(data.From)
	// add the watchers
	watchers := append([]confs.Watcher{}, confs.Watcher{PublicKeyHex: hexPub})
	watchB, _ := json.Marshal(watchers)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	watchHash, err := sh.BlockPut(watchB)
	assert.Nil(t, err)
	confs.Conf.IpfsWatchers = watchHash
	err = confs.Conf.SubmitWatchers()
	assert.Nil(t, err)

	b, _ = json.Marshal(qr)
	req := types.RequestQuery{}
	req.Data = b

	resp := app.Query(req)
	assert.Equal(t, CodeTypeOK, resp.Code)

	qresp := QueryResponse{}
	err = json.Unmarshal(resp.Value, &qresp)
	assert.Nil(t, err)

	assert.Equal(t, expectedCoins, qresp.Coins)
}
