// nolint
package tags

import (
	sdk "github.com/gracenoah/cosmos-sdk/types"
)

var (
	ActionCompleteUnbonding    = "complete-unbonding"
	ActionCompleteRedelegation = "complete-redelegation"

	Action       = sdk.TagAction
	SrcValidator = sdk.TagSrcValidator
	DstValidator = sdk.TagDstValidator
	Delegator    = sdk.TagDelegator
	Moniker      = "moniker"
	Identity     = "identity"
	EndTime      = "end-time"
)
