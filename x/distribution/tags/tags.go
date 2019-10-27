// nolint
package tags

import (
	sdk "github.com/gracenoah/cosmos-sdk/types"
)

// Distribution tx tags
var (
	Rewards    = "rewards"
	Commission = "commission"

	Validator = sdk.TagSrcValidator
	Delegator = sdk.TagDelegator
)
