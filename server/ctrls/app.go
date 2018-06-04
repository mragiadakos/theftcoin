package ctrls

import (
	"github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var _ types.Application = (*TCApplication)(nil)

type TCApplication struct {
	types.BaseApplication

	state State
}

func NewTCApplication() *TCApplication {
	state := loadState(dbm.NewMemDB())
	return &TCApplication{state: state}
}
