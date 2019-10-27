package rpc

import (
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/gracenoah/tendermint/rpc/core/types"
)

var cdc = amino.NewCodec()

func init() {
	ctypes.RegisterAmino(cdc)
}
