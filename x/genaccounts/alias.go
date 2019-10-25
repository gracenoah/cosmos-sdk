// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/gracenoah/cosmos-sdk/x/genaccounts/internal/types
package genaccounts

import (
	"github.com/gracenoah/cosmos-sdk/x/genaccounts/internal/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	// functions aliases
	NewGenesisAccountRaw        = types.NewGenesisAccountRaw
	NewGenesisAccount           = types.NewGenesisAccount
	NewGenesisAccountI          = types.NewGenesisAccountI
	GetGenesisStateFromAppState = types.GetGenesisStateFromAppState
	SetGenesisStateInAppState   = types.SetGenesisStateInAppState
	ValidateGenesis             = types.ValidateGenesis

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	GenesisAccount  = types.GenesisAccount
	GenesisAccounts = types.GenesisAccounts
	GenesisState    = types.GenesisState
)
