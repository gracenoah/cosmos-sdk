package bank

import (
	sdk "github.com/gracenoah/cosmos-sdk/types"
)

// expected crisis keeper
type CrisisKeeper interface {
	RegisterRoute(moduleName, route string, invar sdk.Invariant)
}
