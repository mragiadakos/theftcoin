package ctrls

import (
	"encoding/binary"

	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/abci/types"
)

func (tca *TCApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (tca *TCApplication) Commit() types.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, tca.state.Size)
	tca.state.AppHash = appHash
	tca.state.Height += 1
	saveState(tca.state)
	return types.ResponseCommit{Data: appHash}
}
