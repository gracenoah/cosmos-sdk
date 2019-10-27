package app

import (
	"os"
	"testing"

	"github.com/gracenoah/cosmos-sdk/x/bank"
	"github.com/gracenoah/cosmos-sdk/x/crisis"

	"github.com/stretchr/testify/require"
	"github.com/gracenoah/tendermint/libs/db"
	"github.com/gracenoah/tendermint/libs/log"

	"github.com/gracenoah/cosmos-sdk/codec"
	"github.com/gracenoah/cosmos-sdk/x/auth"
	distr "github.com/gracenoah/cosmos-sdk/x/distribution"
	"github.com/gracenoah/cosmos-sdk/x/gov"
	"github.com/gracenoah/cosmos-sdk/x/mint"
	"github.com/gracenoah/cosmos-sdk/x/slashing"
	"github.com/gracenoah/cosmos-sdk/x/staking"

	abci "github.com/gracenoah/tendermint/abci/types"
)

func setGenesis(gapp *GaiaApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = NewGenesisAccount(acc)
	}

	genesisState := NewGenesisState(
		genaccs,
		auth.DefaultGenesisState(),
		bank.DefaultGenesisState(),
		staking.DefaultGenesisState(),
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		crisis.DefaultGenesisState(),
		slashing.DefaultGenesisState(),
	)

	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	gapp.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	gapp.Commit()

	return nil
}

func TestGaiadExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewGaiaApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	setGenesis(gapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewGaiaApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err := newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
