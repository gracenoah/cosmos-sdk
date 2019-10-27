package store

import (
	dbm "github.com/gracenoah/tendermint/libs/db"

	"github.com/gracenoah/cosmos-sdk/store/rootmulti"
	"github.com/gracenoah/cosmos-sdk/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewPruningOptionsFromString(strategy string) (opt PruningOptions) {
	switch strategy {
	case "nothing":
		opt = PruneNothing
	case "everything":
		opt = PruneEverything
	case "syncable":
		opt = PruneSyncable
	default:
		opt = PruneSyncable
	}
	return
}
