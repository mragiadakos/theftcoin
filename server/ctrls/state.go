package ctrls

import (
	"encoding/json"
	"errors"

	crypto "github.com/libp2p/go-libp2p-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	stateKey = []byte("stateKey")
	coinKey  = []byte("coinKey:")
)

func prefixCoinKey(pubk crypto.PubKey) ([]byte, error) {
	b, err := pubk.Bytes()
	if err != nil {
		return nil, errors.New("The public key is not correct")
	}
	return append(coinKey, b...), nil
}

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

type CoinJson struct {
	Coins float64
}

func (s *State) GetCoins(pubk crypto.PubKey) (CoinJson, error) {
	name, err := prefixCoinKey(pubk)
	if err != nil {
		return CoinJson{}, err
	}
	cj := CoinJson{}
	b := s.db.Get(name)
	json.Unmarshal(b, &cj)
	return cj, nil
}

func (s *State) SetCoins(pubk crypto.PubKey, coins float64) error {
	name, err := prefixCoinKey(pubk)
	if err != nil {
		return err
	}
	cj := CoinJson{Coins: coins}
	b, _ := json.Marshal(cj)
	s.db.Set(name, b)
	return nil
}

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}
